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

	// Export tree as JSON
	huma.Register(api, huma.Operation{
		OperationID: "exportTree",
		Method:      "GET",
		Path:        "/tree/export",
		Summary:     "Export Tree",
		Description: "Returns the current prompt tree as JSON for copying/exporting",
		Tags:        []string{"Tree"},
	}, handler.ExportTree)

	// Import tree from JSON
	huma.Register(api, huma.Operation{
		OperationID:   "importTree",
		Method:        "POST",
		Path:          "/tree/import",
		Summary:       "Import Tree",
		Description:   "Imports a prompt tree from JSON and replaces the current tree",
		Tags:          []string{"Tree"},
		DefaultStatus: 201,
	}, handler.ImportTree)

	// Save current tree
	huma.Register(api, huma.Operation{
		OperationID:   "saveTree",
		Method:        "POST",
		Path:          "/tree/save",
		Summary:       "Save Tree",
		Description:   "Saves the current prompt tree with a name for later retrieval",
		Tags:          []string{"Tree"},
		DefaultStatus: 201,
	}, handler.SaveTree)

	// List saved trees
	huma.Register(api, huma.Operation{
		OperationID: "listSavedTrees",
		Method:      "GET",
		Path:        "/tree/saves",
		Summary:     "List Saved Trees",
		Description: "Returns a list of all saved tree names and metadata",
		Tags:        []string{"Tree"},
	}, handler.ListSavedTrees)

	// Load saved tree
	huma.Register(api, huma.Operation{
		OperationID:   "loadTree",
		Method:        "POST",
		Path:          "/tree/load/{name}",
		Summary:       "Load Tree",
		Description:   "Loads a saved tree and replaces the current tree",
		Tags:          []string{"Tree"},
		DefaultStatus: 201,
	}, handler.LoadTree)

	// Delete saved tree
	huma.Register(api, huma.Operation{
		OperationID: "deleteSavedTree",
		Method:      "DELETE",
		Path:        "/tree/saves/{name}",
		Summary:     "Delete Saved Tree",
		Description: "Deletes a saved tree by name",
		Tags:        []string{"Tree"},
	}, handler.DeleteSavedTree)

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

	// Update prompt
	huma.Register(api, huma.Operation{
		OperationID: "updatePrompt",
		Method:      "PUT",
		Path:        "/prompts/{id}",
		Summary:     "Update Prompt",
		Description: "Updates an existing prompt by its ID",
		Tags:        []string{"Prompts"},
	}, handler.UpdatePrompt)

	// Delete prompt
	huma.Register(api, huma.Operation{
		OperationID: "deletePrompt",
		Method:      "DELETE",
		Path:        "/prompts/{id}",
		Summary:     "Delete Prompt",
		Description: "Deletes a prompt by its ID. This will cascade delete all associated nodes and notes.",
		Tags:        []string{"Prompts"},
	}, handler.DeletePrompt)

	// Update node
	huma.Register(api, huma.Operation{
		OperationID: "updateNode",
		Method:      "PUT",
		Path:        "/prompts/{id}/nodes/{nodeId}",
		Summary:     "Update Node",
		Description: "Updates an existing node (subprompt) by its ID",
		Tags:        []string{"Nodes"},
	}, handler.UpdateNode)

	// Delete node
	huma.Register(api, huma.Operation{
		OperationID: "deleteNode",
		Method:      "DELETE",
		Path:        "/prompts/{id}/nodes/{nodeId}",
		Summary:     "Delete Node",
		Description: "Deletes a node (subprompt) by its ID",
		Tags:        []string{"Nodes"},
	}, handler.DeleteNode)

	// Update note
	huma.Register(api, huma.Operation{
		OperationID: "updateNote",
		Method:      "PUT",
		Path:        "/prompts/{id}/notes/{noteId}",
		Summary:     "Update Note",
		Description: "Updates an existing note by its ID",
		Tags:        []string{"Notes"},
	}, handler.UpdateNote)

	// Delete note
	huma.Register(api, huma.Operation{
		OperationID: "deleteNote",
		Method:      "DELETE",
		Path:        "/prompts/{id}/notes/{noteId}",
		Summary:     "Delete Note",
		Description: "Deletes a note by its ID",
		Tags:        []string{"Notes"},
	}, handler.DeleteNote)
}