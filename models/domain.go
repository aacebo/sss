package models

import (
	"encoding/json"
	"time"
)

type Domain struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Extension string    `json:"extension"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (self Domain) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
