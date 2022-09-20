package models

type Message struct {
	ID  int64  `json:"id" example:"1"`
	Msg string `json:"msg" example:"any text"`
}
