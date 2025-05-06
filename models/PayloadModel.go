package models


type MessagePayload struct {
	Request  Message `json:"request"`
	Response Message `json:"response"`
}
