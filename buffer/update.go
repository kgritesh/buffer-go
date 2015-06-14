package buffer

import (
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

const (
	getUpdateURL          = "updates/%s.json"
	getPendingUpdates     = "profiles/%s/updates/pending.json"
	getSentUpdates        = "profiles/%s/updates/sent.json"
	getUpdateInteractions = "updates/%s/interactions.json"
	reorderUpdates        = "profiles/%s/updates/reorder.json"
	shuffleUpdates        = "profiles/%s/updates/shuffle.json"
	createUpdate          = "updates/create.json"
	editUpdate            = "updates/%s/update.json"
	shareUpdate           = "updates/%s/share.json"
	deleteUpdate          = "updates/%s/destroy.json"
	moveUpdateToTop       = "updates/%s/move_to_top.json"
)

// UpdateService handles communication with the update related
// methods of the Buffer API.
//
// Buffer API docs: https://buffer.com/developers/api/updates
type UpdateService struct {
	client *Client
}

//UpdateStats holds the statistics for each social post
type UpdateStats struct {
	Reach      int `json:"reach"`
	Clicks     int `json:"clicks"`
	Retweets   int `json:"retweets"`
	Favourites int `json:"favourites"`
	Mentions   int `json:"mentions"`
}

//Update represents a single post to a single social media account
type Update struct {
	ID              string      `json:"id,omitempty"`
	CreatedAt       Timestamp   `json:"created_at,omitempty"`
	Day             string      `json:"string,omitempty"`
	DueAt           Timestamp   `json:"due_at,omitempty"`
	DueTime         string      `json:"due_time,omitempty"`
	ProfileID       string      `json:"profile_id,omitempty"`
	ProfileService  string      `json:"profile_service,omitempty"`
	SentAt          Timestamp   `json:"sent_at,omitempty"`
	ServiceUpdateID string      `json:"service_update_id,omitempty"`
	Statistics      UpdateStats `json:"statistics,omitempty"`
	Status          string      `json:"status,omitempty"`
	Text            string      `json:"text,omitempty"`
	TextFormatted   string      `json:"text_formatted,omitempty"`
	UserID          string      `json:"user_id,omitempty"`
	Via             string      `json:"via,omitempty"`
}

//UpdateList represents list of updates
type UpdateList struct {
	Total   int      `json:"total,omitempty"`
	Updates []Update `json:"updates"`
}

// UpdateListOptions specifies the optional parameters to the
// UpdateService.Pending and UpdateService.sent methods
type UpdateListOptions struct {
	Page   int       `url:"page,omitempty"`
	Count  int       `url:"count,omitempty"`
	Since  Timestamp `url:"since,omitempty"`
	UTCSet bool      `url:"utc,omitempty"`
}

//GetUpdate fetches the specified profile
//
//https://api.bufferapp.com/1/updates/:id
func (updateService *UpdateService) GetUpdate(updateID string) (*Update, *http.Response, error) {
	url := fmt.Sprintf(getUpdateURL, updateID)

	req, err := updateService.client.NewRequest("GET", url, nil)

	if err != nil {
		return nil, nil, err
	}

	update := &Update{}
	resp, err := updateService.client.Do(req, update)

	if err != nil {
		return nil, resp, err
	}
	fmt.Println("Update is ", update)
	return update, resp, err
}

// getUpdateList is a helper method to fetch list of updates from the  provided url,
// and fullfilling provided conditions
func (updateService *UpdateService) getUpdateList(url string, profileID string, opt *UpdateListOptions) (*UpdateList, *http.Response, error) {
	url, err := addOptions(url, opt)

	if err != nil {
		return nil, nil, err
	}

	req, err := updateService.client.NewRequest("GET", url, nil)

	if err != nil {
		return nil, nil, err
	}

	updateList := &UpdateList{}
	resp, err := updateService.client.Do(req, updateList)

	if err != nil {
		return nil, resp, err
	}
	fmt.Println("Update List is ", updateList)
	return updateList, resp, err
}

//GetPendingUpdates fetches the pending updates for the specified profile
//
//https://api.bufferapp.com/1/profiles/:id/updates/pending
func (updateService *UpdateService) GetPendingUpdates(profileID string,
	opt *UpdateListOptions) (*UpdateList, *http.Response, error) {
	url := fmt.Sprintf(getPendingUpdates, profileID)
	return updateService.getUpdateList(url, profileID, opt)
}

//GetSentUpdates fetches the sent updates for the specified profile
//
//https://api.bufferapp.com/1/profiles/:id/updates/sent
func (updateService *UpdateService) GetSentUpdates(profileID string,
	opt *UpdateListOptions) (*UpdateList, *http.Response, error) {

	url := fmt.Sprintf(getSentUpdates, profileID)
	return updateService.getUpdateList(url, profileID, opt)
}

func (updateService *UpdateService) makeUpdatePost(url string,
	opt interface{}) (map[string]interface{}, *http.Response, error) {

	req, err := updateService.client.NewRequest("POST", url, opt)

	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]interface{})
	resp, err := updateService.client.Do(req, &result)

	if err != nil {
		return nil, resp, err
	}
	fmt.Println("Response is ", result)
	return result, resp, err

}

//UpdateReorderOptions specifies the optional parameters to
// UpdateService.Reorder methods
type UpdateReorderOptions struct {
	Order  []string `json:"order"`
	Offset int      `json:"offset,omitempty"`
	UTCSet bool     `json:"utc,omitempty"`
}

//ReorderUpdates reorder the pending updates.
//
//https://api.bufferapp.com/1/profiles/:id/updates/reorder.json
func (updateService *UpdateService) ReorderUpdates(profileID string,
	opt *UpdateReorderOptions) (bool, *http.Response, error) {

	url := fmt.Sprintf(reorderUpdates, profileID)
	result, resp, err := updateService.makeUpdatePost(url, opt)

	if err != nil {
		return false, resp, err
	}

	return result["success"].(bool), resp, err
}

//UpdateShuffleOptions specifies the optional parameters to
// UpdateService.Shuffle methods
type UpdateShuffleOptions struct {
	Count  int  `json:"count,omitempty"`
	UTCSet bool `json:"utc,omitempty"`
}

//ShuffleUpdates shuffle the pending updates.
//
//https://api.bufferapp.com/1/profiles/:id/updates/shuffle.json
func (updateService *UpdateService) ShuffleUpdates(profileID string,
	opt *UpdateShuffleOptions) (bool, *http.Response, error) {

	url := fmt.Sprintf(shuffleUpdates, profileID)

	result, resp, err := updateService.makeUpdatePost(url, opt)

	if err != nil {
		return false, resp, err
	}

	return result["success"].(bool), resp, err
}

//UpdateCreateOptions specifies the optional parameters to
// UpdateService.Create methods
type UpdateCreateOptions struct {
	ProfileIDList []string            `url:"profile_ids,brackets"`
	Text          string              `url:"text,omitempty"`
	Shorten       bool                `url:"shorten,omitempty"`
	Now           bool                `url:"bool,omitempty"`
	Top           bool                `url:"top,omitempty"`
	Media         []map[string]string `url:"media,omitempty"`
	Attachment    bool                `url:"attachment,omitempty"`
}

//CreateUpdates creates a new update
//
//https://api.bufferapp.com/1/updates/create.json
func (updateService *UpdateService) CreateUpdates(
	opt *UpdateCreateOptions) (*[]Update, *http.Response, error) {

	url := createUpdate
	req, err := updateService.client.NewRequest("POST", url, opt)

	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]interface{})
	resp, err := updateService.client.Do(req, &result)

	if err != nil {
		return nil, resp, err
	}

	updateList := new([]Update)

	err = mapstructure.Decode(result["updates"], updateList)

	if err != nil {
		return nil, resp, err
	}

	return updateList, resp, err
}

//UpdateEditOptions specifies the optional parameters to
// UpdateService.Edit methods
type UpdateEditOptions struct {
	Text   string              `url:"text"`
	Now    bool                `url:"bool,omitempty"`
	Media  []map[string]string `url:"media,omitempty"`
	UTCSet bool                `url:"utc,omitempty"`
}

//EditUpdate creates a new update
//
//https://api.bufferapp.com/1/updates/update.json
func (updateService *UpdateService) EditUpdate(updateID string,
	opt *UpdateEditOptions) (*Update, *http.Response, error) {

	url := fmt.Sprintf(editUpdate, updateID)
	req, err := updateService.client.NewRequest("POST", url, opt)

	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]interface{})
	resp, err := updateService.client.Do(req, &result)

	if err != nil {
		return nil, resp, err
	}

	update := new(Update)

	err = mapstructure.Decode(result["update"], update)

	if err != nil {
		return nil, resp, err
	}

	return update, resp, err
}
