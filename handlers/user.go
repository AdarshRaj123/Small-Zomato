package handlers

import (
	"SmallZomato/database"
	"SmallZomato/database/dbhelper"
	"SmallZomato/middlewares"
	"SmallZomato/models"
	"SmallZomato/utils"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"net/http"
)

// RegisterUser godoc
// @Summary      Creates A user
// @Description  takes email,password etc and creates the user and sends the session token
// @Tags         small-zomato/registeruser
// @Accept       json
// @Produce      json
// @Param        register body models.RegisterUser true "register request"
// @Success      200  {object} models.UserLoginResponse
// @Failure      400  {object} utils.RequestErr
// @Failure      401  {object} utils.RequestErr
// @Failure      500  {object} utils.RequestErr
// @Router       /small-zomato/register [post]
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	body := models.RegisterUser{}
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	//todo: use errors.new
	if !body.Role.IsValid() {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid role type provided")
		return
	}
	if len(body.Password) < 6 {
		utils.RespondError(w, http.StatusBadRequest, nil, "password must be 6 chars long")
		return
	}

	exists, existsErr := dbhelper.IsUserExists(body.Role, body.Email)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check user existence")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "user already exists")
		return
	}
	hashedPassword, hasErr := utils.HashPassword(body.Password)
	if hasErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, hasErr, "failed to secure password")
		return
	}
	//sessionToken := utils.HashString(body.Email + time.Now().String())
	sessionToken, _ := utils.GenerateJWT(body.Email, body.Role)

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, saveErr := dbhelper.CreateUser(tx, body.Name, body.Email, hashedPassword, body.Role, body.MobileNo)
		if saveErr != nil {
			return saveErr
		}
		roleErr := dbhelper.CreateUserRole(tx, userID, body.Role)
		if roleErr != nil {
			return roleErr
		}
		sessionErr := dbhelper.CreateUserSession(tx, userID, sessionToken)
		if sessionErr != nil {
			return sessionErr
		}
		addressErr := dbhelper.AddUserAddress(tx, userID, body.Latitude, body.Longitude)
		if addressErr != nil {
			return addressErr
		}
		return nil
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to create user")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, struct {
		Token string `json:"token"`
	}{
		Token: sessionToken,
	})
}

// LoginUser godoc
// @Summary      Allow the user (non admin) to log into the system
// @Description  Login API takes in the email and password and returns the session token if login is valid
// @Tags         v1/login
// @Accept       json
// @Produce      json
// @Param        login body models.UserLoginRequest true "Login request"
// @Success      200  {object} models.UserLoginResponse
// @Failure      400  {object} utils.RequestErr
// @Failure      401  {object} utils.RequestErr
// @Failure      500  {object} utils.RequestErr
// @Router       /small-zomato/login [post]
func LoginUser(w http.ResponseWriter, r *http.Request) {
	body := models.UserLoginRequest{}

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	userID, userErr := dbhelper.GetUserIDByPassword(body.Email, body.Password)
	if userErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, userErr, "failed to find user")
		return
	}
	// create user session

	role, _ := dbhelper.GetUserRole(body.Email)
	sessionToken, _ := utils.GenerateJWT(body.Email, models.Role(role))

	//sessionToken := utils.HashString(body.Email + time.Now().String())

	sessionErr := dbhelper.CreateUserSession(database.SmallZomato, userID, sessionToken)
	if sessionErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, sessionErr, "failed to create user session")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, models.UserLoginResponse{Token: sessionToken})
}

// GetUserInfo godoc
// @Summary      Gets the details of the user
// @Description  the details of the user is returned
// @Tags         small-zomato/getUserInfo
// @Produce      json
// @Success      200  {object} models.User
// @Failure      400  {object} utils.RequestErr
// @Failure      401  {object} utils.RequestErr
// @Failure      500  {object} utils.RequestErr
// @Router       /small-zomato/user/info [post]
// @Security     ApiKeyAuth
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	utils.RespondJSON(w, http.StatusOK, userCtx)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-api-key")
	userCtx := middlewares.UserContext(r)
	err := dbhelper.DeleteSessionToken(userCtx.ID, token)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to logout user")
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// AddAddress godoc
// @Summary      it adds the address of the user
// @Description  takes the latitude and longitude of the users as address
// @Tags         Small-Zomato/add-address
// @Accept       json
// @Produce      json
// @Param        user body  models.UserLocation  true "users location"
// @Success      200  {object}  models.UserLocation
// @Failure      400  {object} utils.RequestErr
// @Failure      401  {object} utils.RequestErr
// @Failure      500  {object} utils.RequestErr
// @Router       /small-zomato/user/address [post]
// @Security     ApiKeyAuth
func AddAddress(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	body := models.UserLocation{}
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	res, err := dbhelper.AddAuthUserAddress(userCtx.ID, body.Longitude, body.Longitude)
	if err != nil {

		utils.RespondJSON(w, http.StatusInternalServerError, res)
		return

	}
	utils.RespondJSON(w, http.StatusCreated, res)

}
func GetAll(w http.ResponseWriter, r *http.Request) {
	res, err := dbhelper.GetAllRestaurant()
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, res)
		return
	}
	utils.RespondJSON(w, http.StatusOK, res)

}
func GetDish(w http.ResponseWriter, r *http.Request) {
	body := struct {
		ID string `json:"name"`
	}{}
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	res, err := dbhelper.GetRestaurantDish(body.ID)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, res)
		return
	}
	utils.RespondJSON(w, http.StatusOK, res)

}
func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully Uploaded File\n")

}
