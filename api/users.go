package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type memberResponse struct {
	Ok      bool     `json:"ok"`
	Members []Member `json:"members"`
	Error   string   `json:"error"`
}

type Member struct {
	ID       string      `json:"id"`
	TeamID   string      `json:"team_id"`
	Name     string      `json:"name"`
	Deleted  bool        `json:"deleted"`
	Status   string      `json:"status"`
	Color    string      `json:"color"`
	RealName string      `json:"real_name"`
	TZ       string      `json:"tz"`
	TZLabel  string      `json:"tz_label"`
	TZOffset int32       `json:"tz_offset"`
	IsAdmin  bool        `json:"is_admin"`
	IsOwner  bool        `json:"is_owner"`
	IsBot    bool        `json:"is_bot"`
	Has2FA   bool        `json:"has_2fa"`
	HasFiles bool        `json:"has_files"`
	Profile  UserProfile `json:"profile"`
	Presence string      `json:"presence"`
}

type UserProfile struct {
	AvatarHash string `json:"avatar_hash"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	RealName   string `json:"real_name"`
	Email      string `json:"email"`
	Skype      string `json:"skype"`
	Phone      string `json:"phone"`
	//TODO: image?

}

func GetUsers(token string) ([]Member, error) {
	userURL := "https://slack.com/api/users.list?token=%s&presence=1&pretty=0"
	var members []Member
	listUsersURL := fmt.Sprintf(userURL, token)
	resp, err := http.Get(listUsersURL)
	if err != nil {
		err = fmt.Errorf("Failed getting userlist; ", err)
		return members, err
	}

	decoder := json.NewDecoder(resp.Body)
	response := &memberResponse{}
	err = decoder.Decode(&response)
	if err != nil {
		err = fmt.Errorf("Slack decoding error: ", err)
		return members, err
	}
	if !response.Ok {
		err = fmt.Errorf("User response not ok: ", response.Error)
		return members, err
	}
	return response.Members, nil
}
