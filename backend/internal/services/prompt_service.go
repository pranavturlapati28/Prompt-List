package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pranavturlapati28/merget-takehome/internal/models"
	"github.com/pranavturlapati28/merget-takehome/internal/repository"
)

// Common errors that can be returned by the service
var (
	ErrPromptNotFound = errors.New("prompt not found")
	ErrNodeNotFound   = errors.New("node not found")
	ErrNoteNotFound   = errors.New("note not found")
)

// PromptService contains all business logic for prompts
// It orchestrates between the API layer and the repository
type PromptService struct {
	repo *repository.PromptRepository
}

// NewPromptService creates a new service with its dependencies
func NewPromptService(repo *repository.PromptRepository) *PromptService {
	return &PromptService{repo: repo}
}

// =============================================================================
// TREE OPERATIONS
// =============================================================================

// GetTree returns the complete prompt tree with all nodes
// This is the main data structure for the frontend visualization
func (s *PromptService) GetTree() (*models.TreeResponse, error) {
	// Get all prompts
	prompts, err := s.repo.GetAllPrompts()
	if err != nil {
		return nil, err
	}

	// Build the tree structure
	var promptNodes []models.PromptNode

	for _, p := range prompts {
		// Get nodes for each prompt
		nodes, err := s.repo.GetNodesByPromptID(p.ID)
		if err != nil {
			return nil, err
		}

		// Convert to NodeSummary for the response
		var nodeSummaries []models.NodeSummary
		for _, n := range nodes {
			nodeSummaries = append(nodeSummaries, models.NodeSummary{
				ID:     n.ID,
				Name:   n.Name,
				Action: n.Action,
			})
		}

		// Ensure nodes array is never nil (cleaner JSON)
		if nodeSummaries == nil {
			nodeSummaries = []models.NodeSummary{}
		}

		promptNodes = append(promptNodes, models.PromptNode{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			Nodes:       nodeSummaries,
		})
	}

	// Ensure prompts array is never nil
	if promptNodes == nil {
		promptNodes = []models.PromptNode{}
	}

	return &models.TreeResponse{
		Project:     "Personal Finance Copilot",
		MainRequest: "Build a web app that helps users track spending, set goals, and get AI-powered budgeting advice from categorized transactions.",
		Prompts:     promptNodes,
	}, nil
}

// =============================================================================
// PROMPT OPERATIONS
// =============================================================================

// GetPrompt retrieves a single prompt by ID
func (s *PromptService) GetPrompt(id int) (*models.PromptDetail, error) {
	prompt, err := s.repo.GetPromptByID(id)
	if err != nil {
		return nil, err
	}
	if prompt == nil {
		return nil, ErrPromptNotFound
	}

	return &models.PromptDetail{
		ID:          prompt.ID,
		Title:       prompt.Title,
		Description: prompt.Description,
		ProjectName: prompt.ProjectName,
	}, nil
}

// CreatePrompt creates a new prompt
func (s *PromptService) CreatePrompt(title, description string) (*models.Prompt, error) {
	return s.repo.CreatePrompt(title, description)
}

// UpdatePrompt updates an existing prompt
func (s *PromptService) UpdatePrompt(id int, title, description string) (*models.Prompt, error) {
	// Check if prompt exists
	exists, err := s.repo.PromptExists(id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrPromptNotFound
	}

	return s.repo.UpdatePrompt(id, title, description)
}

// DeletePrompt deletes a prompt by ID
func (s *PromptService) DeletePrompt(id int) error {
	// Check if prompt exists
	exists, err := s.repo.PromptExists(id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrPromptNotFound
	}

	return s.repo.DeletePrompt(id)
}

// =============================================================================
// NODE OPERATIONS
// =============================================================================

// GetPromptNodes retrieves all nodes for a prompt
func (s *PromptService) GetPromptNodes(promptID int) ([]models.Node, error) {
	// First check if the prompt exists
	exists, err := s.repo.PromptExists(promptID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrPromptNotFound
	}

	nodes, err := s.repo.GetNodesByPromptID(promptID)
	if err != nil {
		return nil, err
	}

	// Ensure never nil
	if nodes == nil {
		nodes = []models.Node{}
	}

	return nodes, nil
}

// CreateNode creates a new node for a prompt
func (s *PromptService) CreateNode(promptID int, name, action string) (*models.Node, error) {
	// First check if the prompt exists
	exists, err := s.repo.PromptExists(promptID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrPromptNotFound
	}

	return s.repo.CreateNode(promptID, name, action)
}

// UpdateNode updates an existing node
func (s *PromptService) UpdateNode(nodeID int, name, action string) (*models.Node, error) {
	// Check if node exists
	node, err := s.repo.GetNodeByID(nodeID)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, ErrNodeNotFound
	}

	return s.repo.UpdateNode(nodeID, name, action)
}

// DeleteNode deletes a node by ID
func (s *PromptService) DeleteNode(nodeID int) error {
	// Check if node exists
	node, err := s.repo.GetNodeByID(nodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return ErrNodeNotFound
	}

	return s.repo.DeleteNode(nodeID)
}

// =============================================================================
// NOTE OPERATIONS
// =============================================================================

// GetNotes retrieves all notes for a prompt
func (s *PromptService) GetNotes(promptID int) ([]models.Note, error) {
	// First check if the prompt exists
	exists, err := s.repo.PromptExists(promptID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrPromptNotFound
	}

	notes, err := s.repo.GetNotesByPromptID(promptID)
	if err != nil {
		return nil, err
	}

	// Ensure never nil
	if notes == nil {
		notes = []models.Note{}
	}

	return notes, nil
}

// CreateNote creates a new note for a prompt
func (s *PromptService) CreateNote(promptID int, content string) (*models.Note, error) {
	// First check if the prompt exists
	exists, err := s.repo.PromptExists(promptID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrPromptNotFound
	}

	return s.repo.CreateNote(promptID, content)
}

// UpdateNote updates an existing note
func (s *PromptService) UpdateNote(noteID int, content string) (*models.Note, error) {
	// Check if note exists
	note, err := s.repo.GetNoteByID(noteID)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, ErrNoteNotFound
	}

	return s.repo.UpdateNote(noteID, content)
}

// DeleteNote deletes a note by ID
func (s *PromptService) DeleteNote(noteID int) error {
	// Check if note exists
	note, err := s.repo.GetNoteByID(noteID)
	if err != nil {
		return err
	}
	if note == nil {
		return ErrNoteNotFound
	}

	return s.repo.DeleteNote(noteID)
}

// =============================================================================
// TREE IMPORT/EXPORT/SAVE/LOAD OPERATIONS
// =============================================================================

// ImportTree replaces the current tree with imported data
func (s *PromptService) ImportTree(treeData *models.TreeResponse) error {
	// Validate tree structure
	if treeData.Project == "" {
		return errors.New("project name is required")
	}
	if len(treeData.Prompts) == 0 {
		return errors.New("at least one prompt is required")
	}

	// Validate each prompt
	for _, prompt := range treeData.Prompts {
		if prompt.Title == "" {
			return errors.New("all prompts must have a title")
		}
		// Validate nodes
		for _, node := range prompt.Nodes {
			if node.Name == "" {
				return errors.New("all nodes must have a name")
			}
		}
	}

	return s.repo.ImportTree(treeData)
}

// SaveTree saves the current tree with a name
func (s *PromptService) SaveTree(name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	// Get current tree
	tree, err := s.GetTree()
	if err != nil {
		return fmt.Errorf("failed to get current tree: %w", err)
	}

	// Convert to JSON
	treeJSON, err := json.Marshal(tree)
	if err != nil {
		return fmt.Errorf("failed to marshal tree: %w", err)
	}

	return s.repo.SaveTree(name, string(treeJSON))
}

// LoadTree loads a saved tree and replaces the current tree
func (s *PromptService) LoadTree(name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	// Get saved tree
	savedTree, err := s.repo.GetSavedTree(name)
	if err != nil {
		return err
	}
	if savedTree == nil {
		return errors.New("saved tree not found")
	}

	// Parse JSON
	var treeData models.TreeResponse
	err = json.Unmarshal([]byte(savedTree.TreeData), &treeData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal tree data: %w", err)
	}

	// Import the tree
	return s.ImportTree(&treeData)
}

// ListSavedTrees returns all saved tree names and metadata
func (s *PromptService) ListSavedTrees() ([]models.SavedTreeInfo, error) {
	return s.repo.ListSavedTrees()
}

// DeleteSavedTree deletes a saved tree by name
func (s *PromptService) DeleteSavedTree(name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	return s.repo.DeleteSavedTree(name)
}