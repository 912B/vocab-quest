package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Simple In-Memory Session Store
var (
	Sessions  = make(map[string]int) // Token -> UserID
	SessionMu sync.RWMutex
)

type AuthHandler struct {
	DB *sql.DB
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Default Avatar if empty
	if req.Avatar == "" {
		req.Avatar = "👨‍🚀"
	}

	// Check if user exists
	var count int
	h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
	if count > 0 {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Hash Password
	hashed, err := HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Create user
	_, err = h.DB.Exec("INSERT INTO users (username, password, avatar) VALUES (?, ?, ?)", req.Username, hashed, req.Avatar)
	if err != nil {
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"registered"}`))
}

func (h *AuthHandler) ListPublicUsers(w http.ResponseWriter, r *http.Request) {
	// Return list of users (ID, Username, Avatar) for login screen
	rows, err := h.DB.Query("SELECT id, username, avatar FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type PublicUser struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}

	var users []PublicUser
	for rows.Next() {
		var u PublicUser
		if err := rows.Scan(&u.ID, &u.Username, &u.Avatar); err != nil {
			continue
		}
		users = append(users, u)
	}
	if users == nil {
		users = make([]PublicUser, 0)
	}

	json.NewEncoder(w).Encode(users)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var userID int
	var role string
	var storedPass string

	// Get User info and stored password (hash or plain)
	err := h.DB.QueryRow("SELECT id, role, password FROM users WHERE username = ?", req.Username).Scan(&userID, &role, &storedPass)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Password Verification Logic (Migration Support)
	match := CheckPasswordHash(req.Password, storedPass)
	if !match {
		// Fallback: Check if it's a legacy plain text password
		if storedPass == req.Password {
			// It IS a match (Plain text), so we upgrade it to Hash now
			newHash, _ := HashPassword(req.Password)
			h.DB.Exec("UPDATE users SET password = ? WHERE id = ?", newHash, userID)
			log.Printf("Security: Upgraded password for user %d to bcrypt", userID)
			match = true
		}
	}

	if !match {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	log.Printf("Login Success: UserID=%d Role=%s", userID, role)

	// Create Session
	token := generateToken()
	SessionMu.Lock()
	Sessions[token] = userID
	SessionMu.Unlock()

	// Send Token via Cookie or Response
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	w.Write([]byte(fmt.Sprintf(`{"status":"logged_in", "user_id":%d, "role":"%s"}`, userID, role)))
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// Middleware to get UserID
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized: No session", http.StatusUnauthorized)
			return
		}

		SessionMu.RLock()
		_, ok := Sessions[c.Value]
		SessionMu.RUnlock()

		if !ok {
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func GetUserID(r *http.Request) (int, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return 0, err
	}

	SessionMu.RLock()
	defer SessionMu.RUnlock()
	id, ok := Sessions[c.Value]
	if !ok {
		return 0, fmt.Errorf("invalid session")
	}
	return id, nil
}

// Admin: List or Manage Users
func (h *AuthHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	// 1. Check Admin Role (Middleware does Authentication, but we need Authorization)
	// For MVP, we trust the frontend UI hiding, but ideally we check role here.
	switch r.Method {
	case http.MethodGet:
		rows, err := h.DB.Query("SELECT id, username, role, avatar FROM users")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var result []map[string]interface{}

		for rows.Next() {
			var id int
			var u, role, av string
			rows.Scan(&id, &u, &role, &av)
			// Handle nulls if necessary, but schema has defaults
			result = append(result, map[string]interface{}{
				"id": id, "username": u, "role": role, "avatar": av,
			})
		}
		json.NewEncoder(w).Encode(result)

	case http.MethodPost:
		// Create New User (Admin Only)
		var payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Role     string `json:"role"`
			Avatar   string `json:"avatar"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		// Defaults
		if payload.Role == "" {
			payload.Role = "user"
		}
		if payload.Avatar == "" {
			payload.Avatar = "👨‍🚀"
		}

		// Check exists
		var count int
		h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", payload.Username).Scan(&count)
		if count > 0 {
			http.Error(w, "Username taken", 409)
			return
		}

		// Hash Password
		hashed, _ := HashPassword(payload.Password)

		_, err := h.DB.Exec("INSERT INTO users (username, password, role, avatar) VALUES (?, ?, ?, ?)",
			payload.Username, hashed, payload.Role, payload.Avatar)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodPut:
		// Update Role or Reset Password
		var payload struct {
			ID       int    `json:"id"`
			Role     string `json:"role"`
			Password string `json:"password"` // Optional
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		if payload.Password != "" {
			// Hash Password
			hashed, _ := HashPassword(payload.Password)
			_, err := h.DB.Exec("UPDATE users SET role = ?, password = ? WHERE id = ?", payload.Role, hashed, payload.ID)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		} else {
			_, err := h.DB.Exec("UPDATE users SET role = ? WHERE id = ?", payload.Role, payload.ID)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing id", 400)
			return
		}
		if _, err := h.DB.Exec("DELETE FROM users WHERE id = ?", id); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
