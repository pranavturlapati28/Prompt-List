import { useState } from 'react';
import './NotesList.css';
import NoteEditModal from './NoteEditModal';
import ConfirmModal from './ConfirmModal';
import { updateNote, deleteNote } from '../api/api';

/**
 * NotesList Component
 * 
 * Displays a list of notes/annotations for a prompt
 * 
 * Props:
 * - notes: Array of note objects with id, content, and created_at
 * - promptId: ID of the prompt these notes belong to
 * - onUpdate: Callback to refresh data after updates
 */
function NotesList({ notes, promptId, onUpdate }) {
  const [editingNote, setEditingNote] = useState(null);
  const [confirmDelete, setConfirmDelete] = useState(null);

  const handleUpdateNote = async (content) => {
    if (!editingNote || !promptId) return;
    try {
      await updateNote(promptId, editingNote.id, content);
      if (onUpdate) {
        await onUpdate();
      }
      setEditingNote(null);
    } catch (error) {
      console.error('Error updating note:', error);
      alert('Failed to update note. Please try again.');
    }
  };

  const handleDeleteNote = (note) => {
    if (!promptId) return;
    setConfirmDelete({
      type: 'note',
      item: note,
      promptId,
      message: 'Are you sure you want to delete this note?',
    });
  };

  const confirmDeleteNote = async () => {
    if (!confirmDelete || confirmDelete.type !== 'note' || !promptId) return;
    
    try {
      await deleteNote(promptId, confirmDelete.item.id);
      if (onUpdate) {
        await onUpdate();
      }
    } catch (error) {
      console.error('Error deleting note:', error);
      alert('Failed to delete note. Please try again.');
    }
    
    setConfirmDelete(null);
  };
  // Show empty state if no notes
  if (!notes || notes.length === 0) {
    return (
      <p className="no-notes">
        No notes yet. Add your first annotation above!
      </p>
    );
  }

  // Format date for display
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <>
      <ul className="notes-list">
        {notes.map((note) => (
          <li key={note.id} className="note-item">
            <div className="note-content-wrapper">
              <p className="note-content">{note.content}</p>
              <div className="action-icons">
                <button 
                  className="icon-button" 
                  onClick={() => setEditingNote(note)}
                  title="Edit note"
                >
                  ✏️
                </button>
                <button 
                  className="icon-button icon-button-danger" 
                  onClick={() => handleDeleteNote(note)}
                  title="Delete note"
                >
                  🗑️
                </button>
              </div>
            </div>
            <span className="note-date">{formatDate(note.created_at)}</span>
          </li>
        ))}
      </ul>

      {editingNote && (
        <NoteEditModal
          note={editingNote}
          onClose={() => setEditingNote(null)}
          onSave={handleUpdateNote}
        />
      )}

      {/* Confirmation Modal */}
      <ConfirmModal
        isOpen={confirmDelete !== null}
        title="Confirm Delete"
        message={confirmDelete?.message || ''}
        onConfirm={confirmDeleteNote}
        onCancel={() => setConfirmDelete(null)}
        confirmText="Delete"
        cancelText="Cancel"
      />
    </>
  );
}

export default NotesList;