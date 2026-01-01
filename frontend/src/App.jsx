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

  // Fetch tree data when component mounts
  useEffect(() => {
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
          <h1>{tree?.project || 'Prompt Tree Explorer'}</h1>
          <p>{tree?.mainRequest}</p>
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
          />
        </div>

        {/* Right side: Details panel */}
        <div className="panel-container">
          <SidePanel prompt={selectedPrompt} />
        </div>
      </main>
    </div>
  );
}

export default App;