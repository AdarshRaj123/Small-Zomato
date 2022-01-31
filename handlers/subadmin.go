package handlers

import (
	"SmallZomato/database"
	"SmallZomato/database/dbhelper"
	"SmallZomato/models"
	"SmallZomato/utils"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func AddSubAdminRestaurant(w http.ResponseWriter, r *http.Request){

	body := struct{
		Name string `json:"name"`
		Latitude string `json:"latitude"`
		Longitude string `json:"longitude"`

	}{}
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	res,err:= dbhelper.SubAdminRestaurantAdd(body.Name,body.Latitude,body.Longitude)
	if err!=nil{
		fmt.Println(err)
		utils.RespondJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondJSON(w,http.StatusCreated,res)


}
func AddSubAdminDish(w http.ResponseWriter, r *http.Request){
	body := struct{
		Name string `json:"name"`
		ID  string `json:"res_id"`

	}{}

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	res,err:= dbhelper.AddSubAdminDishToRes(body.Name,body.ID)
	if err!=nil{
		fmt.Println(err)
		utils.RespondJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondJSON(w,http.StatusCreated,res)


}
func AddSubAdminUser(w http.ResponseWriter, r *http.Request){
	body := struct {
		Name     string      `json:"name"`
		Email    string      `json:"email"`
		Password string      `json:"password"`
		Latitude string     `json:"latitude"`
		Longitude string     `json:"longitude"`
		MobileNo string        `json:"mobile_no"`
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
		userID, saveErr := dbhelper.CreateUser(tx, body.Name, body.Email, hashedPassword,"subadmin",body.MobileNo)
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
func GetSubAdminUsers(	w http.ResponseWriter, r *http.Request){
	res,err := dbhelper.GetUser()
	if err!=nil{

		utils.RespondJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondJSON(w,http.StatusOK,res)
}

