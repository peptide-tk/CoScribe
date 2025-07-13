import React from "react";
import { DocumentState } from "../types";

interface DocumentInfoProps {
  document: DocumentState;
  isSaving: boolean;
  lastSaved: Date | null;
  onSave: () => void;
}

export const DocumentInfo: React.FC<DocumentInfoProps> = ({
  document,
  isSaving,
  lastSaved,
  onSave,
}) => {
  return (
    <div className="document-info">
      <h2>{document.title}</h2>
      <p>Version: {document.version}</p>
      <p>Content length: {document.content.length} characters</p>
      <div className="save-status">
        {isSaving ? (
          <span style={{ color: "#ff9800" }}>Saving...</span>
        ) : lastSaved ? (
          <span style={{ color: "#4caf50" }}>
            Saved at {lastSaved.toLocaleTimeString()}
          </span>
        ) : (
          <span style={{ color: "#666" }}>Auto-save enabled</span>
        )}
      </div>
      <button onClick={onSave} disabled={isSaving}>
        {isSaving ? "Saving..." : "Save Now"}
      </button>
    </div>
  );
};
