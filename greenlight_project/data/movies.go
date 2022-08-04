package data

import "time"

type Movie struct {
	ID int64 `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title string `json:"title"`
	Year int32 `json:"year"`
	Runtime int32 `json:"runtime,omitempty,string"`
	Type []string `json:"type"`
	Version int `json:"version"`
}
