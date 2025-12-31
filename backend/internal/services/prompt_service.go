package services

import (
	"errors"

	"github.com/pranavturlapati28/merget-takehome/internal/models"
	"github.com/pranavturlapati28/merget-takehome/internal/repository"
)

// Common errors that can be returned by the service
var (
	ErrPromptNotFound = errors.New("prompt not found")
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
		Project:     "3D Racing Game",
		MainRequest: "Build a 3D racing video game in React Three Fiber where the player drives against AI opponents on a pregenerated racing track.",
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