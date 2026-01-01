import { useState, useEffect } from 'react';
import NotesList from './NotesList';
import AddNote from './AddNote';
import { getNotes, createNote } from '../api/api';
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
 */
function SidePanel({ prompt }) {
  // State for notes
  const [notes, setNotes] = useState([]);
  const [loadingNotes, setLoadingNotes] = useState(false);
  const [notesError, setNotesError] = useState(null);

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
        <h2>{prompt.title}</h2>
        <p className="description">{prompt.description}</p>
      </div>

      {/* Nodes Section */}
      <div className="panel-section">
        <h3>
          <span className="section-icon"></span>
          Steps ({prompt.nodes?.length || 0})
        </h3>

        {prompt.nodes && prompt.nodes.length > 0 ? (
          <ul className="nodes-detail-list">
            {prompt.nodes.map((node, index) => (
              <li key={node.id || index} className="node-detail-item">
                <div className="node-header">
                  <span className="node-number">{index + 1}</span>
                  <strong>{node.name}</strong>
                </div>
                <p className="node-action">{node.action}</p>
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
          <NotesList notes={notes} />
        )}
      </div>
    </div>
  );
}

export default SidePanel;