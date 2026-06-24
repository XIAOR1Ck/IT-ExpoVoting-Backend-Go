package voting

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	projectList, err := h.svc.store.ListProjects()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch projects",})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(projectList)
}
