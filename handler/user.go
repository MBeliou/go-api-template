package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MBeliou/go-api-template/model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/forestgiant/sliceutil"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Signup(c echo.Context) error {
	u := &model.User{}
	if err := c.Bind(u); err != nil {
		return err
	}

	// Validate
	if u.Email == "" || u.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid email or password"}
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), saltCost)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: fmt.Sprintf("couldn't hash password - %s", err)}
	}

	u.Password = string(hashedPassword)

	err = h.DB.Save(u)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: fmt.Sprintf("couldn't create user - %s", err)}
	}

	return c.JSON(http.StatusCreated, u)
}

func (h *Handler) Login(c echo.Context) error {
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: fmt.Sprintf("couldn't create user - %s", err)}
	}

	var user model.User
	err := h.DB.One("Email", u.Email, &user)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: fmt.Sprintf("find user - %s", err)}
	}

	// Check password & its hash
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: fmt.Sprintf("invalid email or password - %s", err)}
	}

	// JWT
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	u.Token, err = token.SignedString([]byte(JwtKey))
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: fmt.Sprintf("signing - %s", err)}

	}

	u.Password = "" // We don't want to send the password back

	return c.JSON(http.StatusOK, u)
}

// TODO: only dev
func (h *Handler) GetAll(c echo.Context) error {
	var users []model.User
	err := h.DB.All(&users)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: fmt.Sprintf("query all - %s", err)}
	}

	return c.JSON(http.StatusOK, users)
}

func makeError(code int, message interface{}) *echo.HTTPError {
	return &echo.HTTPError{Code: code, Message: message}
}

func (h *Handler) Follow(c echo.Context) error {
	userID := userIDFromToken(c)

	id := c.Param("id")
	if id == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid user id"}

	}

	// Add a follower to user
	var user model.User

	// TODO: this should be in a transaction
	err := h.DB.One("ID", userID, &user)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: fmt.Sprintf("find user - %s", err)}
	}

	// User can't be followed twice
	alreadyFollowed := sliceutil.Contains(user.Followers, id)
	if alreadyFollowed {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "already following user"}
	}

	user.Followers = append(user.Followers, id)

	// now we can update the user
	h.DB.UpdateField(&model.User{ID: userID}, "Followers", user.Followers)

	return nil
}

func userIDFromToken(c echo.Context) int {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	fmt.Printf("Obtained claim id: %v\n", claims["id"])

	fmt.Printf("as int: %v\n", int(claims["id"].(float64)))

	return int(claims["id"].(float64))
}
