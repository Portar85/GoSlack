package rtm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

type responseRTMStart struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	Url   string       `json:"url"`
	Self  responseSelf `json:"self"`
}

type responseSelf struct {
	Id string `json:"id"`
}

type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

//Connects to slack RTM, returns url, id, error
func startRTM(token string) (string, string, error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with code %d", resp.StatusCode)
		return "", "", err
	}

	decoder := json.NewDecoder(resp.Body)
	rtmStartResponse := &responseRTMStart{}
	err = decoder.Decode(&rtmStartResponse)
	if err != nil {
		//TODO: :(
		err = fmt.Errorf("Failed decoding json for Slack rtm")
		return "", "", err
	}

	if !rtmStartResponse.Ok {
		err = fmt.Errorf("Slack error: %s", rtmStartResponse.Error)
		return "", "", err
	}

	wsURL := rtmStartResponse.Url
	id := rtmStartResponse.Self.Id
	return wsURL, id, nil
}

func GetMessage(ws *websocket.Conn) (Message, error) {
	var m Message
	err := websocket.JSON.Receive(ws, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

var counter uint64

func SendMessage(ws *websocket.Conn, m Message) error {
	m.Id = atomic.AddUint64(&counter, 1)
	return websocket.JSON.Send(ws, m)
}

func Connect(token string) (*websocket.Conn, string) {
	wsURL, id, err := startRTM(token)
	if err != nil {
		log.Fatal(err)
	}

	ws, err := websocket.Dial(wsURL, "", "https://api.slack.com")
	if err != nil {
		log.Fatal(err)
	}

	return ws, id
}
