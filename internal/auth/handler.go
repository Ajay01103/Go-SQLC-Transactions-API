package auth

import (
	"errors"
	"net/http"

	jsonutil "github.com/Ajay01103/goTransactonsAPI/internal/json"
)

// Handler holds all HTTP handlers for the auth domain.
type Handler struct {
	service Service
}

// NewHandler constructs a Handler with the given auth Service.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type registerRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles POST /auth/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := jsonutil.Read(r, &req); err != nil {
		jsonutil.Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if req.Email == "" || req.Password == "" {
		jsonutil.Write(w, http.StatusBadRequest, map[string]string{"error": "email and password are required"})
		return
	}

	resp, err := h.service.Login(r.Context(), LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		jsonutil.Write(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	jsonutil.Write(w, http.StatusOK, resp)
}

// Register handles POST /auth/register.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := jsonutil.Read(r, &req); err != nil {
		jsonutil.Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		jsonutil.Write(w, http.StatusBadRequest, map[string]string{"error": "name, email and password are required"})
		return
	}

	// default profile picture to empty string if not provided
	profilePicture := req.ProfilePicture

	resp, err := h.service.Register(r.Context(), RegisterInput{
		Name:           req.Name,
		Email:          req.Email,
		Password:       req.Password,
		ProfilePicture: profilePicture,
	})
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			jsonutil.Write(w, http.StatusConflict, map[string]string{"error": "an account with this email already exists"})
			return
		}
		jsonutil.Write(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
		return
	}

	jsonutil.Write(w, http.StatusCreated, resp)
}
