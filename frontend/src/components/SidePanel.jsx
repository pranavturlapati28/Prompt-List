import { useState, useEffect } from 'react';
import NotesList from './NotesList';
import AddNote from './AddNote';
import PromptEditModal from './PromptEditModal';
import NodeEditModal from './NodeEditModal';
import ConfirmModal from './ConfirmModal';
import { getNotes, createNote, updatePrompt, deletePrompt, updateNode, deleteNode, createNode } from '../api/api';
import './SidePanel.css';

/**
 * SidePanel Component
 * 
 * Shows details for the selected prompt:
 * - Title and description
 * - List of child nodes with their actions
 * - Notes section with add/view functionality
 * 
 * Props:
 * - prompt: The currently selected prompt object
 * - onTreeUpdate: Callback to refresh the tree after updates
 */
function SidePanel({ prompt, onTreeUpdate }) {
  // State for notes
  const [notes, setNotes] = useState([]);
  const [loadingNotes, setLoadingNotes] = useState(false);
  const [notesError, setNotesError] = useState(null);

  // State for editing
  const [editingPrompt, setEditingPrompt] = useState(null);
  const [editingNode, setEditingNode] = useState(null);
  const [isCreatingNode, setIsCreatingNode] = useState(false);
  
  // Confirmation modal state
  const [confirmDelete, setConfirmDelete] = useState(null);

  // Fetch notes when selected prompt changes
  useEffect(() => {
    const fetchNotes = async () => {
      if (!prompt?.id) return;

      setLoadingNotes(true);
      setNotesError(null);

      try {
        const data = await getNotes(prompt.id);
        setNotes(data || []);
      } catch (err) {
        console.error('Error fetching notes:', err);
        setNotesError('Failed to load notes');
        setNotes([]);
      } finally {
        setLoadingNotes(false);
      }
    };

    fetchNotes();
  }, [prompt?.id]);

  // Handler for adding a new note
  const handleAddNote = async (content) => {
    try {
      const newNote = await createNote(prompt.id, content);
      // Add new note to the top of the list
      setNotes(prevNotes => [newNote, ...prevNotes]);
      return true; // Success
    } catch (err) {
      console.error('Error adding note:', err);
      return false; // Failed
    }
  };

  // Handler for updating prompt
  const handleUpdatePrompt = async (updatedData) => {
    if (!prompt) return;
    try {
      await updatePrompt(prompt.id, updatedData);
      if (onTreeUpdate) {
        await onTreeUpdate();
      }
      setEditingPrompt(null);
    } catch (error) {
      console.error('Error updating prompt:', error);
      alert('Failed to update prompt. Please try again.');
    }
  };

  // Handler for deleting prompt
  const handleDeletePrompt = () => {
    if (!prompt) return;
    setConfirmDelete({
      type: 'prompt',
      item: prompt,
      message: `Are you sure you want to delete "${prompt.title}"? This will also delete all associated nodes and notes.`,
    });
  };

  // Confirm and execute prompt delete
  const confirmDeletePrompt = async () => {
    if (!confirmDelete || confirmDelete.type !== 'prompt' || !prompt) return;
    
    try {
      await deletePrompt(prompt.id);
      if (onTreeUpdate) {
        await onTreeUpdate();
      }
    } catch (error) {
      console.error('Error deleting prompt:', error);
      alert('Failed to delete prompt. Please try again.');
    }
    
    setConfirmDelete(null);
  };

  // Handler for updating node
  const handleUpdateNode = async (nodeId, updatedData) => {
    if (!prompt) return;
    try {
      await updateNode(prompt.id, nodeId, updatedData);
      if (onTreeUpdate) {
        await onTreeUpdate();
      }
      setEditingNode(null);
    } catch (error) {
      console.error('Error updating node:', error);
      alert('Failed to update node. Please try again.');
    }
  };

  // Handler for creating new node
  const handleCreateNode = async (nodeData) => {
    if (!prompt) return;
    try {
      await createNode(prompt.id, nodeData);
      if (onTreeUpdate) {
        await onTreeUpdate();
      }
      setIsCreatingNode(false);
    } catch (error) {
      console.error('Error creating node:', error);
      alert('Failed to create node. Please try again.');
    }
  };

  // Handler for deleting node
  const handleDeleteNode = (node) => {
    if (!prompt) return;
    setConfirmDelete({
      type: 'node',
      item: node,
      promptId: prompt.id,
      message: `Are you sure you want to delete "${node.name}"?`,
    });
  };

  // Confirm and execute node delete
  const confirmDeleteNode = async () => {
    if (!confirmDelete || confirmDelete.type !== 'node' || !prompt) return;
    
    try {
      await deleteNode(prompt.id, confirmDelete.item.id);
      if (onTreeUpdate) {
        await onTreeUpdate();
      }
    } catch (error) {
      console.error('Error deleting node:', error);
      alert('Failed to delete node. Please try again.');
    }
    
    setConfirmDelete(null);
  };

  // Show empty state if no prompt selected
  if (!prompt) {
    return (
      <div className="side-panel empty">
        <div className="empty-state">
          <span className="empty-icon">Click prompt</span>
          <p>Select a prompt from the tree to view details</p>
        </div>
      </div>
    );
  }

  return (
    <div className="side-panel">
      {/* Header with prompt info */}
      <div className="panel-header">
        <span className="prompt-badge">Prompt {prompt.id}</span>
        <div className="editable-header">
          <h2>{prompt.title}</h2>
          <div className="action-icons">
            <button 
              className="icon-button" 
              onClick={() => setEditingPrompt(prompt)}
              title="Edit prompt"
            >
              ✏️
            </button>
            <button 
              className="icon-button icon-button-danger" 
              onClick={handleDeletePrompt}
              title="Delete prompt"
            >
              🗑️
            </button>
          </div>
        </div>
        <div className="editable-description">
          <p className="description">{prompt.description}</p>
          <div className="action-icons">
            <button 
              className="icon-button" 
              onClick={() => setEditingPrompt(prompt)}
              title="Edit description"
            >
              ✏️
            </button>
          </div>
        </div>
      </div>

      {/* Nodes Section */}
      <div className="panel-section">
        <div className="section-header-with-action">
          <h3>
            <span className="section-icon"></span>
            Steps ({prompt.nodes?.length || 0})
          </h3>
          <button 
            className="add-node-button" 
            onClick={() => setIsCreatingNode(true)}
            title="Add new subprompt"
          >
            +
          </button>
        </div>

        {prompt.nodes && prompt.nodes.length > 0 ? (
          <ul className="nodes-detail-list">
            {prompt.nodes.map((node, index) => (
              <li key={node.id || index} className="node-detail-item">
                <div className="node-header">
                  <span className="node-number">{index + 1}</span>
                  <div className="editable-node-name">
                    <strong>{node.name}</strong>
                    <div className="action-icons">
                      <button 
                        className="icon-button" 
                        onClick={() => setEditingNode({ node, promptId: prompt.id })}
                        title="Edit node name"
                      >
                        ✏️
                      </button>
                      <button 
                        className="icon-button icon-button-danger" 
                        onClick={() => handleDeleteNode(node)}
                        title="Delete node"
                      >
                        🗑️
                      </button>
                    </div>
                  </div>
                </div>
                <div className="editable-node-action">
                  <p className="node-action">{node.action}</p>
                  <div className="action-icons">
                    <button 
                      className="icon-button" 
                      onClick={() => setEditingNode({ node, promptId: prompt.id })}
                      title="Edit node action"
                    >
                      ✏️
                    </button>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        ) : (
          <p className="empty-text">No steps defined for this prompt</p>
        )}
      </div>

      {/* Notes Section */}
      <div className="panel-section">
        <h3>
          <span className="section-icon"></span>
          Notes ({notes.length})
        </h3>

        {/* Add Note Form */}
        <AddNote onAdd={handleAddNote} />

        {/* Notes List */}
        {loadingNotes ? (
          <p className="loading-text">Loading notes...</p>
        ) : notesError ? (
          <p className="error-text">{notesError}</p>
        ) : (
          <NotesList 
            notes={notes} 
            promptId={prompt.id} 
            onUpdate={async () => {
              // Refresh both tree and notes
              if (onTreeUpdate) {
                await onTreeUpdate();
              }
              // Refresh notes
              const data = await getNotes(prompt.id);
              setNotes(data || []);
            }} 
          />
        )}
      </div>

      {/* Edit Modals */}
      {editingPrompt && (
        <PromptEditModal
          prompt={editingPrompt}
          onClose={() => setEditingPrompt(null)}
          onSave={handleUpdatePrompt}
        />
      )}

      {editingNode && (
        <NodeEditModal
          node={editingNode.node}
          onClose={() => setEditingNode(null)}
          onSave={(data) => handleUpdateNode(editingNode.node.id, data)}
        />
      )}

      {isCreatingNode && (
        <NodeEditModal
          node={null}
          onClose={() => setIsCreatingNode(false)}
          onSave={handleCreateNode}
        />
      )}

      {/* Confirmation Modal */}
      <ConfirmModal
        isOpen={confirmDelete !== null}
        title="Confirm Delete"
        message={confirmDelete?.message || ''}
        onConfirm={() => {
          if (confirmDelete?.type === 'node') {
            confirmDeleteNode();
          } else if (confirmDelete?.type === 'prompt') {
            confirmDeletePrompt();
          }
        }}
        onCancel={() => setConfirmDelete(null)}
        confirmText="Delete"
        cancelText="Cancel"
      />
    </div>
  );
}

export default SidePanel;