import { useState, useEffect } from 'react';
import './NodeEditModal.css'; // Reuse the same styles

/**
 * Modal Component for Editing Prompts
 */
function PromptEditModal({ prompt, onClose, onSave }) {
  const [title, setTitle] = useState(prompt?.title || '');
  const [description, setDescription] = useState(prompt?.description || '');

  useEffect(() => {
    if (prompt) {
      setTitle(prompt.title || '');
      setDescription(prompt.description || '');
    }
  }, [prompt]);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (title.trim()) {
      onSave({ title: title.trim(), description: description.trim() });
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Escape') {
      onClose();
    }
  };

  if (!prompt) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()} onKeyDown={handleKeyDown}>
        <div className="modal-header">
          <h2>Edit Prompt</h2>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="modal-body">
            <div className="form-group">
              <label htmlFor="prompt-title">Title</label>
              <input
                id="prompt-title"
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Prompt title"
                required
                autoFocus
              />
            </div>
            <div className="form-group">
              <label htmlFor="prompt-description">Description</label>
              <textarea
                id="prompt-description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Prompt description"
                rows="4"
              />
            </div>
          </div>
          <div className="modal-footer">
            <button type="button" className="btn-secondary" onClick={onClose}>
              Cancel
            </button>
            <button type="submit" className="btn-primary">
              Save
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default PromptEditModal;

