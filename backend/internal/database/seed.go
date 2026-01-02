package database

import (
	"fmt"
	"log"
)

func Seed() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM prompts").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing data: %w", err)
	}

	if count > 0 {
		fmt.Println("Clearing existing data...")
		_, err = DB.Exec("DELETE FROM notes")
		if err != nil {
			log.Printf("Warning: failed to delete notes: %v", err)
		}
		_, err = DB.Exec("DELETE FROM nodes")
		if err != nil {
			log.Printf("Warning: failed to delete nodes: %v", err)
		}
		_, err = DB.Exec("DELETE FROM prompts")
		if err != nil {
			log.Printf("Warning: failed to delete prompts: %v", err)
		}
		fmt.Println("✓ Existing data cleared")
	}

	fmt.Println("Seeding database with initial data...")

	prompts := []struct {
		title       string
		description string
	}{
		{
			"Project Setup",
			"Initialize repo, tooling, and baseline architecture.",
		},
		{
			"Transaction Ingestion",
			"Allow users to bring in transactions and normalize them.",
		},
		{
			"Categorization",
			"Automatically categorize transactions and let users correct them.",
		},
		{
			"Budgeting & Goals",
			"Help users define budgets and track progress.",
		},
		{
			"Insights Copilot",
			"Provide AI summaries and recommendations based on the user's data.",
		},
		{
			"Prompt Tree UI",
			"Build the interactive prompt tree + annotation UX.",
		},
		{
			"Shipping",
			"Polish, testing, and deployment.",
		},
	}

	promptIDs := make([]int, len(prompts))
	for i, p := range prompts {
		var id int
		err := DB.QueryRow(
			"INSERT INTO prompts (title, description, project_name) VALUES ($1, $2, $3) RETURNING id",
			p.title, p.description, "Personal Finance Copilot",
		).Scan(&id)
		if err != nil {
			log.Printf("Warning: failed to insert prompt '%s': %v", p.title, err)
			promptIDs[i] = i + 1 // Fallback to expected ID
		} else {
			promptIDs[i] = id
		}
	}

	// Insert all nodes (subprompts) from the Personal Finance Copilot project
	// Use the actual prompt IDs we got from inserting prompts
	nodes := []struct {
		promptIndex int // Index into promptIDs array
		name        string
		action      string
	}{
		// Prompt 1: Project Setup
		{0, "Scaffold app", "Create a Vite + React project with TypeScript. Add ESLint + Prettier and confirm the dev server runs."},
		{0, "Set up folder structure", "Create /components, /pages, /api, /hooks, /utils, and /styles. Add a simple route layout and global theme variables."},
		{0, "Wire backend client", "Create an API helper module that centralizes fetch calls, base URL, error handling, and JSON parsing."},

		// Prompt 2: Transaction Ingestion
		{1, "CSV upload", "Add a file upload UI that accepts CSV exports. Parse client-side and show a preview table before importing."},
		{1, "Normalization", "Normalize dates, amounts, and merchant names. Convert credits/debits into a consistent signed format."},
		{1, "Deduping", "Implement duplicate detection using (date, amount, merchant) fingerprinting. Show warnings and allow overrides."},
		{1, "Persistence", "Store transactions in your database with user scoping. Add indexes on date and category for fast filtering."},

		// Prompt 3: Categorization
		{2, "Rule-based baseline", "Start with keyword rules (e.g., UBER→Transport). Provide a UI for editing rules and re-running categorization."},
		{2, "ML/LLM categorizer", "Add an optional categorizer endpoint that predicts categories using merchant + memo fields. Log confidence scores."},
		{2, "Manual overrides", "Enable per-transaction category edits. Persist overrides and make sure overrides win over automated labeling."},

		// Prompt 4: Budgeting & Goals
		{3, "Budget builder", "Create a monthly budget interface by category with sliders/inputs. Show total budget vs expected income."},
		{3, "Progress tracking", "Compute spend-to-date and remaining budget. Visualize per-category burn with progress bars."},
		{3, "Savings goals", "Let users define goals (e.g., emergency fund) with target and deadline. Calculate required monthly contribution."},
		{3, "Alerts", "Add budget overrun alerts and unusual spending detection (e.g., 2σ above baseline)."},

		// Prompt 5: Insights Copilot
		{4, "Monthly summary", "Generate a narrative summary of spend patterns, top categories, and notable changes from last month."},
		{4, "What changed?", "Compute deltas by category and merchant. Surface the largest contributors to month-over-month change."},
		{4, "Recommendations", "Suggest 3–5 actionable improvements (e.g., reduce subscriptions). Each suggestion should cite the underlying data."},
		{4, "Explainability", "Add an 'Explain' button that shows exactly which transactions drove each recommendation."},

		// Prompt 6: Prompt Tree UI
		{5, "Tree visualization", "Render prompts as rows with subprompts as nodes. Support expand/collapse per prompt and selection highlight."},
		{5, "Side panel details", "Click any prompt or node to show details. Display title, description, and action text in the side panel."},
		{5, "Notes + annotations", "Allow users to add notes to prompts. Persist notes and show them inline in the side panel."},

		// Prompt 7: Shipping
		{6, "E2E test coverage", "Add tests for importing, categorization edits, budgets, and note creation. Test the main user flows."},
		{6, "Performance pass", "Add pagination/virtualization for large transaction tables. Ensure the tree remains responsive with many nodes."},
		{6, "Deployment", "Deploy backend and frontend. Provide a README with setup, environment variables, and API docs."},
	}

	for _, n := range nodes {
		// Use the actual prompt ID from the inserted prompts
		promptID := promptIDs[n.promptIndex]
		_, err := DB.Exec(
			"INSERT INTO nodes (prompt_id, name, action) VALUES ($1, $2, $3)",
			promptID, n.name, n.action,
		)
		if err != nil {
			log.Printf("Warning: failed to insert node '%s': %v", n.name, err)
		}
	}

	fmt.Println("✓ Database seeded successfully")
	return nil
}