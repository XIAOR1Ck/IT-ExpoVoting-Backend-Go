package voting

import (
	"errors"
)

var (
	UserNotFound = errors.New("User Not Found")
	ProjectAlreadyExists = errors.New("Project Already Exists")
	UserAlreadyVoted = errors.New("The User has Already Voted")
	)
type Voting struct {
	user User
	project Project
}

type User struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Organization string `json:"organization"`
	Address string `json:"address"`
}

type Project struct {
		Id           string `json:"id"`
    Name         string `json:"name"`
    Members      string `json:"members"`
    Organization string `json:"organization"`
    Description  string `json:"description,omitempty"`
    Votes        int    `json:"votes"`}

type VotingOperations interface {
	ListProjects() []Project
	GetUser() User
}
