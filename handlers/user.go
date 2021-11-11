package handlers

import (
	"SmallZomato/database"
	"SmallZomato/database/dbhelper"
	"SmallZomato/middlewares"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"

	"SmallZomato/models"
	"SmallZomato/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Name     string      `json:"name"`
		Email    string      `json:"email"`
		Password string      `json:"password"`
		Latitude string     `json:"latitude"`
		Longitude string     `json:"longitude"`
		Role     models.Role `json:"role"`
	}{}
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	if !body.Role.IsValid() {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid role type provided")
		return
	}
	if len(body.Password) < 6 {
		utils.RespondError(w, http.StatusBadRequest, nil, "password must be 6 chars long")
		return
	}

	exists, existsErr := dbhelper.IsUserExists(body.Role,body.Email)
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
	sessionToken := utils.HashString(body.Email + time.Now().String())
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, saveErr := dbhelper.CreateUser(tx, body.Name, body.Email, hashedPassword)
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
		addressErr := dbhelper.AddUserAddress(tx,userID,body.Latitude,body.Longitude)
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
func LoginUser(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

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
	sessionToken := utils.HashString(body.Email + time.Now().String())
	sessionErr := dbhelper.CreateUserSession(database.SmallZomato, userID, sessionToken)
	if sessionErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, sessionErr, "failed to create user session")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Token string `json:"token"`
	}{
		Token: sessionToken,
	})
}
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
func AddAddress(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	body := struct {
		Latitude    string `json:"latitude"`
		Longitude string `json:"longitude"`
	}{}
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	res,err := dbhelper.AddAuthUserAddress(userCtx.ID, body.Longitude,body.Longitude)
	if err != nil {

		utils.RespondJSON(w, http.StatusOK, res)

	}
	utils.RespondJSON(w, http.StatusCreated, res)

}

