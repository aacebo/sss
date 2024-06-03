package models

import (
	"encoding/json"
	"time"
)

type Link struct {
	FromID    string    `json:"from_id"`
	ToID      string    `json:"to_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (self Link) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
