import { useState, useEffect } from 'react';
import TreeView from './components/TreeView';
import SidePanel from './components/SidePanel';
import TreeManager from './components/TreeManager';
import { getTree } from './api/api';
import './App.css';

function App() {
  const [tree, setTree] = useState(null);
  const [selectedPrompt, setSelectedPrompt] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchTree = async () => {
    try {
      const data = await getTree();
      setTree(data);
      setLoading(false);
      setError(null);
    } catch (err) {
      console.error('Error fetching tree:', err);
      setError('Failed to load prompt tree. Make sure the backend is running on http://localhost:8080');
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTree();
  }, []);

  const handleSelectPrompt = (prompt) => {
    setSelectedPrompt(prompt);
  };

  if (loading) {
    return (
      <div className="loading-screen">
        <div className="loading-spinner"></div>
        <p>Loading prompt tree...</p>
      </div>
    );
  }

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

  return (
    <div className="app">
      <header className="header">
        <div className="header-content">
          <div className="header-left">
            <h1>{tree?.project || 'Prompt Tree Explorer'}</h1>
            <p>{tree?.mainRequest || ''}</p>
            <TreeManager treeData={tree} onTreeUpdate={fetchTree} />
          </div>
          <div className="header-right">
            <div className="instructions">
              <span className="instructions-icon">ℹ</span>
              <span>Right-click nodes to edit or delete</span>
            </div>
          </div>
        </div>
      </header>

      <main className="main">
        <div className="tree-container">
          <TreeView
            prompts={tree?.prompts || []}
            onSelectPrompt={handleSelectPrompt}
            selectedPromptId={selectedPrompt?.id}
            onTreeUpdate={fetchTree}
          />
        </div>

        <div className="panel-container">
          <SidePanel prompt={selectedPrompt} onTreeUpdate={fetchTree} />
        </div>
      </main>
    </div>
  );
}

export default App;