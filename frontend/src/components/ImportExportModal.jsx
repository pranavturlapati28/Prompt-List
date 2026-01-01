import { useState } from 'react';
import './NodeEditModal.css'; // Reuse modal styles

/**
 * Modal Component for Importing/Exporting Tree JSON
 */
function ImportExportModal({ isOpen, mode, treeData, onClose, onImport, onExport }) {
  const [jsonText, setJsonText] = useState('');
  const [error, setError] = useState('');

  const handleExport = () => {
    try {
      const json = JSON.stringify(treeData, null, 2);
      setJsonText(json);
      // Copy to clipboard
      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(json).then(() => {
          alert('Tree JSON copied to clipboard!');
        }).catch(() => {
          // Fallback: select text
          const textarea = document.createElement('textarea');
          textarea.value = json;
          document.body.appendChild(textarea);
          textarea.select();
          document.execCommand('copy');
          document.body.removeChild(textarea);
          alert('Tree JSON copied to clipboard!');
        });
      } else {
        // Fallback for older browsers
        const textarea = document.createElement('textarea');
        textarea.value = json;
        document.body.appendChild(textarea);
        textarea.select();
        document.execCommand('copy');
        document.body.removeChild(textarea);
        alert('Tree JSON copied to clipboard!');
      }
    } catch (err) {
      setError('Failed to export tree: ' + err.message);
    }
  };

  const handleImport = () => {
    setError('');
    try {
      const parsed = JSON.parse(jsonText);
      // Validate structure
      if (!parsed.project || !parsed.prompts || !Array.isArray(parsed.prompts)) {
        throw new Error('Invalid tree structure. Expected: { project, mainRequest, prompts: [...] }');
      }
      
      // Transform the tree structure to match backend expectations
      // Convert "subprompts" to "nodes" if present, and ensure proper structure
      const transformed = {
        project: parsed.project,
        mainRequest: parsed.mainRequest || '',
        prompts: parsed.prompts.map(prompt => {
          // Handle both "nodes" and "subprompts" field names
          const nodes = prompt.nodes || prompt.subprompts || [];
          
          return {
            id: prompt.id || 0,
            title: prompt.title || '',
            description: prompt.description || '',
            nodes: nodes.map(node => ({
              id: node.id || 0,
              name: node.name || '',
              action: node.action || ''
            }))
          };
        })
      };
      
      onImport(transformed);
      onClose();
    } catch (err) {
      setError('Invalid JSON: ' + err.message);
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Escape') {
      onClose();
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()} onKeyDown={handleKeyDown}>
        <div className="modal-header">
          <h2>{mode === 'export' ? 'Export Tree' : 'Import Tree'}</h2>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>
        <div className="modal-body">
          {mode === 'export' ? (
            <div>
              <p style={{ marginBottom: '12px', color: 'var(--text-secondary)', fontSize: '14px' }}>
                Click "Copy JSON" to copy the current tree to your clipboard, or copy the text below.
              </p>
              <div className="form-group">
                <label htmlFor="json-text">Tree JSON</label>
                <textarea
                  id="json-text"
                  value={jsonText || JSON.stringify(treeData, null, 2)}
                  onChange={(e) => setJsonText(e.target.value)}
                  rows="15"
                  readOnly={!jsonText}
                  style={{ fontFamily: 'monospace', fontSize: '12px' }}
                />
              </div>
              {error && <p style={{ color: '#ef4444', fontSize: '13px', marginTop: '8px' }}>{error}</p>}
            </div>
          ) : (
            <div>
              <p style={{ marginBottom: '12px', color: 'var(--text-secondary)', fontSize: '14px' }}>
                Paste your tree JSON below. This will replace the current tree.
              </p>
              <div className="form-group">
                <label htmlFor="json-text">Tree JSON</label>
                <textarea
                  id="json-text"
                  value={jsonText}
                  onChange={(e) => setJsonText(e.target.value)}
                  placeholder='{"project": "...", "mainRequest": "...", "prompts": [...]}'
                  rows="15"
                  style={{ fontFamily: 'monospace', fontSize: '12px' }}
                  autoFocus
                />
              </div>
              {error && <p style={{ color: '#ef4444', fontSize: '13px', marginTop: '8px' }}>{error}</p>}
            </div>
          )}
        </div>
        <div className="modal-footer">
          <button type="button" className="btn-secondary" onClick={onClose}>
            Cancel
          </button>
          {mode === 'export' ? (
            <button type="button" className="btn-primary" onClick={handleExport}>
              Copy JSON
            </button>
          ) : (
            <button type="button" className="btn-primary" onClick={handleImport}>
              Import
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

export default ImportExportModal;

