package dbhelper

import (
	"SmallZomato/database"
	"SmallZomato/models"
	"SmallZomato/utils"
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
func AddDishToRes(name,id string)(map[string]string ,error){
	SQL2 := `INSERT INTO  dishes(name,res_id,created_by) VALUES($1,$2,$3)`
	_,err2:=database.SmallZomato.Exec(SQL2,name,id,"admin")
	res := make(map[string]string)
	if err2!=nil{

		res["error"]= "Problem in insertion"

		return res,err2


	}
	res["success"]="inserted successfully"
	return res,nil

}
func GetAllUser()([]models.UserDetail,error){
	//language SQL
	SQL := `SELECT name,email,created_at FROM user_profile`
	res := make([]models.UserDetail,0)
	err:=database.SmallZomato.Select(&res,SQL)

	if err!=nil{

		return nil,err
	}
	return res,nil

}
func CalculateDistance(id,latitude,longitude string)float64{
	//language SQL
	SQL:=`SELECT latitude,longitude FROM user_address WHERE user_id =$1`
	var data models.Location
	err :=database.SmallZomato.Get(&data,SQL,id)
	if err!=nil{

	}
	  dist := utils.GetDistanceFromLat(data,latitude,longitude)

	  return dist

}
