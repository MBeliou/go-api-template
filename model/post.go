package model

import "time"

type (
	Post struct {
		ID        uint64 `storm:"id,increment"`
		To        int    `storm:"unique"`
		From      int    `storm:"unique"`
		Message   string
		CreatedAt time.Time `storm:"index"`
	}
)
