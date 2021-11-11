package models

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleSubAdmin Role ="subadmin"
	RoleUser  Role = "user"

)

func (r Role) IsValid() bool {
	return r == RoleAdmin || r == RoleUser||r==RoleSubAdmin
}
type User struct {
	ID        string     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Email     string     `json:"email" db:"email"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	Roles     []UserRole `json:"roles" db:"-"`
}
type UserRole struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	Role      Role      `json:"role" db:"role"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
type Restaurant struct{
	Name  string `json:"name db:"name"`
	Latitude string `json:"latitude" db:"latitude"`
	Longitude string `json:"longitude" db:"longitude"`
}
type Dish struct {
	ID string `json:"name" db:"name"`
}