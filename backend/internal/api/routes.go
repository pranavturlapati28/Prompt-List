package api

import (
	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes sets up all API routes with Huma
// Huma automatically generates OpenAPI documentation from these definitions
func RegisterRoutes(api huma.API, handler *Handler) {

	// Health check endpoint
	huma.Register(api, huma.Operation{
		OperationID: "health",
		Method:      "GET",
		Path:        "/health",
		Summary:     "Health Check",
		Description: "Returns the health status of the API",
		Tags:        []string{"Health"},
	}, handler.Health)

	// Get full tree
	huma.Register(api, huma.Operation{
		OperationID: "getTree",
		Method:      "GET",
		Path:        "/tree",
		Summary:     "Get Prompt Tree",
		Description: "Returns the complete prompt tree with all prompts and their nodes for visualization",
		Tags:        []string{"Tree"},
	}, handler.GetTree)

	// Get single prompt
	huma.Register(api, huma.Operation{
		OperationID: "getPrompt",
		Method:      "GET",
		Path:        "/prompts/{id}",
		Summary:     "Get Prompt",
		Description: "Returns a single prompt by its ID",
		Tags:        []string{"Prompts"},
	}, handler.GetPrompt)

	// Create prompt
	huma.Register(api, huma.Operation{
		OperationID:   "createPrompt",
		Method:        "POST",
		Path:          "/prompts/{id}",
		Summary:       "Create Prompt",
		Description:   "Creates a new prompt. The ID in the path is ignored.",
		Tags:          []string{"Prompts"},
		DefaultStatus: 201,
	}, handler.CreatePrompt)

	// Get nodes for a prompt
	huma.Register(api, huma.Operation{
		OperationID: "getPromptNodes",
		Method:      "GET",
		Path:        "/prompts/{id}/nodes",
		Summary:     "Get Prompt Nodes",
		Description: "Returns all nodes (subprompts) for a specific prompt",
		Tags:        []string{"Nodes"},
	}, handler.GetPromptNodes)

	// Create node for a prompt
	huma.Register(api, huma.Operation{
		OperationID:   "createNode",
		Method:        "POST",
		Path:          "/prompts/{id}/nodes",
		Summary:       "Create Node",
		Description:   "Creates a new node (subprompt) for a specific prompt",
		Tags:          []string{"Nodes"},
		DefaultStatus: 201,
	}, handler.CreateNode)

	// Get notes for a prompt
	huma.Register(api, huma.Operation{
		OperationID: "getNotes",
		Method:      "GET",
		Path:        "/prompts/{id}/notes",
		Summary:     "Get Notes",
		Description: "Returns all user annotations for a specific prompt",
		Tags:        []string{"Notes"},
	}, handler.GetNotes)

	// Create note for a prompt
	huma.Register(api, huma.Operation{
		OperationID:   "createNote",
		Method:        "POST",
		Path:          "/prompts/{id}/notes",
		Summary:       "Create Note",
		Description:   "Creates a new annotation for a specific prompt",
		Tags:          []string{"Notes"},
		DefaultStatus: 201,
	}, handler.CreateNote)
}