package handlers

import (
	"SmallZomato/database/dbhelper"
	"SmallZomato/utils"
	"fmt"
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
