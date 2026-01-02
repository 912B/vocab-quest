package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"vocab-quest/server/database"
	"vocab-quest/server/handlers"
	"vocab-quest/server/services"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 1. Initialize Database
	db, err := database.InitDB("./vocab.db")
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer db.Close()

	// 1.5 Load Config
	config, err := LoadConfig("server_config.json")
	if err != nil {
		log.Printf("Config: %v", err)
	} else {
		log.Printf("Config loaded. Port: %s", config.Port)
		// Process Admin Reset
		if config.AdminReset.Enabled {
			log.Println("Config: Processing Admin Reset...")
			hashed, _ := handlers.HashPassword("admin") // Default to 'admin'

			// Check if user exists
			var count int
			db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", config.AdminReset.Username).Scan(&count)

			if count == 0 {
				// Create
				_, err := db.Exec("INSERT INTO users (username, password, role, avatar) VALUES (?, ?, 'admin', '👨‍🚀')", config.AdminReset.Username, hashed)
				if err != nil {
					log.Printf("Config Reset Error: Failed to create admin: %v", err)
				} else {
					log.Println("Config Reset: Admin created.")
				}
			} else {
				// Update
				_, err := db.Exec("UPDATE users SET password = ? WHERE username = ?", hashed, config.AdminReset.Username)
				if err != nil {
					log.Printf("Config Reset Error: Failed to update admin: %v", err)
				} else {
					log.Println("Config Reset: Admin password updated.")
				}
			}

			// Disable Reset in Config
			config.AdminReset.Enabled = false
			SaveConfig("server_config.json", config)
			log.Println("Config: Reset flag disabled and file updated.")
		}
	}

	// 2. Setup API Routes
	authHandler := &handlers.AuthHandler{DB: db}
	http.HandleFunc("/api/register", authHandler.Register)
	http.HandleFunc("/api/login", authHandler.Login)
	http.HandleFunc("/api/users/public", authHandler.ListPublicUsers)
	http.HandleFunc("/api/users", handlers.Authenticate(authHandler.HandleUsers)) // New Admin Route

	gameHandler := &handlers.GameHandler{
		DB:     db,
		Engine: services.NewLearningEngine(db),
	}
	// Protected Routes (Simple manual check or Middleware wrapper if we had a router...
	// Stdlib Wrapper: handlers.Authenticate(gameHandler.GetSession))
	http.HandleFunc("/api/session", handlers.Authenticate(gameHandler.GetSession))
	http.HandleFunc("/api/result", handlers.Authenticate(gameHandler.SubmitResult))
	http.HandleFunc("/api/stats", handlers.Authenticate(gameHandler.HandleGetStats))

	// Management API
	dictHandler := &handlers.DictHandler{DB: db}
	http.HandleFunc("/api/dictionaries", handlers.Authenticate(dictHandler.List))
	http.HandleFunc("/api/dictionaries/words", handlers.Authenticate(dictHandler.ListWords))
	http.HandleFunc("/api/dictionaries/active", handlers.Authenticate(dictHandler.SetActive))
	http.HandleFunc("/api/dictionaries/import", handlers.Authenticate(dictHandler.ImportWords))
	http.HandleFunc("/api/dictionaries/template", handlers.Authenticate(dictHandler.DownloadTemplate))

	wordHandler := &handlers.WordHandler{DB: db}
	http.HandleFunc("/api/words", handlers.Authenticate(wordHandler.Handle))

	// 3. Setup Static File Server with No-Cache Middleware
	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/", NoCache(fs))

	// 3. Start Server
	port := "8081"
	if config != nil && config.Port != "" {
		port = config.Port
	}

	log.Printf("Server starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// Config Structures
type ServerConfig struct {
	Port       string           `json:"port"`
	AdminReset AdminResetConfig `json:"admin_reset"`
}

type AdminResetConfig struct {
	Enabled  bool   `json:"enabled"`
	Username string `json:"username"`
	Comment  string `json:"_comment,omitempty"`
}

func LoadConfig(filename string) (*ServerConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Config: File not found, creating default template...")
			defaultConfig := &ServerConfig{
				Port: "8081",
				AdminReset: AdminResetConfig{
					Enabled:  false,
					Username: "admin",
					Comment:  "Set enabled to true to reset password to 'admin' on restart.",
				},
			}
			if err := SaveConfig(filename, defaultConfig); err != nil {
				return nil, err
			}
			return defaultConfig, nil
		}
		return nil, err // Other errors (permission, etc)
	}
	defer file.Close()

	var config ServerConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func SaveConfig(filename string, config *ServerConfig) error {
	file, err := os.Create(filename) // Truncates
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

// NoCache Wrapper Middleware
func NoCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		h.ServeHTTP(w, r)
	})
}
