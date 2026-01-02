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
		fmt.Println("✓ Data already seeded (skipping)")
		return nil
	}

	fmt.Println("Seeding database with initial data...")

	_, err = DB.Exec(`
		UPDATE project_settings 
		SET project_name = $1, main_request = $2 
		WHERE id = 1
	`, "3D Racing Game", "Build a 3D racing video game in React Three Fiber where the player drives against AI opponents on a pregenerated racing track.")
	if err != nil {
		log.Printf("Warning: failed to update project_settings: %v", err)
	}

	prompts := []struct {
		title       string
		description string
	}{
		{
			"Project Setup",
			"Initialize the development environment and install all required dependencies.",
		},
		{
			"3D Environment",
			"Build the visual atmosphere and world surrounding the race track.",
		},
		{
			"Racing Track",
			"Generate the prebuilt racing circuit with all necessary geometry and race markers.",
		},
		{
			"Player Vehicle",
			"Implement the user-controlled car with physics and camera.",
		},
		{
			"AI Opponents",
			"Create computer-controlled vehicles that race against the player.",
		},
		{
			"Game Systems",
			"Implement core racing game logic and state management.",
		},
		{
			"UI / HUD",
			"Build the heads-up display and menu interfaces.",
		},
	}

	promptIDs := make([]int, len(prompts))
	for i, p := range prompts {
		var id int
		err := DB.QueryRow(
			"INSERT INTO prompts (title, description, project_name) VALUES ($1, $2, $3) RETURNING id",
			p.title, p.description, "3D Racing Game",
		).Scan(&id)
		if err != nil {
			log.Printf("Warning: failed to insert prompt '%s': %v", p.title, err)
			promptIDs[i] = i + 1
		} else {
			promptIDs[i] = id
		}
	}

	nodes := []struct {
		promptIndex int
		name        string
		action      string
	}{
		{0, "npm create vite", "Scaffold a new React project using Vite as the build tool. Select React + TypeScript template for type safety and fast HMR."},
		{0, "Install dependencies", "Add core 3D libraries: @react-three/fiber for React Three.js bindings, @react-three/drei for helpful abstractions, and @react-three/rapier for physics simulation."},
		{0, "Create folder structure", "Organize the project into logical directories: components/ for React/3D components, hooks/ for custom hooks, utils/ for helper functions, assets/ for models and textures, and stores/ for game state management."},

		{1, "HDRI skybox", "Load a high dynamic range image as the scene environment using drei's <Environment> component. This provides realistic sky rendering and image-based lighting for reflections on car surfaces."},
		{1, "Lighting setup", "Configure a directional light to simulate the sun with shadows enabled. Add ambient light for fill. Adjust intensity and shadow map resolution for performance balance."},
		{1, "Ground plane", "Create a large ground mesh extending beyond the track with a grass or terrain texture. Apply a repeating material and ensure it receives shadows from vehicles and track elements."},

		{2, "Track spline", "Define a CatmullRomCurve3 path using an array of Vector3 control points that form the racing line. This spline serves as the foundation for track generation and AI navigation."},
		{2, "Track mesh", "Extrude a road cross-section shape along the spline to create the track surface geometry. Apply asphalt texture with UV mapping that follows the curve. Add the mesh as a physics collider."},
		{2, "Barriers", "Generate wall meshes along both edges of the track using offset splines. Add RigidBody colliders to prevent cars from leaving the circuit. Style with tire wall or concrete barrier textures."},
		{2, "Checkpoints", "Create invisible trigger zones at regular intervals around the track using sensor colliders. These track player progress and prevent lap-skipping by requiring sequential checkpoint passage."},
		{2, "Start/finish line", "Place a visual marker mesh (checkered pattern) at the race origin. Add a dedicated trigger zone that increments lap count when crossed after completing all checkpoints."},

		{3, "Car model", "Load a 3D car model (GLTF/GLB format) using useGLTF hook. Ensure the model has separate wheel meshes for rotation animation. Apply materials and set appropriate scale."},
		{3, "Vehicle physics", "Create a dynamic RigidBody for the car chassis. Implement a raycast vehicle controller with four wheel configurations including suspension stiffness, friction, and roll influence parameters."},
		{3, "Keyboard controls", "Set up input handling for WASD or arrow keys. Map vertical axis to acceleration/braking force applied to wheels. Map horizontal axis to steering angle with smooth interpolation."},
		{3, "Chase camera", "Implement a third-person camera that follows behind the player car. Use useFrame to smoothly lerp camera position and look-at target. Add slight lag for dynamic feel during turns."},

		{4, "Spawn AI cars", "Instantiate multiple opponent vehicles at staggered starting positions on the grid. Use the same car model with different color materials. Each AI car gets its own RigidBody and state."},
		{4, "Pathfinding", "Implement spline-following behavior where AI cars steer toward the next waypoint along the track curve. Sample points ahead on the spline and calculate steering angle to reach them."},
		{4, "Speed AI", "Add randomized speed multipliers to each AI car for varied difficulty. Implement acceleration curves and braking logic when approaching sharp turns based on track curvature analysis."},
		{4, "Collision avoidance", "Cast rays forward and to sides from each AI car. When obstacles are detected, apply lateral steering adjustments to avoid collisions with walls and other vehicles."},

		{5, "Lap counting", "Track each vehicle's checkpoint progress in a state store. When a car crosses the finish line trigger with all checkpoints cleared, increment their lap counter and reset checkpoint flags."},
		{5, "Position tracking", "Calculate race positions by comparing each car's lap count and progress percentage along the track spline. Update positions in real-time and store for UI display."},
		{5, "Race timer", "Start a countdown timer at race begin (3-2-1-GO sequence). Track elapsed race time and individual lap times. Store best lap time for display. Pause timer when race ends."},
		{5, "Win/lose logic", "Define race completion as finishing a set number of laps (e.g., 3). Determine final standings when all cars finish or timeout. Trigger end-race state with results display."},

		{6, "Speedometer", "Create an overlay component displaying current player speed. Calculate from vehicle velocity magnitude. Style as digital readout or analog gauge with needle animation."},
		{6, "Position display", "Show player's current race position prominently (e.g., '2nd / 4'). Update in real-time as positions change. Add ordinal suffix formatting (1st, 2nd, 3rd)."},
		{6, "Lap counter", "Display current lap number and total laps (e.g., 'Lap 2 / 3'). Show current lap time and best lap time below. Flash or highlight on new best lap."},
		{6, "Mini-map", "Render a top-down 2D view of the track in a corner overlay. Show dots for all car positions color-coded by player/AI. Rotate map to match player heading or keep north-up."},
		{6, "Menus", "Create start screen with race configuration options. Implement pause menu with resume/restart/quit options. Build results screen showing final standings, times, and replay option."},
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