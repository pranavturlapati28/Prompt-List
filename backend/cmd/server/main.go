package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/pranavturlapati28/merget-takehome/internal/api"
	"github.com/pranavturlapati28/merget-takehome/internal/config"
	"github.com/pranavturlapati28/merget-takehome/internal/database"
	"github.com/pranavturlapati28/merget-takehome/internal/repository"
	"github.com/pranavturlapati28/merget-takehome/internal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	err = database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	err = database.Migrate()
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	err = database.Seed()
	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	repo := repository.NewPromptRepository()
	service := services.NewPromptService(repo)
	handler := api.NewHandler(service)

	router := chi.NewMux()

	router.Use(middleware.Recoverer)
	router.Use(corsMiddleware)
	router.Use(apiKeyMiddleware(cfg))
	router.Use(middleware.Logger)

	humaConfig := huma.DefaultConfig("Prompt Tree API", "1.0.0")
	humaConfig.Info.Description = "API for exploring and annotating hierarchical prompt trees"
	humaConfig.Info.Contact = &huma.Contact{
		Name:  "Pranav Turlapati",
		Email: "pranav@example.com",
	}

	humaAPI := humachi.New(router, humaConfig)
	api.RegisterRoutes(humaAPI, handler)

	printStartupBanner(cfg.Port, cfg)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/events" || r.URL.Path == "/test-sse" {
			next.ServeHTTP(w, r)
			return
		}
		
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"http://localhost:5173",
			"http://localhost:3000",
			"https://frontend-709459926380.us-central1.run.app",
		}
		
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func apiKeyMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := []string{
				"http://localhost:5173",
				"http://localhost:3000",
				"https://frontend-709459926380.us-central1.run.app",
			}
			
			origin := r.Header.Get("Origin")
			isAllowedOrigin := false
			
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					isAllowedOrigin = true
					break
				}
			}
			
			if cfg.APIKey == "" {
				next.ServeHTTP(w, r)
				return
			}
			
			if isAllowedOrigin {
				next.ServeHTTP(w, r)
				return
			}
			
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error":"API key required. Use Authorization: Bearer <your-api-key>"}`)
				return
			}
			
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error":"Invalid authorization format. Use Authorization: Bearer <your-api-key>"}`)
				return
			}
			
			if parts[1] != cfg.APIKey {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error":"Invalid API key"}`)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

func printStartupBanner(port string, cfg *config.Config) {
	fmt.Println("")
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              Prompt Tree API Server                        ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  Server:      http://localhost:%s                            ║\n", port)
	fmt.Printf("║  API Docs:    http://localhost:%s/docs                       ║\n", port)
	if cfg.APIKey != "" {
		fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
		fmt.Println("║  Security:                                                   ║")
		fmt.Println("║    API Key:     Required for external requests              ║")
		fmt.Println("║    Frontend:    Whitelisted (no API key needed)              ║")
		fmt.Println("║    Usage:       Authorization: Bearer <your-api-key>         ║")
	}
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Println("║  Endpoints:                                                   ║")
	fmt.Println("║    GET    /health              Health check                   ║")
	fmt.Println("║    GET    /tree                Full prompt tree               ║")
	fmt.Println("║    GET    /tree/export         Export tree as JSON            ║")
	fmt.Println("║    POST   /tree/import         Import tree from JSON         ║")
	fmt.Println("║    POST   /tree/save           Save current tree              ║")
	fmt.Println("║    GET    /tree/saves          List saved trees               ║")
	fmt.Println("║    POST   /tree/load/{name}    Load saved tree                ║")
	fmt.Println("║    DELETE /tree/saves/{name}   Delete saved tree              ║")
	fmt.Println("║    GET    /prompts/{id}        Single prompt                  ║")
	fmt.Println("║    POST   /prompts/{id}        Create prompt                  ║")
	fmt.Println("║    PUT    /prompts/{id}        Update prompt                  ║")
	fmt.Println("║    DELETE /prompts/{id}        Delete prompt                  ║")
	fmt.Println("║    GET    /prompts/{id}/nodes  Get nodes                      ║")
	fmt.Println("║    POST   /prompts/{id}/nodes  Create node                    ║")
	fmt.Println("║    PUT    /prompts/{id}/nodes/{nodeId} Update node           ║")
	fmt.Println("║    DELETE /prompts/{id}/nodes/{nodeId} Delete node            ║")
	fmt.Println("║    GET    /prompts/{id}/notes  Get notes                      ║")
	fmt.Println("║    POST   /prompts/{id}/notes  Create note                    ║")
	fmt.Println("║    PUT    /prompts/{id}/notes/{noteId} Update note           ║")
	fmt.Println("║    DELETE /prompts/{id}/notes/{noteId} Delete note            ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println("")
}