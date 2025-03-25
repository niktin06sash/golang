package model

import "github.com/google/uuid"

type Person struct {
	Id       uuid.UUID `json:"userid"`
	Name     string    `json:"username"`
	Email    string    `json:"useremail"`
	Password string    `json:"userpassword"`
}
