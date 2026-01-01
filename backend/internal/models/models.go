package models

import "time"

// =============================================================================
// DATABASE MODELS
// These represent how data is stored in the database
// =============================================================================

// Prompt represents a main prompt in the tree (e.g., "Project Setup")
type Prompt struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProjectName string `json:"project_name,omitempty"`
}

// Node represents a subprompt/step under a prompt (e.g., "npm create vite")
type Node struct {
	ID       int    `json:"id,omitempty"`
	PromptID int    `json:"prompt_id,omitempty"`
	Name     string `json:"name"`
	Action   string `json:"action"`
}

// Note represents a user's annotation on a prompt
type Note struct {
	ID        int       `json:"id"`
	PromptID  int       `json:"prompt_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// SavedTree represents a saved prompt tree configuration
type SavedTree struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	TreeData  string    `json:"tree_data"` // JSON string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// =============================================================================
// API RESPONSE MODELS
// These represent how data is sent to the frontend
// =============================================================================

// TreeResponse is returned by GET /tree
// It contains the full tree structure for visualization
type TreeResponse struct {
	Project     string       `json:"project" doc:"Project name"`
	MainRequest string       `json:"mainRequest" doc:"Main project description"`
	Prompts     []PromptNode `json:"prompts" doc:"List of prompts with their nodes"`
}

// PromptNode represents a prompt in the tree response
type PromptNode struct {
	ID          int           `json:"id" doc:"Prompt ID"`
	Title       string        `json:"title" doc:"Prompt title"`
	Description string        `json:"description,omitempty" doc:"Prompt description"`
	Nodes       []NodeSummary `json:"nodes,omitempty" doc:"Child nodes of this prompt"`
}

// NodeSummary is a simplified node for the tree view
type NodeSummary struct {
	ID     int    `json:"id" doc:"Node ID"`
	Name   string `json:"name" doc:"Node name"`
	Action string `json:"action,omitempty" doc:"Node action description"`
}

// PromptDetail is returned by GET /prompts/:id
type PromptDetail struct {
	ID          int    `json:"id" doc:"Prompt ID"`
	Title       string `json:"title" doc:"Prompt title"`
	Description string `json:"description" doc:"Prompt description"`
	ProjectName string `json:"project_name" doc:"Parent project name"`
}

// =============================================================================
// API REQUEST MODELS
// These represent data sent from the frontend
// =============================================================================

// CreatePromptRequest is the body for POST /prompts/:id
type CreatePromptRequest struct {
	Title       string `json:"title" minLength:"1" doc:"Title of the prompt (required)"`
	Description string `json:"description,omitempty" doc:"Description of the prompt"`
}

// CreateNodeRequest is the body for POST /prompts/:id/nodes
type CreateNodeRequest struct {
	Name   string `json:"name" minLength:"1" doc:"Name of the node (required)"`
	Action string `json:"action,omitempty" doc:"Action description"`
}

// CreateNoteRequest is the body for POST /prompts/:id/notes
type CreateNoteRequest struct {
	Content string `json:"content" minLength:"1" doc:"Note content (required)"`
}

// UpdatePromptRequest is the body for PUT /prompts/:id
type UpdatePromptRequest struct {
	Title       string `json:"title,omitempty" doc:"Title of the prompt"`
	Description string `json:"description,omitempty" doc:"Description of the prompt"`
}

// UpdateNodeRequest is the body for PUT /prompts/:id/nodes/:nodeId
type UpdateNodeRequest struct {
	Name   string `json:"name,omitempty" doc:"Name of the node"`
	Action string `json:"action,omitempty" doc:"Action description"`
}

// UpdateNoteRequest is the body for PUT /prompts/:id/notes/:noteId
type UpdateNoteRequest struct {
	Content string `json:"content" minLength:"1" doc:"Note content (required)"`
}

// ImportTreeRequest is the body for POST /tree/import
type ImportTreeRequest struct {
	Tree TreeResponse `json:"tree" doc:"Complete tree structure to import"`
}

// SaveTreeRequest is the body for POST /tree/save
type SaveTreeRequest struct {
	Name string `json:"name" minLength:"1" doc:"Name for the saved tree (required)"`
}

// SavedTreeInfo represents metadata about a saved tree
type SavedTreeInfo struct {
	Name      string    `json:"name" doc:"Name of the saved tree"`
	CreatedAt time.Time `json:"created_at" doc:"When the tree was saved"`
	UpdatedAt time.Time `json:"updated_at" doc:"When the tree was last updated"`
}

// SavedTreeListResponse contains a list of saved trees
type SavedTreeListResponse struct {
	Trees []SavedTreeInfo `json:"trees" doc:"List of saved trees"`
}