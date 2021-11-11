package dbhelper

import (
	"SmallZomato/database"
)

func RestaurantAdd(name,latitude,longitude string)(map[string]string ,error){

	SQL2 := `INSERT INTO  restaurants(name,latitude,longitude,created_by) VALUES($1,$2,$3,$4)`
	_,err2:=database.SmallZomato.Exec(SQL2,name,latitude,longitude,"admin")

	res := make(map[string]string)
	if err2!=nil{

			res["error"]= "Problem in insertion"

			return res,err2


	}
	res["success"]="inserted successfully"
	return res,nil

}