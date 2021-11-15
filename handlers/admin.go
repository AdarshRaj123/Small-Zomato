package handlers

import (
	"SmallZomato/database"
	"SmallZomato/database/dbhelper"
	"SmallZomato/middlewares"
	"SmallZomato/models"
	"SmallZomato/utils"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func AddRestaurant(w http.ResponseWriter, r *http.Request){

	body := struct{
		Name string `json:"name"`
		Latitude string `json:"latitude"`
		Longitude string `json:"longitude"`

	}{}
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	  res,err:= dbhelper.RestaurantAdd(body.Name,body.Latitude,body.Longitude)
	  if err!=nil{
		  fmt.Println(err)
		  utils.RespondJSON(w, http.StatusInternalServerError, err)
		  return
	  }
	utils.RespondJSON(w,http.StatusOK,res)


}
func AddDish(w http.ResponseWriter, r *http.Request){
	body := struct{
		Name string `json:"name"`
		 ID  string `json:"res_id"`
	}{}

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	res,err:= dbhelper.AddDishToRes(body.Name,body.ID)
	if err!=nil{
		fmt.Println(err)
		utils.RespondJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondJSON(w,http.StatusOK,res)


}
func GetUsers(w http.ResponseWriter, r *http.Request){
	 res,err := dbhelper.GetAllUser()

	if err!=nil{

		utils.RespondJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondJSON(w,http.StatusOK,res)
}
func AddUser(w http.ResponseWriter, r *http.Request){
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

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, saveErr := dbhelper.CreateUser(tx, body.Name, body.Email, hashedPassword,"admin")
		if saveErr != nil {
			return saveErr
		}
		roleErr := dbhelper.CreateUserRole(tx, userID, body.Role)
		if roleErr != nil {
			return roleErr
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
		message string `json:"message"`
	}{
		message: "user created",
	})
}
func GetDistance(w http.ResponseWriter, r *http.Request){
	userCtx := middlewares.UserContext(r)
	body := struct {

		Latitude string     `json:"latitude"`
		Longitude string     `json:"longitude"`

	}{}
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	res:= dbhelper.CalculateDistance(userCtx.ID,body.Latitude,body.Longitude)
	jdata := make(map[string]float64)
	jdata["distance"]=res
	utils.RespondJSON(w,http.StatusOK,jdata)

}