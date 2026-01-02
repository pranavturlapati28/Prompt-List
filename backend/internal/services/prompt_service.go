package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/pranavturlapati28/merget-takehome/internal/models"
	"github.com/pranavturlapati28/merget-takehome/internal/repository"
)

var (
	ErrPromptNotFound = errors.New("prompt not found")
	ErrNodeNotFound   = errors.New("node not found")
	ErrNoteNotFound   = errors.New("note not found")
)

type PromptService struct {
	repo *repository.PromptRepository
}

func NewPromptService(repo *repository.PromptRepository) *PromptService {
	return &PromptService{repo: repo}
}

func (s *PromptService) GetTree() (*models.TreeResponse, error) {
	prompts, err := s.repo.GetAllPrompts()
	if err != nil {
		return nil, err
	}

	var promptNodes []models.PromptNode

	for _, p := range prompts {
		nodes, err := s.repo.GetNodesByPromptID(p.ID)
		if err != nil {
			return nil, err
		}

		var nodeSummaries []models.NodeSummary
		for _, n := range nodes {
			nodeSummaries = append(nodeSummaries, models.NodeSummary{
				ID:     n.ID,
				Name:   n.Name,
				Action: n.Action,
			})
		}

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

	if promptNodes == nil {
		promptNodes = []models.PromptNode{}
	}

	projectName, mainRequest, err := s.repo.GetProjectSettings()
	if err != nil {
		log.Printf("Warning: Failed to get project settings: %v\n", err)
		projectName = "Personal Finance Copilot"
		mainRequest = "Build a web app that helps users track spending, set goals, and get AI-powered budgeting advice from categorized transactions."
	} else {
		log.Printf("Retrieved project settings: %s - %s\n", projectName, mainRequest)
	}

	return &models.TreeResponse{
		Project:     projectName,
		MainRequest: mainRequest,
		Prompts:     promptNodes,
	}, nil
}

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

func (s *PromptService) CreatePrompt(title, description string) (*models.Prompt, error) {
	return s.repo.CreatePrompt(title, description)
}

func (s *PromptService) UpdatePrompt(id int, title, description string) (*models.Prompt, error) {
	exists, err := s.repo.PromptExists(id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrPromptNotFound
	}

	return s.repo.UpdatePrompt(id, title, description)
}

func (s *PromptService) DeletePrompt(id int) error {
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

	if nodes == nil {
		nodes = []models.Node{}
	}

	return nodes, nil
}

func (s *PromptService) CreateNode(promptID int, name, action string) (*models.Node, error) {
	exists, err := s.repo.PromptExists(promptID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrPromptNotFound
	}

	return s.repo.CreateNode(promptID, name, action)
}

func (s *PromptService) UpdateNode(nodeID int, name, action string) (*models.Node, error) {
	node, err := s.repo.GetNodeByID(nodeID)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, ErrNodeNotFound
	}

	return s.repo.UpdateNode(nodeID, name, action)
}

func (s *PromptService) DeleteNode(nodeID int) error {
	node, err := s.repo.GetNodeByID(nodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return ErrNodeNotFound
	}

	return s.repo.DeleteNode(nodeID)
}

func (s *PromptService) GetNotes(promptID int) ([]models.Note, error) {
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

func (s *PromptService) CreateNote(promptID int, content string) (*models.Note, error) {
	exists, err := s.repo.PromptExists(promptID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrPromptNotFound
	}

	return s.repo.CreateNote(promptID, content)
}

func (s *PromptService) UpdateNote(noteID int, content string) (*models.Note, error) {
	note, err := s.repo.GetNoteByID(noteID)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, ErrNoteNotFound
	}

	return s.repo.UpdateNote(noteID, content)
}

func (s *PromptService) DeleteNote(noteID int) error {
	note, err := s.repo.GetNoteByID(noteID)
	if err != nil {
		return err
	}
	if note == nil {
		return ErrNoteNotFound
	}

	return s.repo.DeleteNote(noteID)
}

func (s *PromptService) ImportTree(treeData *models.TreeResponse) error {
	if treeData.Project == "" {
		return errors.New("project name is required")
	}
	if len(treeData.Prompts) == 0 {
		return errors.New("at least one prompt is required")
	}

	for _, prompt := range treeData.Prompts {
		if prompt.Title == "" {
			return errors.New("all prompts must have a title")
		}
		for _, node := range prompt.Nodes {
			if node.Name == "" {
				return errors.New("all nodes must have a name")
			}
		}
	}

	log.Printf("Importing tree with project: %s, mainRequest: %s\n", treeData.Project, treeData.MainRequest)
	err := s.repo.ImportTree(treeData)
	if err != nil {
		log.Printf("Error importing tree: %v\n", err)
		return err
	}
	log.Printf("Tree imported successfully\n")
	return nil
}

func (s *PromptService) SaveTree(name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	// Get current tree
	tree, err := s.GetTree()
	if err != nil {
		return fmt.Errorf("failed to get current tree: %w", err)
	}

	treeJSON, err := json.Marshal(tree)
	if err != nil {
		return fmt.Errorf("failed to marshal tree: %w", err)
	}

	return s.repo.SaveTree(name, string(treeJSON))
}

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

	var treeData models.TreeResponse
	err = json.Unmarshal([]byte(savedTree.TreeData), &treeData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal tree data: %w", err)
	}

	return s.ImportTree(&treeData)
}

func (s *PromptService) ListSavedTrees() ([]models.SavedTreeInfo, error) {
	return s.repo.ListSavedTrees()
}

func (s *PromptService) DeleteSavedTree(name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	return s.repo.DeleteSavedTree(name)
}