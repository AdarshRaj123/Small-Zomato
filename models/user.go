package models

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleSubAdmin Role = "subadmin"
	RoleUser     Role = "user"
)

func (r Role) IsValid() bool {
	return r == RoleAdmin || r == RoleUser || r == RoleSubAdmin
}

// User contains the info of the user
type User struct {
	//ID of the user
	//Required: true
	ID string `json:"id" db:"id"`
	//Name of the user
	//Required: true
	Name string `json:"name" db:"name"`
	//Email of the user
	//Required: true
	Email string `json:"email" db:"email"`
	//CreatedAt of the user
	//Required: true
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	//Roles of the user
	//Required: true
	Roles []UserRole `json:"roles" db:"-"`
} // @ name User
type UserRole struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	Role      Role      `json:"role" db:"role"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
type Restaurant struct {
	Name      string `json:"name" db:"name"`
	Latitude  string `json:"latitude" db:"latitude"`
	Longitude string `json:"longitude" db:"longitude"`
}
type Dish struct {
	Name string `json:"name" db:"name"`
}

// UserDetail models contains the details of the user
type UserDetail struct {
	//Name of the user
	//Required:true
	Name string `json:"name" db:"name"`
	//Email of the user
	//Required: true
	Email string `json:"email" db:"email"`
	//CreatedAt is when the user was created
	//Required: false
	CreatedAt string `json:"created_at" db:"created_at"`
	//MobileNo is the mobile no of the user
	//Required: true
	MobileNo string `json:"mobile_no" db:"mobile_no"`
} // @ name UserDetail

type Location struct {
	Latitude  string `json:"latitude" db:"latitude"`
	Longitude string `json:"longitude" db:"longitude"`
}
type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

// UserLoginRequest models contains the body data for the login request
type UserLoginRequest struct {
	// Email address of the user
	// Required: true
	// Example: test@someemail.com
	Email string `json:"email"`
	// Password for the user ac
	// Required: true
	// Example: some-password
	Password string `json:"password"`
} // @ name UserLoginRequest
// UserLoginResponse models contains the token when user is successfully authorized
type UserLoginResponse struct {
	// Token for the user session, should be sent in x-api-key for the authorized requests
	// Example: 8e77e71abe427ced1c93d883aeeddfa57ce39b787f229caaf176fdd71353f3466d340a2cdb5a219c429c53ad37f2f144c7ce01b985b6b33e397c4b8fd1433cc3
	Token string `json:"token"`
} // @ name UserLoginResponse

//RegisterUser modal contains the register user body
type RegisterUser struct {
	//Name of the user
	//Required:true
	Name string `json:"name"`
	//Email of the user
	//Required: true
	Email string `json:"email"`
	//Password of the user
	//Required: true
	Password string `json:"password"`
	//MobileNo of the user
	//Required: true
	MobileNo string `json:"mobile_no"`
	//Latitude of the user
	Latitude string `json:"latitude"`
	//Longitude of the user
	//Required: true
	Longitude string `json:"longitude"`
	//Role of the user
	//Required: true
	Role Role `json:"role"`
} //@name RegisterUser

// UserLocation models contains the location of the user
type UserLocation struct {
	// Latitude is the latitude of the user
	// Required: true
	Latitude string `json:"latitude"`
	// Longitude is the longitude of the user
	// Required: true
	Longitude string `json:"longitude"`
} //@name UserLocation
