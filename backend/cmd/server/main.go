package main

import (
	"fmt"
	"log"
	"net/http"

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
	// =========================================================================
	// STEP 1: Load Configuration
	// =========================================================================
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// =========================================================================
	// STEP 2: Connect to Database
	// =========================================================================
	err = database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close() // Ensure connection is closed when app exits

	// =========================================================================
	// STEP 3: Run Migrations (create tables)
	// =========================================================================
	err = database.Migrate()
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// =========================================================================
	// STEP 4: Seed Initial Data
	// =========================================================================
	err = database.Seed()
	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// =========================================================================
	// STEP 5: Initialize Application Layers (Dependency Injection)
	// =========================================================================
	// Create instances and wire dependencies
	repo := repository.NewPromptRepository()    // Data access layer
	service := services.NewPromptService(repo)  // Business logic layer
	handler := api.NewHandler(service)          // HTTP handler layer

	// =========================================================================
	// STEP 6: Create HTTP Router
	// =========================================================================
	router := chi.NewMux()

	// Add middleware
	router.Use(middleware.Recoverer) // Recover from panics
	router.Use(corsMiddleware)       // Handle CORS for frontend
	
	// Add logger middleware
	router.Use(middleware.Logger)    // Log all requests

	// =========================================================================
	// STEP 7: Create Huma API with OpenAPI Documentation
	// =========================================================================
	humaConfig := huma.DefaultConfig("Prompt Tree API", "1.0.0")
	humaConfig.Info.Description = "API for exploring and annotating hierarchical prompt trees"
	humaConfig.Info.Contact = &huma.Contact{
		Name:  "Pranav Turlapati",
		Email: "pranav@example.com",
	}

	humaAPI := humachi.New(router, humaConfig)

	// =========================================================================
	// STEP 9: Register Huma Routes
	// =========================================================================
	api.RegisterRoutes(humaAPI, handler)

	// =========================================================================
	// STEP 11: Start Server
	// =========================================================================
	printStartupBanner(cfg.Port)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}

// corsMiddleware handles Cross-Origin Resource Sharing
// This allows the frontend (on a different port) to call the API
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin (in production, be more specific)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// printStartupBanner prints a nice startup message
func printStartupBanner(port string) {
	fmt.Println("")
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              Prompt Tree API Server                        ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  Server:      http://localhost:%s                            ║\n", port)
	fmt.Printf("║  API Docs:    http://localhost:%s/docs                       ║\n", port)
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