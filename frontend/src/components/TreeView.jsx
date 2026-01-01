import { useState, useRef, useEffect } from 'react';
import './TreeView.css';

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
function TreeView({ prompts, projectName, onSelectPrompt, selectedPromptId }) {
  const [expandedIds, setExpandedIds] = useState(
    new Set(prompts.map(p => p.id))
  );

  // Track the start index for each prompt's visible nodes
  const [nodeStartIndex, setNodeStartIndex] = useState({});

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
                >
                  {prompt.id}
                </div>

                {/* Prompt Title */}
                <div
                  className="prompt-title"
                  onClick={() => handlePromptClick(prompt)}
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
                        <div key={node.id || nodeIndex} className="node-item">
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
    </div>
  );
}

export default TreeView;