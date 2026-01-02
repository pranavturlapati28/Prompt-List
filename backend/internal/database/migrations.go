package database

import "fmt"

func Migrate() error {
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

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS project_settings (
			id INTEGER PRIMARY KEY DEFAULT 1,
			project_name VARCHAR(255) NOT NULL DEFAULT 'Personal Finance Copilot',
			main_request TEXT NOT NULL DEFAULT 'Build a web app that helps users track spending, set goals, and get AI-powered budgeting advice from categorized transactions.',
			CONSTRAINT single_row CHECK (id = 1)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create project_settings table: %w", err)
	}
	
	_, err = DB.Exec(`
		INSERT INTO project_settings (id, project_name, main_request)
		VALUES (1, '3D Racing Game', 'Build a 3D racing video game in React Three Fiber where the player drives against AI opponents on a pregenerated racing track.')
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("failed to initialize project_settings: %w", err)
	}
	fmt.Println("✓ Project settings table ready")

	return nil
}