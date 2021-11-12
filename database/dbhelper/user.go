package dbhelper

import (
	"SmallZomato/database"
	"SmallZomato/models"
	"SmallZomato/utils"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)
func DeleteSessionToken(userID, token string) error {
	// language=SQL
	SQL := `DELETE FROM user_session WHERE user_id = $1 AND session_token = $2`
	_, err := database.SmallZomato.Exec(SQL, userID, token)
	return err
}
func IsUserExists(role models.Role,email string) (bool, error) {
	// language=SQL
	SQL := `SELECT u.id FROM user_profile AS u INNER JOIN user_roles AS v ON u.id = v.user_id  WHERE v.role =$1 AND u.email= $2 AND u.archived_at IS NULL`
	var id string
	err := database.SmallZomato.Get(&id, SQL, role,email)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)

		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}
func GetUserRoles(userID string) ([]models.UserRole, error) {
	// language=SQL
	SQL := `SELECT id, user_id, role FROM user_roles WHERE user_id = $1 AND archived_at IS NULL`
	roles := make([]models.UserRole, 0)
	err := database.SmallZomato.Select(&roles, SQL, userID)
	return roles, err
}
func CreateUser(db sqlx.Ext, name string, email string, password string,created models.Role) (string, error) {
	// language=SQL
	SQL := `INSERT INTO user_profile(name, email, password,created_by) VALUES ($1, TRIM(LOWER($2)), $3,$4) RETURNING id`
	var userID string
	if err := db.QueryRowx(SQL, name, email, password,created).Scan(&userID); err != nil {
		return "", err
	}
	return userID, nil
}
func CreateUserRole(db sqlx.Ext, userID string, role models.Role) error {
	// language=SQL
	SQL := `INSERT INTO user_roles(user_id, role) VALUES ($1, $2)`
	_, err := db.Exec(SQL, userID, role)
	return err
}

func CreateUserSession(db sqlx.Ext, userID, sessionToken string) error {
	// language=SQL
	SQL := `INSERT INTO user_session(user_id, session_token) VALUES ($1, $2)`
	_, err := db.Exec(SQL, userID, sessionToken)
	return err
}
func GetUserBySession(sessionToken string) (*models.User, error) {
	// language=SQL
	SQL := `SELECT 
       			u.id, 
       			u.name, 
       			u.email, 
       			u.created_at 
			FROM user_profile u
			JOIN user_session us on u.id = us.user_id
			WHERE u.archived_at IS NULL AND us.session_token = $1`
	var user models.User
	err := database.SmallZomato.Get(&user, SQL, sessionToken)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	roles, roleErr := GetUserRoles(user.ID)
	if roleErr != nil {
		return nil, roleErr
	}
	user.Roles = roles
	return &user, nil
}
func AddUserAddress(db sqlx.Ext,userID,latitude,longitude string)error{
	fmt.Println(userID)
	//language SQL
	SQL := `INSERT INTO user_address (user_id,latitude,longitude) VALUES($1,$2,$3)`
	_,err := db.Exec(SQL,userID,latitude,longitude)
	if err!=nil{
		return err
	}
	return nil
}
func GetUserIDByPassword(email, password string) (string, error) {
// language=SQL
SQL := `SELECT
				u.id,
       			u.password
       		FROM
				user_profile u
			WHERE
				u.archived_at IS NULL
				AND u.email = TRIM(LOWER($1))`
var userID string
var passwordHash string
err := database.SmallZomato.QueryRowx(SQL, email).Scan(&userID, &passwordHash)
if err != nil && err != sql.ErrNoRows {
return "", err
}
if err == sql.ErrNoRows {
return "", nil
}
// compare password
if passwordErr := utils.CheckPassword(password, passwordHash); passwordErr != nil {
return "", passwordErr
}
return userID, nil
}


func AddAuthUserAddress(userID,latitude,longitude string)(map[string]string,error){
	//language = SQL
	SQL := `INSERT INTO user_address (user_id,latitude,longitude) VALUES($1,$2,$3)`
	_,err := database.SmallZomato.Exec(SQL,userID,latitude,longitude)
	res := make(map[string]string)
	if err!=nil{
		res["error"]="Error in insertion of address"
		return res,err
	}
	res["success"]="Address Inserted Successfully"
	return res,nil

}
func GetAllRestaurant()([]models.Restaurant,error){
	SQL := `SELECT name,latitude,longitude from restaurants`
	res :=make([]models.Restaurant,0)

	err := database.SmallZomato.Select(&res,SQL)
	fmt.Println(res)
	if err!=nil{
		return nil,err
	}
	return res,nil

}
func GetRestaurantDish(id string)([]models.Dish,error){
	SQL := `SELECT v.name FROM dishes AS v INNER JOIN restaurants as r ON v.res_id=r.id WHERE r.id =$1`
	res := make([]models.Dish,0)
	err:=database.SmallZomato.Select(&res,SQL,id)
	if err!=nil{
		fmt.Println(err)
		return nil,err
	}
	return res,nil


}