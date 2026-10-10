package voting

import (
	"errors"
)

var (
	UserNotFound = errors.New("User Not Found")
	ProjectAlreadyExists = errors.New("Project Already Exists")
	UserAlreadyVoted = errors.New("The User has Already Voted")
	UserAlreadyExists = errors.New("The user is already registered")
	AdminAlreadyExists = errors.New("Admin already exists for this Email")
	AdminNotFound = errors.New("No admin Found for the given credentials")
	)


type User struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Organization string `json:"organization"`
	Address string `json:"address"`
	Qrcode string `json:"qrcode"`
}

type Project struct {
		Id           string `json:"id"`
    Name         string `json:"name"`
    Members      string `json:"members"`
    Organization string `json:"organization"`
    Description  string `json:"description,omitempty"`
}
type VotingOperations interface {
	ListProjects() []Project
	GetUser() User
}
 type Admin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
