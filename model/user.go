package model

// /!\ an ID in bbolt is a uint64, which might not be the case in another database.

// User represents a simple user with its followers. It also contains the JWT token.
type (
	User struct {
		ID        int    `storm:"id,increment,index"`
		Email     string `storm:"unique"`
		Password  string
		Token     string
		Followers []string
	}
)
