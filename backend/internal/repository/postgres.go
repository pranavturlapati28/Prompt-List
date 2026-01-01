package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pranavturlapati28/merget-takehome/internal/database"
	"github.com/pranavturlapati28/merget-takehome/internal/models"
)

// PromptRepository handles all database operations
// This is the only layer that knows about SQL
type PromptRepository struct{}

// NewPromptRepository creates a new repository instance
func NewPromptRepository() *PromptRepository {
	return &PromptRepository{}
}

// =============================================================================
// PROMPT OPERATIONS
// =============================================================================

// GetAllPrompts retrieves all prompts from the database
func (r *PromptRepository) GetAllPrompts() ([]models.Prompt, error) {
	query := `
		SELECT id, title, description, project_name 
		FROM prompts 
		ORDER BY id
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close() // Always close rows to release connection

	var prompts []models.Prompt
	for rows.Next() {
		var p models.Prompt
		err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.ProjectName)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		prompts = append(prompts, p)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return prompts, nil
}

// GetPromptByID retrieves a single prompt by its ID
func (r *PromptRepository) GetPromptByID(id int) (*models.Prompt, error) {
	query := `
		SELECT id, title, description, project_name 
		FROM prompts 
		WHERE id = $1
	`

	var p models.Prompt
	err := database.DB.QueryRow(query, id).Scan(
		&p.ID, &p.Title, &p.Description, &p.ProjectName,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found (not an error)
	}
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &p, nil
}

// CreatePrompt inserts a new prompt into the database
func (r *PromptRepository) CreatePrompt(title, description string) (*models.Prompt, error) {
	query := `
		INSERT INTO prompts (title, description, project_name) 
		VALUES ($1, $2, $3) 
		RETURNING id
	`

	var id int
	err := database.DB.QueryRow(query, title, description, "3D Racing Game").Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	return &models.Prompt{
		ID:          id,
		Title:       title,
		Description: description,
		ProjectName: "3D Racing Game",
	}, nil
}

// PromptExists checks if a prompt with the given ID exists
func (r *PromptRepository) PromptExists(id int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM prompts WHERE id = $1)"
	err := database.DB.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("exists check failed: %w", err)
	}
	return exists, nil
}

// UpdatePrompt updates an existing prompt
func (r *PromptRepository) UpdatePrompt(id int, title, description string) (*models.Prompt, error) {
	// Build dynamic UPDATE query based on provided fields
	query := "UPDATE prompts SET"
	var args []interface{}
	argPos := 1

	if title != "" {
		query += fmt.Sprintf(" title = $%d", argPos)
		args = append(args, title)
		argPos++
	}

	if description != "" {
		if len(args) > 0 {
			query += ","
		}
		query += fmt.Sprintf(" description = $%d", argPos)
		args = append(args, description)
		argPos++
	}

	if len(args) == 0 {
		// No fields to update, just return the existing prompt
		return r.GetPromptByID(id)
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, title, description, project_name", argPos)
	args = append(args, id)

	var p models.Prompt
	err := database.DB.QueryRow(query, args...).Scan(&p.ID, &p.Title, &p.Description, &p.ProjectName)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}

	return &p, nil
}

// DeletePrompt deletes a prompt by ID (cascade will delete related nodes and notes)
func (r *PromptRepository) DeletePrompt(id int) error {
	query := "DELETE FROM prompts WHERE id = $1"
	result, err := database.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Not found
	}

	return nil
}

// =============================================================================
// NODE OPERATIONS
// =============================================================================

// GetNodesByPromptID retrieves all nodes for a specific prompt
func (r *PromptRepository) GetNodesByPromptID(promptID int) ([]models.Node, error) {
	query := `
		SELECT id, prompt_id, name, action 
		FROM nodes 
		WHERE prompt_id = $1 
		ORDER BY id
	`

	rows, err := database.DB.Query(query, promptID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var nodes []models.Node
	for rows.Next() {
		var n models.Node
		err := rows.Scan(&n.ID, &n.PromptID, &n.Name, &n.Action)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		nodes = append(nodes, n)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return nodes, nil
}

// CreateNode inserts a new node for a prompt
func (r *PromptRepository) CreateNode(promptID int, name, action string) (*models.Node, error) {
	query := `
		INSERT INTO nodes (prompt_id, name, action) 
		VALUES ($1, $2, $3) 
		RETURNING id
	`

	var id int
	err := database.DB.QueryRow(query, promptID, name, action).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	return &models.Node{
		ID:       id,
		PromptID: promptID,
		Name:     name,
		Action:   action,
	}, nil
}

// GetNodeByID retrieves a single node by its ID
func (r *PromptRepository) GetNodeByID(nodeID int) (*models.Node, error) {
	query := `
		SELECT id, prompt_id, name, action 
		FROM nodes 
		WHERE id = $1
	`

	var n models.Node
	err := database.DB.QueryRow(query, nodeID).Scan(&n.ID, &n.PromptID, &n.Name, &n.Action)

	if err == sql.ErrNoRows {
		return nil, nil // Not found (not an error)
	}
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &n, nil
}

// UpdateNode updates an existing node
func (r *PromptRepository) UpdateNode(nodeID int, name, action string) (*models.Node, error) {
	// Build dynamic UPDATE query based on provided fields
	query := "UPDATE nodes SET"
	var args []interface{}
	argPos := 1

	if name != "" {
		query += fmt.Sprintf(" name = $%d", argPos)
		args = append(args, name)
		argPos++
	}

	if action != "" {
		if len(args) > 0 {
			query += ","
		}
		query += fmt.Sprintf(" action = $%d", argPos)
		args = append(args, action)
		argPos++
	}

	if len(args) == 0 {
		// No fields to update, just return the existing node
		return r.GetNodeByID(nodeID)
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, prompt_id, name, action", argPos)
	args = append(args, nodeID)

	var n models.Node
	err := database.DB.QueryRow(query, args...).Scan(&n.ID, &n.PromptID, &n.Name, &n.Action)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}

	return &n, nil
}

// DeleteNode deletes a node by ID
func (r *PromptRepository) DeleteNode(nodeID int) error {
	query := "DELETE FROM nodes WHERE id = $1"
	result, err := database.DB.Exec(query, nodeID)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Not found
	}

	return nil
}

// =============================================================================
// NOTE OPERATIONS
// =============================================================================

// GetNotesByPromptID retrieves all notes for a specific prompt
func (r *PromptRepository) GetNotesByPromptID(promptID int) ([]models.Note, error) {
	query := `
		SELECT id, prompt_id, content, created_at 
		FROM notes 
		WHERE prompt_id = $1 
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query, promptID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		err := rows.Scan(&n.ID, &n.PromptID, &n.Content, &n.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		notes = append(notes, n)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return notes, nil
}

// CreateNote inserts a new note for a prompt
func (r *PromptRepository) CreateNote(promptID int, content string) (*models.Note, error) {
	query := `
		INSERT INTO notes (prompt_id, content, created_at) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`

	var id int
	var createdAt time.Time
	err := database.DB.QueryRow(query, promptID, content, time.Now()).Scan(&id, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	return &models.Note{
		ID:        id,
		PromptID:  promptID,
		Content:   content,
		CreatedAt: createdAt,
	}, nil
}

// GetNoteByID retrieves a single note by its ID
func (r *PromptRepository) GetNoteByID(noteID int) (*models.Note, error) {
	query := `
		SELECT id, prompt_id, content, created_at 
		FROM notes 
		WHERE id = $1
	`

	var n models.Note
	err := database.DB.QueryRow(query, noteID).Scan(&n.ID, &n.PromptID, &n.Content, &n.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // Not found (not an error)
	}
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &n, nil
}

// UpdateNote updates an existing note
func (r *PromptRepository) UpdateNote(noteID int, content string) (*models.Note, error) {
	query := `
		UPDATE notes 
		SET content = $1 
		WHERE id = $2 
		RETURNING id, prompt_id, content, created_at
	`

	var n models.Note
	err := database.DB.QueryRow(query, content, noteID).Scan(&n.ID, &n.PromptID, &n.Content, &n.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}

	return &n, nil
}

// DeleteNote deletes a note by ID
func (r *PromptRepository) DeleteNote(noteID int) error {
	query := "DELETE FROM notes WHERE id = $1"
	result, err := database.DB.Exec(query, noteID)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Not found
	}

	return nil
}