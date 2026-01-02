package models

import "time"

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

type Note struct {
	ID        int       `json:"id"`
	PromptID  int       `json:"prompt_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type SavedTree struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	TreeData  string    `json:"tree_data"` // JSON string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TreeResponse struct {
	Project     string       `json:"project" doc:"Project name"`
	MainRequest string       `json:"mainRequest" doc:"Main project description"`
	Prompts     []PromptNode `json:"prompts" doc:"List of prompts with their nodes"`
}

type PromptNode struct {
	ID          int           `json:"id" doc:"Prompt ID"`
	Title       string        `json:"title" doc:"Prompt title"`
	Description string        `json:"description,omitempty" doc:"Prompt description"`
	Nodes       []NodeSummary `json:"nodes,omitempty" doc:"Child nodes of this prompt"`
}

type NodeSummary struct {
	ID     int    `json:"id" doc:"Node ID"`
	Name   string `json:"name" doc:"Node name"`
	Action string `json:"action,omitempty" doc:"Node action description"`
}

type PromptDetail struct {
	ID          int    `json:"id" doc:"Prompt ID"`
	Title       string `json:"title" doc:"Prompt title"`
	Description string `json:"description" doc:"Prompt description"`
	ProjectName string `json:"project_name" doc:"Parent project name"`
}

type CreatePromptRequest struct {
	Title       string `json:"title" minLength:"1" doc:"Title of the prompt (required)"`
	Description string `json:"description,omitempty" doc:"Description of the prompt"`
}

type CreateNodeRequest struct {
	Name   string `json:"name" minLength:"1" doc:"Name of the node (required)"`
	Action string `json:"action,omitempty" doc:"Action description"`
}

type CreateNoteRequest struct {
	Content string `json:"content" minLength:"1" doc:"Note content (required)"`
}

type UpdatePromptRequest struct {
	Title       string `json:"title,omitempty" doc:"Title of the prompt"`
	Description string `json:"description,omitempty" doc:"Description of the prompt"`
}

type UpdateNodeRequest struct {
	Name   string `json:"name,omitempty" doc:"Name of the node"`
	Action string `json:"action,omitempty" doc:"Action description"`
}

type UpdateNoteRequest struct {
	Content string `json:"content" minLength:"1" doc:"Note content (required)"`
}

type ImportTreeRequest struct {
	Tree TreeResponse `json:"tree" doc:"Complete tree structure to import"`
}

type SaveTreeRequest struct {
	Name string `json:"name" minLength:"1" doc:"Name for the saved tree (required)"`
}

type SavedTreeInfo struct {
	Name      string    `json:"name" doc:"Name of the saved tree"`
	CreatedAt time.Time `json:"created_at" doc:"When the tree was saved"`
	UpdatedAt time.Time `json:"updated_at" doc:"When the tree was last updated"`
}

type SavedTreeListResponse struct {
	Trees []SavedTreeInfo `json:"trees" doc:"List of saved trees"`
}