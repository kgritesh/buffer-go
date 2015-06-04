package buffer

import (
	"fmt"
	"net/http"
)

const (
	getUserURL     = "user.json"
	deauthorizeURL = "user/deauthorize.json"
)

// UserService handles communication with the user related
// methods of the User API.
//
// GitHub API docs: http://developer.github.com/v3/repos/
type UserService struct {
	client *Client
}

//User Respresents a Buffer User
type User struct {
	ID         string    `json:"id,omitempty"`
	ActivityAt Timestamp `json:"activity_at,omitempty"`
	CreatedAt  Timestamp `json:"created_at,omitempty"`
	Plan       string    `json:"plan,omitempty"`
	Timezone   string    `json:"timezone,omitempty"`
}

//Get the current User Details
func (userService *UserService) Get() (*User, *http.Response, error) {
	req, err := userService.client.NewRequest("GET", getUserURL, nil)

	if err != nil {
		return nil, nil, err
	}

	user := &User{}
	resp, err := userService.client.Do(req, user)

	if err != nil {
		return nil, resp, err
	}
	fmt.Println("User is ", user)
	return user, resp, err
}
