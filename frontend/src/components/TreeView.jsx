import { useState, useRef, useEffect } from 'react';
import './TreeView.css';
import NodeContextMenu from './NodeContextMenu';
import NodeEditModal from './NodeEditModal';
import PromptEditModal from './PromptEditModal';
import ConfirmModal from './ConfirmModal';
import { updateNode, deleteNode, updatePrompt, deletePrompt } from '../api/api';

/**
 * Generate a color based on index
 */
function getColor(index) {
  const colors = [
    '#6366f1', // indigo
    '#f59e0b', // amber
    '#ef4444', // red
    '#8b5cf6', // violet
    '#ec4899', // pink
    '#06b6d4', // cyan
    '#84cc16', // lime
  ];
  return colors[index % colors.length];
}

/**
 * TreeView Component
 * Displays prompts in a vertical timeline with horizontal node branches
 */
function TreeView({ prompts, projectName, onSelectPrompt, selectedPromptId, onTreeUpdate }) {
  const [expandedIds, setExpandedIds] = useState(
    new Set(prompts.map(p => p.id))
  );

  // Track the start index for each prompt's visible nodes
  const [nodeStartIndex, setNodeStartIndex] = useState({});

  // Context menu state
  const [contextMenu, setContextMenu] = useState(null);
  const [editingNode, setEditingNode] = useState(null);
  const [promptContextMenu, setPromptContextMenu] = useState(null);
  const [editingPrompt, setEditingPrompt] = useState(null);
  
  // Confirmation modal state
  const [confirmDelete, setConfirmDelete] = useState(null);

  // Track container width to calculate max visible nodes
  const [containerWidth, setContainerWidth] = useState(0);
  const containerRef = useRef(null);

  // Measure container width on mount and resize
  useEffect(() => {
    const updateWidth = () => {
      if (containerRef.current) {
        const treeView = containerRef.current.closest('.tree-view');
        if (treeView) {
          // Get the tree-view width and subtract the fixed elements
          // Circle (36px) + gap (8px) + Title (180px) + gap (8px) + branch padding (30px) = 262px
          const fixedWidth = 262;
          const availableWidth = treeView.clientWidth - fixedWidth;
          setContainerWidth(availableWidth);
        }
      }
    };

    updateWidth();
    window.addEventListener('resize', updateWidth);
    return () => window.removeEventListener('resize', updateWidth);
  }, []);

  const toggleExpand = (id) => {
    setExpandedIds(prev => {
      const newSet = new Set(prev);
      if (newSet.has(id)) {
        newSet.delete(id);
      } else {
        newSet.add(id);
      }
      return newSet;
    });
  };

  const handlePromptClick = (prompt) => {
    onSelectPrompt(prompt);
  };

  const handleCircleClick = (prompt, event) => {
    event.stopPropagation();
    toggleExpand(prompt.id);
    onSelectPrompt(prompt);
  };

  const goToNextNodes = (promptId, currentStart, maxVisible) => {
    setNodeStartIndex(prev => ({
      ...prev,
      [promptId]: currentStart + maxVisible
    }));
  };

  const goToPrevNodes = (promptId, currentStart, maxVisible) => {
    setNodeStartIndex(prev => ({
      ...prev,
      [promptId]: Math.max(0, currentStart - maxVisible)
    }));
  };

  // Handle right-click on node
  const handleNodeRightClick = (e, node, promptId) => {
    e.preventDefault();
    e.stopPropagation();
    setContextMenu({
      x: e.clientX,
      y: e.clientY,
      node,
      promptId,
    });
  };

  // Handle node delete
  const handleDeleteNode = () => {
    if (!contextMenu) return;
    const { node, promptId } = contextMenu;
    setConfirmDelete({
      type: 'node',
      item: node,
      promptId,
      message: `Are you sure you want to delete "${node.name}"?`,
    });
    setContextMenu(null);
  };

  // Confirm and execute node delete
  const confirmDeleteNode = async () => {
    if (!confirmDelete || confirmDelete.type !== 'node') return;
    
    try {
      await deleteNode(confirmDelete.promptId, confirmDelete.item.id);
      // Refresh the tree
      if (onTreeUpdate) {
        await onTreeUpdate();
      }
    } catch (error) {
      console.error('Error deleting node:', error);
      alert('Failed to delete node. Please try again.');
    }
    
    setConfirmDelete(null);
  };

  // Handle node edit
  const handleEditNode = () => {
    if (!contextMenu) return;
    setEditingNode({
      node: contextMenu.node,
      promptId: contextMenu.promptId,
    });
    setContextMenu(null);
  };

  // Handle save edited node
  const handleSaveNode = async (updatedData) => {
    if (!editingNode) return;

    try {
      await updateNode(editingNode.promptId, editingNode.node.id, updatedData);
      // Refresh the tree
      if (onTreeUpdate) {
        await onTreeUpdate();
      }
      setEditingNode(null);
    } catch (error) {
      console.error('Error updating node:', error);
      alert('Failed to update node. Please try again.');
    }
  };

  // Handle right-click on prompt
  const handlePromptRightClick = (e, prompt) => {
    e.preventDefault();
    e.stopPropagation();
    setPromptContextMenu({
      x: e.clientX,
      y: e.clientY,
      prompt,
    });
  };

  // Handle prompt delete
  const handleDeletePrompt = () => {
    if (!promptContextMenu) return;
    const { prompt } = promptContextMenu;
    setConfirmDelete({
      type: 'prompt',
      item: prompt,
      message: `Are you sure you want to delete "${prompt.title}"? This will also delete all associated nodes and notes.`,
    });
    setPromptContextMenu(null);
  };

  // Confirm and execute prompt delete
  const confirmDeletePrompt = async () => {
    if (!confirmDelete || confirmDelete.type !== 'prompt') return;
    
    try {
      await deletePrompt(confirmDelete.item.id);
      // Refresh the tree
      if (onTreeUpdate) {
        await onTreeUpdate();
      }
      // Clear selection if deleted prompt was selected
      if (selectedPromptId === confirmDelete.item.id) {
        onSelectPrompt(null);
      }
    } catch (error) {
      console.error('Error deleting prompt:', error);
      alert('Failed to delete prompt. Please try again.');
    }
    
    setConfirmDelete(null);
  };

  // Handle prompt edit
  const handleEditPrompt = () => {
    if (!promptContextMenu) return;
    setEditingPrompt(promptContextMenu.prompt);
    setPromptContextMenu(null);
  };

  // Handle save edited prompt
  const handleSavePrompt = async (updatedData) => {
    if (!editingPrompt) return;

    try {
      await updatePrompt(editingPrompt.id, updatedData);
      // Refresh the tree
      if (onTreeUpdate) {
        await onTreeUpdate();
      }
      setEditingPrompt(null);
    } catch (error) {
      console.error('Error updating prompt:', error);
      alert('Failed to update prompt. Please try again.');
    }
  };

  // Calculate max visible nodes based on container width
  const calculateMaxVisibleNodes = () => {
    if (containerWidth <= 0) return 10; // Default fallback

    // CSS values from TreeView.css
    const gap = 80;           // gap between nodes (.nodes-list gap: 80px)
    const diameter = 12;      // node circle diameter (.node-circle width: 12px)
    const branchPadding = 30; // .branch-container padding-left: 30px
    const nodesMargin = -25;  // .nodes-list margin-left: -25px

    // Effective start position where nodes begin
    const startPosition = branchPadding + nodesMargin; // 30 + (-25) = 5px

    // Formula: (container_length - start_length) / (gap + diameter)
    const maxNodes = Math.floor((containerWidth - startPosition) / (gap + diameter));

    // Ensure at least 5 nodes and at most 15 to keep it reasonable
    return Math.max(5, Math.min(15, maxNodes));
  };

  return (
    <div className="tree-view">
      {/* Header */}
      <div className="tree-header">
        <h1 className="tree-title">{projectName || 'Prompt Tree'}</h1>
      </div>

      {/* Tree Structure */}
      <div className="tree-structure" ref={containerRef}>
        {/* Prompts */}
        {prompts.map((prompt, index) => {
          const color = getColor(index);
          const isExpanded = expandedIds.has(prompt.id);
          const isSelected = selectedPromptId === prompt.id;
          const hasNodes = prompt.nodes && prompt.nodes.length > 0;

          // Calculate visible nodes with pagination based on available width
          const MAX_VISIBLE_NODES = calculateMaxVisibleNodes();
          const startIdx = nodeStartIndex[prompt.id] || 0;
          const hasBack = startIdx > 0;

          // If we have a back arrow, we have one less slot for nodes
          const availableSlots = hasBack ? MAX_VISIBLE_NODES - 1 : MAX_VISIBLE_NODES;
          const endIdx = startIdx + availableSlots;
          const hasNext = hasNodes && endIdx < prompt.nodes.length;

          // If we have a next arrow, we need one less node
          const actualEndIdx = hasNext ? endIdx - 1 : endIdx;
          const visibleNodes = hasNodes ? prompt.nodes.slice(startIdx, actualEndIdx) : [];

          // Calculate total items (back arrow + nodes + next arrow)
          const totalItems = (hasBack ? 1 : 0) + visibleNodes.length + (hasNext ? 1 : 0);

          return (
            <div key={prompt.id} className="prompt-item">
              <div className="prompt-row">
                {/* Numbered Circle */}
                <div
                  className={`prompt-circle ${isSelected ? 'selected' : ''}`}
                  style={{ backgroundColor: color }}
                  onClick={(e) => handleCircleClick(prompt, e)}
                  onContextMenu={(e) => handlePromptRightClick(e, prompt)}
                >
                  {index + 1}
                </div>

                {/* Prompt Title */}
                <div
                  className="prompt-title"
                  onClick={() => handlePromptClick(prompt)}
                  onContextMenu={(e) => handlePromptRightClick(e, prompt)}
                >
                  {prompt.title}
                </div>

                {/* Horizontal Branch with Nodes */}
                {isExpanded && hasNodes && (
                  <div className="branch-container">
                    <div
                      className="nodes-list"
                      style={{
                        color: color,
                        '--line-width': `${80 * (totalItems-0.)}px`
                      }}
                    >
                      {/* Back Arrow */}
                      {hasBack && (
                        <div className="node-item arrow-node">
                          <button
                            className="node-arrow"
                            style={{ borderColor: color, color: color }}
                            onClick={() => goToPrevNodes(prompt.id, startIdx, availableSlots)}
                            aria-label="Previous nodes"
                          >
                            <svg viewBox="0 0 24 24" width="12" height="12">
                              <polyline points="15 18 9 12 15 6" stroke="currentColor" strokeWidth="2" fill="none"/>
                            </svg>
                          </button>
                          <div className="node-label">Back</div>
                        </div>
                      )}

                      {/* Visible Nodes */}
                      {visibleNodes.map((node, nodeIndex) => (
                        <div 
                          key={node.id || nodeIndex} 
                          className="node-item"
                          onContextMenu={(e) => handleNodeRightClick(e, node, prompt.id)}
                        >
                          <div
                            className="node-circle"
                            style={{
                              borderColor: color
                            }}
                          ></div>
                          <div className="node-label">{node.name}</div>
                        </div>
                      ))}

                      {/* Next Arrow */}
                      {hasNext && (
                        <div className="node-item arrow-node">
                          <button
                            className="node-arrow"
                            style={{ borderColor: color, color: color }}
                            onClick={() => goToNextNodes(prompt.id, startIdx, availableSlots)}
                            aria-label="Next nodes"
                          >
                            <svg viewBox="0 0 24 24" width="12" height="12">
                              <polyline points="9 18 15 12 9 6" stroke="currentColor" strokeWidth="2" fill="none"/>
                            </svg>
                          </button>
                          <div className="node-label">Next</div>
                        </div>
                      )}
                    </div>
                  </div>
                )}
              </div>
            </div>
          );
        })}

        {/* Complete Game Footer */}
        <div className="tree-footer">
          <div className="complete-item">
            <div className="complete-circle">
              <svg viewBox="0 0 24 24">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
            </div>
            <div className="complete-title">Complete Game</div>
          </div>
        </div>
      </div>

      {/* Context Menu */}
      {contextMenu && (
        <NodeContextMenu
          x={contextMenu.x}
          y={contextMenu.y}
          onClose={() => setContextMenu(null)}
          onEdit={handleEditNode}
          onDelete={handleDeleteNode}
        />
      )}

      {/* Edit Modal */}
      {editingNode && (
        <NodeEditModal
          node={editingNode.node}
          onClose={() => setEditingNode(null)}
          onSave={handleSaveNode}
        />
      )}

      {/* Prompt Context Menu */}
      {promptContextMenu && (
        <NodeContextMenu
          x={promptContextMenu.x}
          y={promptContextMenu.y}
          onClose={() => setPromptContextMenu(null)}
          onEdit={handleEditPrompt}
          onDelete={handleDeletePrompt}
        />
      )}

      {/* Prompt Edit Modal */}
      {editingPrompt && (
        <PromptEditModal
          prompt={editingPrompt}
          onClose={() => setEditingPrompt(null)}
          onSave={handleSavePrompt}
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

export default TreeView;