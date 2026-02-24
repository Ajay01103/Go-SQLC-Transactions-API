package users

import (
	"net/http"

	"github.com/Ajay01103/goTransactonsAPI/internal/auth"
	jsonutil "github.com/Ajay01103/goTransactonsAPI/internal/json"
)

// Handler holds all HTTP handlers for the users domain.
type Handler struct {
	service Service
}

// NewHandler constructs a Handler with the given users Service.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetCurrentUser handles GET /users/current-user.
// Requires a valid Bearer JWT set by auth.RequireAuth middleware.
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		jsonutil.Write(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	user, err := h.service.GetCurrentUser(r.Context(), userID)
	if err != nil {
		jsonutil.Write(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}

	jsonutil.Write(w, http.StatusOK, user)
}
