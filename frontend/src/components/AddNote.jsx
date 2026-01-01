import { useState } from 'react';
import './AddNote.css';

/**
 * AddNote Component
 * 
 * Form for adding a new note to a prompt
 * 
 * Props:
 * - onAdd: Async callback function that receives the note content
 *          Returns true on success, false on failure
 */
function AddNote({ onAdd }) {
  const [content, setContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState(null);

  // Handle form submission
  const handleSubmit = async (e) => {
    e.preventDefault();

    // Validate input
    const trimmedContent = content.trim();
    if (!trimmedContent) {
      setError('Please enter a note');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    // Call parent handler
    const success = await onAdd(trimmedContent);

    if (success) {
      setContent(''); // Clear form on success
    } else {
      setError('Failed to add note. Please try again.');
    }

    setIsSubmitting(false);
  };

  return (
    <form className="add-note-form" onSubmit={handleSubmit}>
      <textarea
        value={content}
        onChange={(e) => {
          setContent(e.target.value);
          setError(null); // Clear error when typing
        }}
        placeholder="Add a note or annotation..."
        rows={3}
        disabled={isSubmitting}
      />
      
      {error && <span className="form-error">{error}</span>}
      
      <button 
        type="submit" 
        disabled={isSubmitting || !content.trim()}
      >
        {isSubmitting ? 'Adding...' : 'Add Note'}
      </button>
    </form>
  );
}

export default AddNote;