package api

import (
	"context"
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