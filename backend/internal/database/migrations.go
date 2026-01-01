package database

import "fmt"

// Migrate creates all required database tables
func Migrate() error {
	// Create prompts table
	// This stores the main prompts like "Project Setup", "3D Environment", etc.
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS prompts (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			project_name VARCHAR(255) DEFAULT '3D Racing Game'
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create prompts table: %w", err)
	}
	fmt.Println("✓ Prompts table ready")

	// Create nodes table
	// This stores the subprompts/steps under each prompt
	// Foreign key ensures referential integrity with prompts table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS nodes (
			id SERIAL PRIMARY KEY,
			prompt_id INTEGER REFERENCES prompts(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			action TEXT
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create nodes table: %w", err)
	}
	fmt.Println("✓ Nodes table ready")

	// Create notes table
	// This stores user annotations on prompts
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id SERIAL PRIMARY KEY,
			prompt_id INTEGER REFERENCES prompts(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create notes table: %w", err)
	}
	fmt.Println("✓ Notes table ready")

	// Create saved_trees table
	// This stores saved prompt tree configurations
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS saved_trees (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			tree_data JSONB NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create saved_trees table: %w", err)
	}
	fmt.Println("✓ Saved trees table ready")

	return nil
}