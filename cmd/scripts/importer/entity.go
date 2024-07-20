package main

type tInput struct {
	Name     string     `json:"name"`
	Type     string     `json:"type"`
	ID       int        `json:"id"`
	Messages []tMessage `json:"messages"`
}

type tMessage struct {
	ID           int         `json:"id"`
	Type         string      `json:"type"`
	Date         string      `json:"date"`
	DateUnixtime string      `json:"date_unixtime"`
	Actor        string      `json:"actor,omitempty"`
	ActorID      string      `json:"actor_id,omitempty"`
	Action       string      `json:"action,omitempty"`
	Title        string      `json:"title,omitempty"`
	Photo        string      `json:"photo,omitempty"`
	Text         interface{} `json:"text"`
	TextEntities []struct {
		Type string `json:"type"`
		Text string `json:"text"`
		Href string `json:"href,omitempty"`
	} `json:"text_entities"`
	Edited         string `json:"edited,omitempty"`
	EditedUnixtime string `json:"edited_unixtime,omitempty"`
	From           string `json:"from,omitempty"`
	FromID         string `json:"from_id,omitempty"`
	Author         string `json:"author,omitempty"`
}
