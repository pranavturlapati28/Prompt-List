package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pranavturlapati28/merget-takehome/internal/database"
	"github.com/pranavturlapati28/merget-takehome/internal/models"
)

type PromptRepository struct{}

func NewPromptRepository() *PromptRepository {
	return &PromptRepository{}
}

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
	defer rows.Close()

	var prompts []models.Prompt
	for rows.Next() {
		var p models.Prompt
		err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.ProjectName)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		prompts = append(prompts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return prompts, nil
}

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

func (r *PromptRepository) PromptExists(id int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM prompts WHERE id = $1)"
	err := database.DB.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("exists check failed: %w", err)
	}
	return exists, nil
}

func (r *PromptRepository) UpdatePrompt(id int, title, description string) (*models.Prompt, error) {
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

func (r *PromptRepository) UpdateNode(nodeID int, name, action string) (*models.Node, error) {
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


// SaveTree saves a tree configuration with a name
func (r *PromptRepository) SaveTree(name string, treeData string) error {
	query := `
		INSERT INTO saved_trees (name, tree_data, updated_at)
		VALUES ($1, $2::jsonb, CURRENT_TIMESTAMP)
		ON CONFLICT (name) 
		DO UPDATE SET 
			tree_data = EXCLUDED.tree_data,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := database.DB.Exec(query, name, treeData)
	if err != nil {
		return fmt.Errorf("save failed: %w", err)
	}

	return nil
}

func (r *PromptRepository) GetSavedTree(name string) (*models.SavedTree, error) {
	query := `
		SELECT id, name, tree_data::text, created_at, updated_at
		FROM saved_trees
		WHERE name = $1
	`

	var st models.SavedTree
	err := database.DB.QueryRow(query, name).Scan(&st.ID, &st.Name, &st.TreeData, &st.CreatedAt, &st.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &st, nil
}

func (r *PromptRepository) ListSavedTrees() ([]models.SavedTreeInfo, error) {
	query := `
		SELECT name, created_at, updated_at
		FROM saved_trees
		ORDER BY updated_at DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var trees []models.SavedTreeInfo
	for rows.Next() {
		var st models.SavedTreeInfo
		err := rows.Scan(&st.Name, &st.CreatedAt, &st.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		trees = append(trees, st)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return trees, nil
}

func (r *PromptRepository) DeleteSavedTree(name string) error {
	query := "DELETE FROM saved_trees WHERE name = $1"
	result, err := database.DB.Exec(query, name)
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

func (r *PromptRepository) ImportTree(treeData *models.TreeResponse) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("transaction begin failed: %w", err)
	}
	defer tx.Rollback()

	fmt.Printf("About to update project_settings with: project=%s, mainRequest=%s\n", treeData.Project, treeData.MainRequest)
	result, err := tx.Exec(`
		INSERT INTO project_settings (id, project_name, main_request)
		VALUES (1, $1, $2)
		ON CONFLICT (id) 
		DO UPDATE SET 
			project_name = EXCLUDED.project_name,
			main_request = EXCLUDED.main_request
	`, treeData.Project, treeData.MainRequest)
	if err != nil {
		fmt.Printf("ERROR updating project_settings: %v\n", err)
		return fmt.Errorf("update project settings failed: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("SUCCESS: Updated project_settings: rows affected = %d, project = %s\n", rowsAffected, treeData.Project)

	_, err = tx.Exec("TRUNCATE TABLE notes, nodes, prompts CASCADE")
	if err != nil {
		return fmt.Errorf("truncate failed: %w", err)
	}

	for _, promptNode := range treeData.Prompts {
		var newID int
		err := tx.QueryRow(`
			INSERT INTO prompts (title, description, project_name)
			VALUES ($1, $2, $3)
			RETURNING id
		`, promptNode.Title, promptNode.Description, treeData.Project).Scan(&newID)

		if err != nil {
			return fmt.Errorf("insert prompt failed: %w", err)
		}

		for _, nodeSummary := range promptNode.Nodes {
			_, err = tx.Exec(`
				INSERT INTO nodes (prompt_id, name, action)
				VALUES ($1, $2, $3)
			`, newID, nodeSummary.Name, nodeSummary.Action)

			if err != nil {
				return fmt.Errorf("insert node failed: %w", err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		fmt.Printf("ERROR committing transaction: %v\n", err)
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	fmt.Printf("=== ImportTree SUCCESS: project=%s ===\n", treeData.Project)
	return nil
}

func (r *PromptRepository) GetProjectSettings() (string, string, error) {
	var projectName, mainRequest string
	query := "SELECT project_name, main_request FROM project_settings WHERE id = 1"
	err := database.DB.QueryRow(query).Scan(&projectName, &mainRequest)
	if err != nil {
		return "", "", fmt.Errorf("query failed: %w", err)
	}
	return projectName, mainRequest, nil
}