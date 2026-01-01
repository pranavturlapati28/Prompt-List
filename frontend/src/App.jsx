import { useState, useEffect } from 'react';
import TreeView from './components/TreeView';
import SidePanel from './components/SidePanel';
import { getTree } from './api/api';
import './App.css';

/**
 * Main Application Component
 * 
 * Layout:
 * - Header with project title
 * - Main area split into:
 *   - TreeView (left): Interactive tree visualization
 *   - SidePanel (right): Details and notes for selected prompt
 */
function App() {
  // State for the tree data from API
  const [tree, setTree] = useState(null);
  
  // State for which prompt is currently selected
  const [selectedPrompt, setSelectedPrompt] = useState(null);
  
  // Loading state while fetching data
  const [loading, setLoading] = useState(true);
  
  // Error state if API call fails
  const [error, setError] = useState(null);

  // Fetch tree data function
  const fetchTree = async () => {
    try {
      const data = await getTree();
      setTree(data);
      setLoading(false);
    } catch (err) {
      console.error('Error fetching tree:', err);
      setError('Failed to load prompt tree. Make sure the backend is running on http://localhost:8080');
      setLoading(false);
    }
  };

  // Fetch tree data when component mounts
  useEffect(() => {
    fetchTree();
  }, []); // Empty dependency array = run once on mount

  // Handler when user clicks a prompt in the tree
  const handleSelectPrompt = (prompt) => {
    setSelectedPrompt(prompt);
  };

  // Show loading spinner while fetching
  if (loading) {
    return (
      <div className="loading-screen">
        <div className="loading-spinner"></div>
        <p>Loading prompt tree...</p>
      </div>
    );
  }

  // Show error message if fetch failed
  if (error) {
    return (
      <div className="error-screen">
        <h2>Connection Error</h2>
        <p>{error}</p>
        <button onClick={() => window.location.reload()}>
          Retry
        </button>
      </div>
    );
  }

  // Main app UI
  return (
    <div className="app">
      {/* Header */}
      <header className="header">
        <div className="header-content">
          <div className="header-left">
            <h1>{tree?.project || 'Prompt Tree Explorer'}</h1>
            <p>{tree?.mainRequest}</p>
          </div>
          <div className="header-right">
            <div className="instructions">
              <span className="instructions-icon">ℹ</span>
              <span>Right-click nodes to edit or delete</span>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="main">
        {/* Left side: Tree visualization */}
        <div className="tree-container">
          <TreeView
            prompts={tree?.prompts || []}
            onSelectPrompt={handleSelectPrompt}
            selectedPromptId={selectedPrompt?.id}
            onTreeUpdate={fetchTree}
          />
        </div>

        {/* Right side: Details panel */}
        <div className="panel-container">
          <SidePanel prompt={selectedPrompt} onTreeUpdate={fetchTree} />
        </div>
      </main>
    </div>
  );
}

export default App;