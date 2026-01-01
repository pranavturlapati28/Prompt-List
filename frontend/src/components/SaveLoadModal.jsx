import { useState, useEffect } from 'react';
import './NodeEditModal.css'; // Reuse modal styles

/**
 * Modal Component for Saving/Loading Trees
 */
function SaveLoadModal({ isOpen, savedTrees, onClose, onSave, onLoad, onDelete }) {
  const [saveName, setSaveName] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    if (isOpen) {
      setSaveName('');
      setError('');
    }
  }, [isOpen]);

  const handleSave = () => {
    setError('');
    if (!saveName.trim()) {
      setError('Please enter a name for the saved tree');
      return;
    }
    onSave(saveName.trim());
    setSaveName('');
  };

  const handleLoad = (name) => {
    setError('');
    onLoad(name);
  };

  const handleDelete = (name) => {
    if (window.confirm(`Are you sure you want to delete "${name}"?`)) {
      onDelete(name);
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Escape') {
      onClose();
    }
  };

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

  if (!isOpen) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()} onKeyDown={handleKeyDown}>
        <div className="modal-header">
          <h2>Save & Load Trees</h2>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>
        <div className="modal-body">
          {/* Save Section */}
          <div style={{ marginBottom: '24px' }}>
            <h3 style={{ fontSize: '16px', fontWeight: '600', marginBottom: '12px', color: 'var(--text-primary)' }}>
              Save Current Tree
            </h3>
            <div className="form-group" style={{ marginBottom: '0' }}>
              <label htmlFor="save-name">Name</label>
              <input
                id="save-name"
                type="text"
                value={saveName}
                onChange={(e) => setSaveName(e.target.value)}
                placeholder="Enter a name for this tree"
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    handleSave();
                  }
                }}
              />
            </div>
            <button
              type="button"
              className="btn-primary"
              onClick={handleSave}
              style={{ marginTop: '12px', width: '100%' }}
            >
              Save Tree
            </button>
            {error && <p style={{ color: '#ef4444', fontSize: '13px', marginTop: '8px' }}>{error}</p>}
          </div>

          {/* Saved Trees List */}
          <div>
            <h3 style={{ fontSize: '16px', fontWeight: '600', marginBottom: '12px', color: 'var(--text-primary)' }}>
              Saved Trees ({savedTrees.length})
            </h3>
            {savedTrees.length === 0 ? (
              <p style={{ color: 'var(--text-secondary)', fontSize: '14px', fontStyle: 'italic' }}>
                No saved trees yet. Save a tree above to get started.
              </p>
            ) : (
              <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
                {savedTrees.map((tree) => (
                  <li
                    key={tree.name}
                    style={{
                      padding: '12px',
                      marginBottom: '8px',
                      background: 'var(--bg-tertiary)',
                      border: '1px solid var(--border-color)',
                      borderRadius: '6px',
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                    }}
                  >
                    <div style={{ flex: 1 }}>
                      <div style={{ fontWeight: '500', color: 'var(--text-primary)', marginBottom: '4px' }}>
                        {tree.name}
                      </div>
                      <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>
                        Updated: {formatDate(tree.updated_at)}
                      </div>
                    </div>
                    <div style={{ display: 'flex', gap: '8px' }}>
                      <button
                        type="button"
                        className="btn-primary"
                        onClick={() => handleLoad(tree.name)}
                        style={{ padding: '6px 12px', fontSize: '13px' }}
                      >
                        Load
                      </button>
                      <button
                        type="button"
                        className="btn-secondary"
                        onClick={() => handleDelete(tree.name)}
                        style={{ padding: '6px 12px', fontSize: '13px', color: '#ef4444' }}
                      >
                        Delete
                      </button>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
        <div className="modal-footer">
          <button type="button" className="btn-secondary" onClick={onClose}>
            Close
          </button>
        </div>
      </div>
    </div>
  );
}

export default SaveLoadModal;

