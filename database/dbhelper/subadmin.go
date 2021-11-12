package dbhelper

import (
	"SmallZomato/database"
	"SmallZomato/models"
)

func SubAdminRestaurantAdd(name,latitude,longitude string)(map[string]string ,error){
	SQL2 := `INSERT INTO  restaurants(name,latitude,longitude,created_by) VALUES($1,$2,$3,$4)`
	_,err2:=database.SmallZomato.Exec(SQL2,name,latitude,longitude,"subadmin")

	res := make(map[string]string)
	if err2!=nil{

		res["error"]= "Problem in insertion"

		return res,err2


	}
	res["success"]="inserted successfully"
	return res,nil


}
func AddSubAdminDishToRes(name,id string)(map[string]string ,error){
	SQL2 := `INSERT INTO  dishes(name,res_id,created_by) VALUES($1,$2,$3)`
	_,err2:=database.SmallZomato.Exec(SQL2,name,id,"subadmin")
	res := make(map[string]string)
	if err2!=nil{

		res["error"]= "Problem in insertion"

		return res,err2


	}
	res["success"]="inserted successfully"
	return res,nil

}
func GetUser()([]models.UserDetail,error){
	SQL := `SELECT name,email,created_at FROM user_profile WHERE created_by =$1`
	res := make([]models.UserDetail,0)
	err:=database.SmallZomato.Select(&res,SQL,"subadmin")

	if err!=nil{

		return nil,err
	}
	return res,nil

}
