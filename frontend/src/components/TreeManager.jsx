import { useState } from 'react';
import ImportExportModal from './ImportExportModal';
import SaveLoadModal from './SaveLoadModal';
import { exportTree, importTree, saveTree, listSavedTrees, loadTree, deleteSavedTree } from '../api/api';
import './TreeManager.css';

/**
 * TreeManager Component
 * Button and modals for managing tree import/export and save/load
 */
function TreeManager({ treeData, onTreeUpdate }) {
  const [showImportExport, setShowImportExport] = useState(false);
  const [showSaveLoad, setShowSaveLoad] = useState(false);
  const [importExportMode, setImportExportMode] = useState('export'); // 'export' or 'import'
  const [savedTrees, setSavedTrees] = useState([]);
  const [loading, setLoading] = useState(false);

  const handleExportClick = () => {
    setImportExportMode('export');
    setShowImportExport(true);
  };

  const handleImportClick = () => {
    setImportExportMode('import');
    setShowImportExport(true);
  };

  const handleSaveLoadClick = async () => {
    setShowSaveLoad(true);
    await refreshSavedTrees();
  };

  const refreshSavedTrees = async () => {
    try {
      const data = await listSavedTrees();
      setSavedTrees(data.trees || []);
    } catch (error) {
      console.error('Error fetching saved trees:', error);
    }
  };

  const handleImport = async (treeData) => {
    setLoading(true);
    try {
      await importTree(treeData);
      if (onTreeUpdate) {
        await onTreeUpdate();
      }
      alert('Tree imported successfully!');
    } catch (error) {
      console.error('Error importing tree:', error);
      console.error('Error response:', error.response?.data);
      
      // Build detailed error message
      let errorMsg = 'Failed to import tree: ';
      if (error.response?.data) {
        if (error.response.data.detail) {
          errorMsg += error.response.data.detail;
        } else if (error.response.data.errors && Array.isArray(error.response.data.errors)) {
          errorMsg += error.response.data.errors.map(e => e.message || e).join(', ');
        } else {
          errorMsg += JSON.stringify(error.response.data);
        }
      } else {
        errorMsg += error.message;
      }
      
      alert(errorMsg);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (name) => {
    setLoading(true);
    try {
      await saveTree(name);
      await refreshSavedTrees();
      alert('Tree saved successfully!');
    } catch (error) {
      console.error('Error saving tree:', error);
      alert('Failed to save tree: ' + (error.response?.data?.detail || error.message));
    } finally {
      setLoading(false);
    }
  };

  const handleLoad = async (name) => {
    setLoading(true);
    try {
      await loadTree(name);
      if (onTreeUpdate) {
        await onTreeUpdate();
      }
      setShowSaveLoad(false);
      alert('Tree loaded successfully!');
    } catch (error) {
      console.error('Error loading tree:', error);
      alert('Failed to load tree: ' + (error.response?.data?.detail || error.message));
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (name) => {
    setLoading(true);
    try {
      await deleteSavedTree(name);
      await refreshSavedTrees();
    } catch (error) {
      console.error('Error deleting tree:', error);
      alert('Failed to delete tree: ' + (error.response?.data?.detail || error.message));
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <div className="tree-manager">
        <button
          className="tree-manager-button"
          onClick={handleExportClick}
          title="Export tree as JSON"
        >
          Export
        </button>
        <button
          className="tree-manager-button"
          onClick={handleImportClick}
          title="Import tree from JSON"
        >
          Import
        </button>
        <button
          className="tree-manager-button"
          onClick={handleSaveLoadClick}
          title="Save or load a tree"
        >
          Save/Load
        </button>
      </div>

      <ImportExportModal
        isOpen={showImportExport}
        mode={importExportMode}
        treeData={treeData}
        onClose={() => setShowImportExport(false)}
        onImport={handleImport}
        onExport={() => {}}
      />

      <SaveLoadModal
        isOpen={showSaveLoad}
        savedTrees={savedTrees}
        onClose={() => setShowSaveLoad(false)}
        onSave={handleSave}
        onLoad={handleLoad}
        onDelete={handleDelete}
      />
    </>
  );
}

export default TreeManager;

