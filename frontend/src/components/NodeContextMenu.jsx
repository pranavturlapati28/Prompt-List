import { useEffect, useRef } from 'react';
import './NodeContextMenu.css';

/**
 * Context Menu Component for Nodes
 * Displays a right-click menu with options to edit or delete a node
 */
function NodeContextMenu({ x, y, onClose, onEdit, onDelete }) {
  const menuRef = useRef(null);

  // Close menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (menuRef.current && !menuRef.current.contains(event.target)) {
        onClose();
      }
    };

    // Close on escape key
    const handleEscape = (event) => {
      if (event.key === 'Escape') {
        onClose();
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    document.addEventListener('keydown', handleEscape);

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
      document.removeEventListener('keydown', handleEscape);
    };
  }, [onClose]);

  // Adjust position if menu would go off screen
  useEffect(() => {
    if (menuRef.current) {
      const rect = menuRef.current.getBoundingClientRect();
      const windowWidth = window.innerWidth;
      const windowHeight = window.innerHeight;

      if (rect.right > windowWidth) {
        menuRef.current.style.left = `${x - rect.width}px`;
      }
      if (rect.bottom > windowHeight) {
        menuRef.current.style.top = `${y - rect.height}px`;
      }
    }
  }, [x, y]);

  return (
    <div
      ref={menuRef}
      className="node-context-menu"
      style={{ left: `${x}px`, top: `${y}px` }}
    >
      <button
        className="context-menu-item"
        onClick={() => {
          onEdit();
          onClose();
        }}
      >
        <span className="context-menu-icon"></span>
        Edit
      </button>
      <button
        className="context-menu-item context-menu-item-danger"
        onClick={() => {
          onDelete();
          onClose();
        }}
      >
        <span className="context-menu-icon"></span>
        Delete
      </button>
    </div>
  );
}

export default NodeContextMenu;

