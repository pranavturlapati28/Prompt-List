import { useState, useEffect } from 'react';
import './NodeEditModal.css';

/**
 * Modal Component for Editing Nodes
 */
function NodeEditModal({ node, onClose, onSave }) {
  const [name, setName] = useState(node?.name || '');
  const [action, setAction] = useState(node?.action || '');

  useEffect(() => {
    if (node) {
      setName(node.name || '');
      setAction(node.action || '');
    } else {
      // Reset for creating new node
      setName('');
      setAction('');
    }
  }, [node]);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (name.trim()) {
      onSave({ name: name.trim(), action: action.trim() });
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Escape') {
      onClose();
    }
  };

  const isCreating = !node;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()} onKeyDown={handleKeyDown}>
        <div className="modal-header">
          <h2>{isCreating ? 'Add New Subprompt' : 'Edit Node'}</h2>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="modal-body">
            <div className="form-group">
              <label htmlFor="node-name">Name</label>
              <input
                id="node-name"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Node name"
                required
                autoFocus
              />
            </div>
            <div className="form-group">
              <label htmlFor="node-action">Action</label>
              <textarea
                id="node-action"
                value={action}
                onChange={(e) => setAction(e.target.value)}
                placeholder="Action description"
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

export default NodeEditModal;

