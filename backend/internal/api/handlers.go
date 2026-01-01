package api

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/pranavturlapati28/merget-takehome/internal/models"
	"github.com/pranavturlapati28/merget-takehome/internal/services"
)

// Handler contains all HTTP handlers
// It depends only on the service layer
type Handler struct {
	service *services.PromptService
}

// NewHandler creates a new handler with its dependencies
func NewHandler(service *services.PromptService) *Handler {
	return &Handler{service: service}
}

// =============================================================================
// HUMA REQUEST/RESPONSE TYPES
// Huma uses these structs to generate OpenAPI documentation and validate input
// =============================================================================

// --- Health ---

type HealthOutput struct {
	Body struct {
		Status  string `json:"status" example:"ok" doc:"Health status"`
		Message string `json:"message" example:"Prompt Tree API is running" doc:"Status message"`
	}
}

// --- Tree ---

type TreeOutput struct {
	Body models.TreeResponse
}

// --- Prompt ---

type PromptPathParams struct {
	ID int `path:"id" minimum:"1" doc:"Prompt ID"`
}

type GetPromptOutput struct {
	Body models.PromptDetail
}

type CreatePromptInput struct {
	ID   int `path:"id" doc:"Ignored for creation"`
	Body models.CreatePromptRequest
}

type CreatePromptOutput struct {
	Body models.Prompt
}

// --- Nodes ---

type GetNodesOutput struct {
	Body []models.Node
}

type CreateNodeInput struct {
	ID   int `path:"id" minimum:"1" doc:"Prompt ID to add node to"`
	Body models.CreateNodeRequest
}

type CreateNodeOutput struct {
	Body models.Node
}

// --- Notes ---

type GetNotesOutput struct {
	Body []models.Note
}

type CreateNoteInput struct {
	ID   int `path:"id" minimum:"1" doc:"Prompt ID to add note to"`
	Body models.CreateNoteRequest
}

type CreateNoteOutput struct {
	Body models.Note
}

// --- Update/Delete Input/Output Types ---

type UpdatePromptInput struct {
	ID   int `path:"id" minimum:"1" doc:"Prompt ID"`
	Body models.UpdatePromptRequest
}

type UpdatePromptOutput struct {
	Body models.Prompt
}

type NodePathParams struct {
	ID     int `path:"id" minimum:"1" doc:"Prompt ID"`
	NodeID int `path:"nodeId" minimum:"1" doc:"Node ID"`
}

type UpdateNodeInput struct {
	ID     int `path:"id" minimum:"1" doc:"Prompt ID"`
	NodeID int `path:"nodeId" minimum:"1" doc:"Node ID"`
	Body   models.UpdateNodeRequest
}

type UpdateNodeOutput struct {
	Body models.Node
}

type NotePathParams struct {
	ID     int `path:"id" minimum:"1" doc:"Prompt ID"`
	NoteID int `path:"noteId" minimum:"1" doc:"Note ID"`
}

type UpdateNoteInput struct {
	ID     int `path:"id" minimum:"1" doc:"Prompt ID"`
	NoteID int `path:"noteId" minimum:"1" doc:"Note ID"`
	Body   models.UpdateNoteRequest
}

type UpdateNoteOutput struct {
	Body models.Note
}

// =============================================================================
// HANDLERS
// Each handler is a thin wrapper that:
// 1. Extracts input from the request
// 2. Calls the service
// 3. Returns the response or error
// =============================================================================

// Health returns the API health status
func (h *Handler) Health(ctx context.Context, input *struct{}) (*HealthOutput, error) {
	resp := &HealthOutput{}
	resp.Body.Status = "ok"
	resp.Body.Message = "Prompt Tree API is running"
	return resp, nil
}

// GetTree returns the full prompt tree
func (h *Handler) GetTree(ctx context.Context, input *struct{}) (*TreeOutput, error) {
	tree, err := h.service.GetTree()
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to fetch tree", err)
	}
	return &TreeOutput{Body: *tree}, nil
}

// GetPrompt returns a single prompt by ID
func (h *Handler) GetPrompt(ctx context.Context, input *PromptPathParams) (*GetPromptOutput, error) {
	prompt, err := h.service.GetPrompt(input.ID)

	if errors.Is(err, services.ErrPromptNotFound) {
		return nil, huma.Error404NotFound("Prompt not found")
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to fetch prompt", err)
	}

	return &GetPromptOutput{Body: *prompt}, nil
}

// CreatePrompt creates a new prompt
func (h *Handler) CreatePrompt(ctx context.Context, input *CreatePromptInput) (*CreatePromptOutput, error) {
	prompt, err := h.service.CreatePrompt(input.Body.Title, input.Body.Description)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create prompt", err)
	}
	return &CreatePromptOutput{Body: *prompt}, nil
}

// GetPromptNodes returns all nodes for a prompt
func (h *Handler) GetPromptNodes(ctx context.Context, input *PromptPathParams) (*GetNodesOutput, error) {
	nodes, err := h.service.GetPromptNodes(input.ID)

	if errors.Is(err, services.ErrPromptNotFound) {
		return nil, huma.Error404NotFound("Prompt not found")
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to fetch nodes", err)
	}

	return &GetNodesOutput{Body: nodes}, nil
}

// CreateNode creates a new node for a prompt
func (h *Handler) CreateNode(ctx context.Context, input *CreateNodeInput) (*CreateNodeOutput, error) {
	node, err := h.service.CreateNode(input.ID, input.Body.Name, input.Body.Action)

	if errors.Is(err, services.ErrPromptNotFound) {
		return nil, huma.Error404NotFound("Prompt not found")
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create node", err)
	}

	return &CreateNodeOutput{Body: *node}, nil
}

// GetNotes returns all notes for a prompt
func (h *Handler) GetNotes(ctx context.Context, input *PromptPathParams) (*GetNotesOutput, error) {
	notes, err := h.service.GetNotes(input.ID)

	if errors.Is(err, services.ErrPromptNotFound) {
		return nil, huma.Error404NotFound("Prompt not found")
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to fetch notes", err)
	}

	return &GetNotesOutput{Body: notes}, nil
}

// CreateNote creates a new note for a prompt
func (h *Handler) CreateNote(ctx context.Context, input *CreateNoteInput) (*CreateNoteOutput, error) {
	note, err := h.service.CreateNote(input.ID, input.Body.Content)

	if errors.Is(err, services.ErrPromptNotFound) {
		return nil, huma.Error404NotFound("Prompt not found")
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create note", err)
	}

	return &CreateNoteOutput{Body: *note}, nil
}

// UpdatePrompt updates an existing prompt
func (h *Handler) UpdatePrompt(ctx context.Context, input *UpdatePromptInput) (*UpdatePromptOutput, error) {
	prompt, err := h.service.UpdatePrompt(input.ID, input.Body.Title, input.Body.Description)

	if errors.Is(err, services.ErrPromptNotFound) {
		return nil, huma.Error404NotFound("Prompt not found")
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update prompt", err)
	}

	return &UpdatePromptOutput{Body: *prompt}, nil
}

// DeletePrompt deletes a prompt by ID
func (h *Handler) DeletePrompt(ctx context.Context, input *PromptPathParams) (*struct{}, error) {
	err := h.service.DeletePrompt(input.ID)

	if errors.Is(err, services.ErrPromptNotFound) {
		return nil, huma.Error404NotFound("Prompt not found")
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to delete prompt", err)
	}

	return &struct{}{}, nil
}

// UpdateNode updates an existing node
func (h *Handler) UpdateNode(ctx context.Context, input *UpdateNodeInput) (*UpdateNodeOutput, error) {
	node, err := h.service.UpdateNode(input.NodeID, input.Body.Name, input.Body.Action)

	if errors.Is(err, services.ErrNodeNotFound) {
		return nil, huma.Error404NotFound("Node not found")
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update node", err)
	}

	return &UpdateNodeOutput{Body: *node}, nil
}

// DeleteNode deletes a node by ID
func (h *Handler) DeleteNode(ctx context.Context, input *NodePathParams) (*struct{}, error) {
	err := h.service.DeleteNode(input.NodeID)

	if errors.Is(err, services.ErrNodeNotFound) {
		return nil, huma.Error404NotFound("Node not found")
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to delete node", err)
	}

	return &struct{}{}, nil
}

// UpdateNote updates an existing note
func (h *Handler) UpdateNote(ctx context.Context, input *UpdateNoteInput) (*UpdateNoteOutput, error) {
	note, err := h.service.UpdateNote(input.NoteID, input.Body.Content)

	if errors.Is(err, services.ErrNoteNotFound) {
		return nil, huma.Error404NotFound("Note not found")
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update note", err)
	}

	return &UpdateNoteOutput{Body: *note}, nil
}

// DeleteNote deletes a note by ID
func (h *Handler) DeleteNote(ctx context.Context, input *NotePathParams) (*struct{}, error) {
	err := h.service.DeleteNote(input.NoteID)

	if errors.Is(err, services.ErrNoteNotFound) {
		return nil, huma.Error404NotFound("Note not found")
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to delete note", err)
	}

	return &struct{}{}, nil
}

// =============================================================================
// TREE IMPORT/EXPORT/SAVE/LOAD HANDLERS
// =============================================================================

// ImportTreeInput is the input for POST /tree/import
type ImportTreeInput struct {
	Body models.ImportTreeRequest
}

// ImportTreeOutput indicates success
type ImportTreeOutput struct {
	Body struct {
		Message string `json:"message" example:"Tree imported successfully"`
	}
}

// ExportTreeOutput returns the current tree as JSON
type ExportTreeOutput struct {
	Body models.TreeResponse
}

// SaveTreeInput is the input for POST /tree/save
type SaveTreeInput struct {
	Body models.SaveTreeRequest
}

// SaveTreeOutput indicates success
type SaveTreeOutput struct {
	Body struct {
		Message string `json:"message" example:"Tree saved successfully"`
	}
}

// SavedTreeListOutput returns list of saved trees
type SavedTreeListOutput struct {
	Body models.SavedTreeListResponse
}

// LoadTreePathParams is the path params for POST /tree/load/{name}
type LoadTreePathParams struct {
	Name string `path:"name" doc:"Name of the saved tree"`
}

// LoadTreeOutput indicates success
type LoadTreeOutput struct {
	Body struct {
		Message string `json:"message" example:"Tree loaded successfully"`
	}
}

// DeleteSavedTreePathParams is the path params for DELETE /tree/saves/{name}
type DeleteSavedTreePathParams struct {
	Name string `path:"name" doc:"Name of the saved tree"`
}

// ImportTree imports a tree from JSON and replaces the current tree
func (h *Handler) ImportTree(ctx context.Context, input *ImportTreeInput) (*ImportTreeOutput, error) {
	err := h.service.ImportTree(&input.Body.Tree)
	if err != nil {
		return nil, huma.Error400BadRequest("Failed to import tree", err)
	}

	resp := &ImportTreeOutput{}
	resp.Body.Message = "Tree imported successfully"
	return resp, nil
}

// ExportTree returns the current tree as JSON
func (h *Handler) ExportTree(ctx context.Context, input *struct{}) (*ExportTreeOutput, error) {
	tree, err := h.service.GetTree()
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to export tree", err)
	}
	return &ExportTreeOutput{Body: *tree}, nil
}

// SaveTree saves the current tree with a name
func (h *Handler) SaveTree(ctx context.Context, input *SaveTreeInput) (*SaveTreeOutput, error) {
	err := h.service.SaveTree(input.Body.Name)
	if err != nil {
		return nil, huma.Error400BadRequest("Failed to save tree", err)
	}

	resp := &SaveTreeOutput{}
	resp.Body.Message = "Tree saved successfully"
	return resp, nil
}

// ListSavedTrees returns all saved tree names
func (h *Handler) ListSavedTrees(ctx context.Context, input *struct{}) (*SavedTreeListOutput, error) {
	trees, err := h.service.ListSavedTrees()
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to list saved trees", err)
	}

	return &SavedTreeListOutput{
		Body: models.SavedTreeListResponse{Trees: trees},
	}, nil
}

// LoadTree loads a saved tree and replaces the current tree
func (h *Handler) LoadTree(ctx context.Context, input *LoadTreePathParams) (*LoadTreeOutput, error) {
	err := h.service.LoadTree(input.Name)
	if err != nil {
		if err.Error() == "saved tree not found" {
			return nil, huma.Error404NotFound("Saved tree not found")
		}
		return nil, huma.Error400BadRequest("Failed to load tree", err)
	}

	resp := &LoadTreeOutput{}
	resp.Body.Message = "Tree loaded successfully"
	return resp, nil
}

// DeleteSavedTree deletes a saved tree by name
func (h *Handler) DeleteSavedTree(ctx context.Context, input *DeleteSavedTreePathParams) (*struct{}, error) {
	err := h.service.DeleteSavedTree(input.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Saved tree not found")
		}
		return nil, huma.Error500InternalServerError("Failed to delete saved tree", err)
	}

	return &struct{}{}, nil
}