import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const getTree = async () => {
  const response = await api.get('/tree');
  return response.data;
};

export const exportTree = async () => {
  const response = await api.get('/tree/export');
  return response.data;
};

export const importTree = async (treeData) => {
  const response = await api.post('/tree/import', { tree: treeData });
  return response.data;
};

export const saveTree = async (name) => {
  const response = await api.post('/tree/save', { name });
  return response.data;
};

export const listSavedTrees = async () => {
  const response = await api.get('/tree/saves');
  return response.data;
};

export const loadTree = async (name) => {
  const response = await api.post(`/tree/load/${encodeURIComponent(name)}`);
  return response.data;
};

export const deleteSavedTree = async (name) => {
  await api.delete(`/tree/saves/${encodeURIComponent(name)}`);
};

export const getPrompt = async (id) => {
  const response = await api.get(`/prompts/${id}`);
  return response.data;
};

export const createPrompt = async (prompt) => {
  const response = await api.post('/prompts/0', prompt);
  return response.data;
};

export const updatePrompt = async (promptId, prompt) => {
  const response = await api.put(`/prompts/${promptId}`, prompt);
  return response.data;
};

export const deletePrompt = async (promptId) => {
  await api.delete(`/prompts/${promptId}`);
};

export const getPromptNodes = async (id) => {
  const response = await api.get(`/prompts/${id}/nodes`);
  return response.data;
};

export const createNode = async (promptId, node) => {
  const response = await api.post(`/prompts/${promptId}/nodes`, node);
  return response.data;
};

export const updateNode = async (promptId, nodeId, node) => {
  const response = await api.put(`/prompts/${promptId}/nodes/${nodeId}`, node);
  return response.data;
};

export const deleteNode = async (promptId, nodeId) => {
  await api.delete(`/prompts/${promptId}/nodes/${nodeId}`);
};

export const getNotes = async (promptId) => {
  const response = await api.get(`/prompts/${promptId}/notes`);
  return response.data;
};

export const createNote = async (promptId, content) => {
  const response = await api.post(`/prompts/${promptId}/notes`, { content });
  return response.data;
};

export const updateNote = async (promptId, noteId, content) => {
  const response = await api.put(`/prompts/${promptId}/notes/${noteId}`, { content });
  return response.data;
};

export const deleteNote = async (promptId, noteId) => {
  await api.delete(`/prompts/${promptId}/notes/${noteId}`);
};

export default api;