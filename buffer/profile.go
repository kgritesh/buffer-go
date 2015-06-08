package buffer

import (
	"fmt"
	"net/http"
)

const (
	getProfileListURL      = "profiles.json"
	getProfileURL          = "profiles/%s.json"
	getProfileSchedulesURL = "profiles/%s/schedules.json"
)

// ProfileService handles communication with the profile related
// methods of the Buffer API.
//
// Buffer API docs: https://buffer.com/developers/api/profiles
type ProfileService struct {
	client *Client
}

//Schedule is an individual posting schedule which consists of days and times
type Schedule struct {
	Days  []string `json:"days"`
	Times []string `json:"times"`
}

//Statistics contains profile stats
type Statistics struct {
	Followers int `json:"followers"`
}

//Profile Represents a Buffer Profile
type Profile struct {
	Avatar            string     `json:"avatar,omitempty"`
	CreatedAt         Timestamp  `json:"created_at,omitempty"`
	Default           bool       `json:"default"`
	FormattedUsername string     `json:"formatted_username,omitempty"`
	ID                string     `json:"id,omitempty"`
	Schedules         []Schedule `json:"schedules,omitempty"`
	Service           string     `json:"service,omitempty"`
	ServiceID         string     `json:"service_id,omitempty"`
	ServiceUsername   string     `json:"service_username,omitempty"`
	Stats             Statistics `json:"statistics,omitempty"`
	TeamMembers       []string   `json:"team_members,omitempty"`
	Timezone          string     `json:"timezone,omitempty"`
	UserID            string     `json:"user_id,omitempty"`
}

//ListProfiles the current User all Profiles
func (profileService *ProfileService) ListProfiles() ([]Profile, *http.Response, error) {
	req, err := profileService.client.NewRequest("GET", getProfileListURL, nil)

	if err != nil {
		return nil, nil, err
	}

	profiles := new([]Profile)
	resp, err := profileService.client.Do(req, profiles)

	if err != nil {
		return nil, resp, err
	}
	fmt.Println("Profiles are ", profiles)
	return *profiles, resp, err
}

//GetProfile fetches the specified profile
//
//https://api.bufferapp.com/1/profiles/<profile_id>.json
func (profileService *ProfileService) GetProfile(profileID string) (*Profile, *http.Response, error) {
	url := fmt.Sprintf(getProfileURL, profileID)

	req, err := profileService.client.NewRequest("GET", url, nil)

	if err != nil {
		return nil, nil, err
	}

	profile := &Profile{}
	resp, err := profileService.client.Do(req, profile)

	if err != nil {
		return nil, resp, err
	}
	fmt.Println("Profile is ", profile)
	return profile, resp, err
}

//GetSchedule fetches all the posting schedule for the specified profile
//
//https://api.bufferapp.com/1/profiles/<profile_id>/schedules.json
func (profileService *ProfileService) GetSchedule(profileID string) ([]Schedule, *http.Response, error) {
	url := fmt.Sprintf(getProfileSchedulesURL, profileID)

	req, err := profileService.client.NewRequest("GET", url, nil)

	if err != nil {
		return nil, nil, err
	}

	schedules := new([]Schedule)
	resp, err := profileService.client.Do(req, schedules)

	if err != nil {
		return nil, resp, err
	}
	fmt.Println("Schedules are", *schedules)
	return *schedules, resp, err
}
