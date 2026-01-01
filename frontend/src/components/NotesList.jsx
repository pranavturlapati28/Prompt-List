import './NotesList.css';

/**
 * NotesList Component
 * 
 * Displays a list of notes/annotations for a prompt
 * 
 * Props:
 * - notes: Array of note objects with id, content, and created_at
 */
function NotesList({ notes }) {
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
    <ul className="notes-list">
      {notes.map((note) => (
        <li key={note.id} className="note-item">
          <p className="note-content">{note.content}</p>
          <span className="note-date">{formatDate(note.created_at)}</span>
        </li>
      ))}
    </ul>
  );
}

export default NotesList;