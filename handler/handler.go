package handler

import (
	"github.com/asdine/storm"
	"golang.org/x/crypto/bcrypt"
)

const (
	saltCost = bcrypt.DefaultCost
	JwtKey   = "mysecretkey" // Should come from somewhere else
)

type (
	Handler struct {
		DB *storm.DB
	}
)
