package voting

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) { 

    w.Header().Set("Content-Type", "application/json")

    // Read body
    var req struct {
        Name         string `json:"name"`
        Email        string `json:"email"`
        Organization string `json:"organization"`
        Address      string `json:"address"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
        return
    }

    // Validate
    if req.Name == "" {
        http.Error(w, "Name is required", http.StatusBadRequest)
        return
    }
    if req.Email == "" {
        http.Error(w, "Email is required", http.StatusBadRequest)

        return
    }
	usrid := uuid.New().String()
	qrPath, err := h.svc.GenerateQR(usrid)
	if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(map[string]string{
        "error": err.Error(),
    })
    return  // <-- MISSING THIS!
}

    // Create user with generated ID
    user := &User{
        Id:           usrid,
        Name:         req.Name,
        Email:        req.Email,
        Organization: req.Organization,
        Address:      req.Address,
       	Qrcode: qrPath, 
    }

	newUser, err := h.svc.store.AddUser(user) 
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
    "error": err.Error(),
})
	}
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(newUser)
}

//Cast Vote Handler
func (h *Handler) CastVote(w http.ResponseWriter, r *http.Request) {
	var vote struct {
		UserID string `json:"user_id"`
		ProjectID string `json:"project_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
        http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
        return
    }
	err := h.svc.store.CastVote(vote.UserID, vote.ProjectID)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error(),})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"success": "Vote has been recorded",
	})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	//url params
	params := r.URL.Query()
	email := params.Get("email")
	if email == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "email is required",})
		return	
	}
	user, err := h.svc.store.GetUserByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error(),})
		return
	}
	w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (h *Handler) LoadBulk (w http.ResponseWriter, r *http.Request) {
	filepath := "./internal/voting/data/projects.json"
	prjList, err := LoadProjects(filepath)
	if err != nil {
	w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error(),})
		return
	}
	for _, prj := range prjList{
		err := h.svc.store.AddProjects(&prj)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error(),})
		return
		}
	}
		w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"success": "All projects added",})
}

func (h *Handler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()
	var admin Admin
	if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
        http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
        return
    }
	if err := validate.Struct(admin); err != nil {
		// validation failed
    http.Error(w, "Invalid Email Or Password", http.StatusBadRequest)
		return
	}
 password_hash, err := h.svc.HashPassword(admin.Password)
	if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Admin Creation Failed",})
		return
	}
	rdsErr := h.svc.store.CreateAdmin(admin.Email, password_hash)	
	if rdsErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": rdsErr.Error(),})
		return
	}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"success": "Admin Account Created Successfully"})

}

func (h *Handler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()
	var admin Admin
	if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
        http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
        return
    }
	if err := validate.Struct(admin); err != nil {
		// validation failed
    http.Error(w, "Invalid Email Or Password", http.StatusBadRequest)
		return
	}
	log.Printf("Email: %s",admin.Email)
 	stored_hash, fetcherr := h.svc.store.GetAdminPassword(admin.Email) 
	if fetcherr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "fetching failed",})
		return
	}
	checkerr := bcrypt.CompareHashAndPassword([]byte(stored_hash), []byte(admin.Password))
	if checkerr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid Password",})
		return
	}
	
	token, jwterr := h.svc.GenerateJWT(admin.Email)
	if jwterr != nil {
			w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Login Failed",})
		return
	}
	expiresAt := time.Now().Add(72 * time.Hour)	
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
	})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"success": "Login Successful"})

}
