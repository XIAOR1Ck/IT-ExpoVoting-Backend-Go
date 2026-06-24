package voting

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
    // Parse JSON body
    var req Project
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
        return
    }

    // Validate required fields
    if req.Name == "" || req.Members == "" || req.Organization == "" {
        http.Error(w, "Missing required fields: name, members, organization", http.StatusBadRequest)
        return
    }

    // Create Project struct with auto-generated ID
    project := &Project{
        Id:           uuid.New().String(),
        Name:         req.Name,
        Members:      req.Members,
        Organization: req.Organization,
        Description:  req.Description,
        Votes:        0, // Initialize votes to 0
    }

    // Store the project
    if err := h.svc.store.AddProjects(project); err != nil {
        http.Error(w, "Failed to create project: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Return success response with the created project
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(project)
}
