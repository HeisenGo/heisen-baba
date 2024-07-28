package handlers

import (
	presenter "authservice/api/http/handlers/presentor"
	"authservice/internal/user"
	"authservice/service"
	"errors"
	"fmt"

	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)
func RegisterUser(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.UserRegisterReq

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}
		err := BodyValidator(req)
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		u := presenter.UserRegisterToUserDomain(&req)

		newUser, err := authService.CreateUser(c.Context(), u)
		if err != nil {
			if errors.Is(err, user.ErrInvalidEmail) || errors.Is(err, user.ErrInvalidPassword) {
				return presenter.BadRequest(c, err)
			}
			if errors.Is(err, user.ErrEmailAlreadyExists) {
				return presenter.Conflict(c, err)
			}

			return presenter.InternalServerError(c, err)
		}

		return presenter.Created(c, "user successfully registered", fiber.Map{
			"user_id": newUser.ID,
		})
	}
}
func LoginUser(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req presenter.UserLoginReq

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err := BodyValidator(req)
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		c.Cookie(&fiber.Cookie{
			Name:        "X-Session-ID",
			Value:       fmt.Sprint(time.Now().UnixNano()),
			HTTPOnly:    true,
			SessionOnly: true,
		})

		authToken, err := authService.Login(c.Context(), req.Email, req.Password)
		if err != nil {

			return presenter.BadRequest(c, err)
		}
		return SendUserToken(c, authToken)
	}
}
func RefreshToken(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		refToken := c.GetReqHeaders()["Authorization"][0]
		if len(refToken) == 0 {
			return SendError(c, errors.New("token should be provided"), fiber.StatusBadRequest)
		}
		pureToken := strings.Split(refToken, " ")[1]
		authToken, err := authService.RefreshAuth(c.UserContext(), pureToken)
		if err != nil {

			return presenter.Unauthorized(c, err)
		}

		return SendUserToken(c, authToken)
	}
}


package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/your-repo/auth/internal/user"
    "github.com/your-repo/auth/pkg/jwt"
)

type AuthHandler struct {
    userService *user.Service
}

func NewAuthHandler(userService *user.Service) *AuthHandler {
    return &AuthHandler{userService: userService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user, err := h.userService.CreateUser(r.Context(), req.Username, req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user, err := h.userService.AuthenticateUser(r.Context(), req.Username, req.Password)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // Generate JWT token
    token, err := jwt.GenerateToken(user.ID, getUserRoles(user))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func getUserRoles(user *entities.User) []string {
    var roles []string
    for _, role := range user.Roles {
        roles = append(roles, role.Name)
    }
    return roles
}