package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MBeliou/go-api-template/model"
	"github.com/asdine/storm/q"
	"github.com/labstack/echo"
)

func (h *Handler) CreatePost(c echo.Context) error {
	u := &model.User{
		ID: userIDFromToken(c),
	}

	p := &model.Post{
		From: u.ID,
	}

	if err := c.Bind(p); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: fmt.Sprintf("post invalid - %s", err)}
	}

	// Validate
	if p.To == 0 || p.Message == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid to or message fields"}
	}

	// Find user from db
	if err := h.DB.One("ID", u.ID, u); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: fmt.Sprintf("user not found - %s", err)}
	}

	err := h.DB.Save(p)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: fmt.Sprintf("couldn't create post - %s", err)}
	}

	return c.JSON(http.StatusCreated, p)

}

func (h *Handler) FetchPosts(c echo.Context) error {
	userID := userIDFromToken(c)
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	// Defaults
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 100
	}

	// Retrieve from DB
	posts := []*model.Post{}

	if err := h.DB.Select(q.Eq("To", userID)).Limit(limit).Skip((page - 1) * limit).Find(&posts); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: fmt.Sprintf("couldn't get posts - %s", err)}
	}

	return c.JSON(http.StatusOK, posts)
}

func (h *Handler) FetchAllPosts(c echo.Context) error {
	/*
		userID := userIDFromToken(c)
		page, _ := strconv.Atoi(c.QueryParam("page"))
		limit, _ := strconv.Atoi(c.QueryParam("limit"))

		// Defaults
		if page == 0 {
			page = 1
		}

		if limit == 0 {
			limit = 100
		}
	*/

	// Retrieve from DB
	posts := []*model.Post{}

	h.DB.All(&posts)

	return c.JSON(http.StatusOK, posts)
}
