import axios from 'axios';

// Get API URL from environment variable, fallback to localhost
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// Create axios instance with default config
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// =============================================================================
// TREE ENDPOINTS
// =============================================================================

/**
 * GET /tree - Fetch the complete prompt tree
 * Used to render the main tree visualization
 */
export const getTree = async () => {
  const response = await api.get('/tree');
  return response.data;
};

// =============================================================================
// PROMPT ENDPOINTS
// =============================================================================

/**
 * GET /prompts/:id - Fetch a single prompt
 * Used when clicking on a prompt to show details
 */
export const getPrompt = async (id) => {
  const response = await api.get(`/prompts/${id}`);
  return response.data;
};

/**
 * POST /prompts/:id - Create a new prompt
 */
export const createPrompt = async (prompt) => {
  const response = await api.post('/prompts/0', prompt);
  return response.data;
};

/**
 * PUT /prompts/:id - Update a prompt
 */
export const updatePrompt = async (promptId, prompt) => {
  const response = await api.put(`/prompts/${promptId}`, prompt);
  return response.data;
};

/**
 * DELETE /prompts/:id - Delete a prompt
 */
export const deletePrompt = async (promptId) => {
  await api.delete(`/prompts/${promptId}`);
};

// =============================================================================
// NODE ENDPOINTS
// =============================================================================

/**
 * GET /prompts/:id/nodes - Fetch nodes for a prompt
 */
export const getPromptNodes = async (id) => {
  const response = await api.get(`/prompts/${id}/nodes`);
  return response.data;
};

/**
 * POST /prompts/:id/nodes - Create a new node
 */
export const createNode = async (promptId, node) => {
  const response = await api.post(`/prompts/${promptId}/nodes`, node);
  return response.data;
};

/**
 * PUT /prompts/:id/nodes/:nodeId - Update a node
 */
export const updateNode = async (promptId, nodeId, node) => {
  const response = await api.put(`/prompts/${promptId}/nodes/${nodeId}`, node);
  return response.data;
};

/**
 * DELETE /prompts/:id/nodes/:nodeId - Delete a node
 */
export const deleteNode = async (promptId, nodeId) => {
  await api.delete(`/prompts/${promptId}/nodes/${nodeId}`);
};

// =============================================================================
// NOTE ENDPOINTS
// =============================================================================

/**
 * GET /prompts/:id/notes - Fetch notes for a prompt
 * Used to display user annotations in the side panel
 */
export const getNotes = async (promptId) => {
  const response = await api.get(`/prompts/${promptId}/notes`);
  return response.data;
};

/**
 * POST /prompts/:id/notes - Create a new note
 * Used when user adds an annotation
 */
export const createNote = async (promptId, content) => {
  const response = await api.post(`/prompts/${promptId}/notes`, { content });
  return response.data;
};

/**
 * PUT /prompts/:id/notes/:noteId - Update a note
 */
export const updateNote = async (promptId, noteId, content) => {
  const response = await api.put(`/prompts/${promptId}/notes/${noteId}`, { content });
  return response.data;
};

/**
 * DELETE /prompts/:id/notes/:noteId - Delete a note
 */
export const deleteNote = async (promptId, noteId) => {
  await api.delete(`/prompts/${promptId}/notes/${noteId}`);
};

export default api;