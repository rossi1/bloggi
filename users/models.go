package users

import (
	"time"
)

type Response map[string]interface{}

type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	FullName  string    `json:"fullname" validate:"required" gorm:"type:varchar(100)"`
	Username  string    `json:"username" validate:"required,unique_username" gorm:"type:varchar(100);unique_index"`
	Password  string    `json:"password,-" validate:"required" gorm:"type:varchar(1000)" `
}

// RequestResponse
func (u User) requestResponse(token string) Response {

	switch {
	case token == "":
		response := Response{
			"username": u.Username,
			"fullname": u.FullName,
			"id":       u.ID,
		}
		return response

	default:
		response := Response{
			"username": u.Username,
			"fullname": u.FullName,
			"id":       u.ID,
			"token":    token,
		}
		return response

	}

}
