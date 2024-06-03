package models

import (
	"encoding/json"
	"time"
)

type Page struct {
	ID        string    `json:"id"`
	DomainID  string    `json:"domain_id"`
	Title     *string   `json:"title,omitempty"`
	Url       string    `json:"url"`
	Address   string    `json:"address"`
	Size      int64     `json:"size"`
	ElapseMs  int64     `json:"elapse_ms"`
	LinkCount int       `json:"link_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (self Page) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
