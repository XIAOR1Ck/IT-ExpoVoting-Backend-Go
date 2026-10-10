package voting

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	qrcode "github.com/skip2/go-qrcode"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
	
	"encoding/json"
)

type Service struct{
	store *RedisStore
}


func NewService(rdb *RedisStore) *Service {
	return &Service{
		store: rdb,
	}
}

func (s *Service) GenerateQR(id string) (string, error) {
	URL := fmt.Sprintf("http://placeholder.com/vote/userid?=%s", id)
	outputdir := "./QrCodes"
	filename := fmt.Sprintf("%s.png", id)
		err := generateQRCode(URL, outputdir, filename)
if err != nil {
        fmt.Printf("Error generating QR code: %v\n", err)
        return "" , err
    }
    
	fileURL := filepath.Join(outputdir, filename)
	return fileURL, nil

}

func generateQRCode(content, outputDir, filename string) error {
    // Create the output directory if it doesn't exist
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return fmt.Errorf("failed to create directory: %v", err)
    }
    
    // Full path for the output file
    outputPath := filepath.Join(outputDir, filename)
    
    // Generate the QR code with custom settings
    // Parameters: content, size, quality, and output path
    err := qrcode.WriteFile(content, qrcode.Medium, 256, outputPath)
    if err != nil {
        return fmt.Errorf("failed to generate QR code: %v", err)
    }
    
    return nil
}

func (s *Service) HashPassword(password string) (string, error) {
bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

var jwtSecret = []byte("11-Expo-JWT")

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token
func (s *Service) GenerateJWT(email string) (string, error) {

	claims := Claims{
		UserID: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

// ValidateToken validates a JWT token and returns claims
func (s *Service) ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Ensure correct signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return jwtSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *Service) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := cookie.Value

		_, valerr := s.ValidateJWT(token)	
		if valerr != nil {
	w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "fetching failed",})
		return
		}

		next.ServeHTTP(w, r)
	})
}
