import React, { useState, useRef, useCallback, useEffect } from "react";
import "./App.css";
import { useWebSocket } from "./hooks/useWebSocket";
import { DocumentState, SaveDocumentResponse } from "./types";

function App() {
  const [document, setDocument] = useState<DocumentState>({
    id: "sample-doc",
    title: "Loading...",
    content: "",
    version: 0,
  });
  const [userId] = useState(`user-${Date.now()}`);
  const [isSaving, setIsSaving] = useState(false);
  const [lastSaved, setLastSaved] = useState<Date | null>(null);
  const autoSaveTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const { connected, sendRealtimeUpdate } = useWebSocket({
    documentId: document.id,
    userId,
    onDocumentUpdate: (updatedDoc) => {
      setDocument(updatedDoc);
    },
    onError: (error) => {
      console.error("WebSocket error:", error);
    },
  });

  useEffect(() => {
    return () => {
      if (autoSaveTimeoutRef.current) {
        clearTimeout(autoSaveTimeoutRef.current);
      }
    };
  }, []);

  const handleTextChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newContent = event.target.value;

    setDocument((prev) => ({
      ...prev,
      content: newContent,
    }));

    sendRealtimeUpdate(newContent);
    scheduleAutoSave();
  };

  const saveDocument = useCallback(
    async (isAutoSave = false) => {
      if (isSaving) return;

      setIsSaving(true);

      try {
        const response = await fetch(
          `http://localhost:8080/api/document/${document.id}`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              content: document.content,
            }),
          }
        );

        const result: SaveDocumentResponse = await response.json();

        if (response.ok) {
          setDocument((prev) => ({
            ...prev,
            version: result.version,
          }));
          setLastSaved(new Date());
        } else {
          console.error("Failed to save document:", result);
        }
      } catch (error) {
        console.error("Save error:", error);
      } finally {
        setIsSaving(false);
      }
    },
    [document.id, document.content, isSaving]
  );

  const scheduleAutoSave = useCallback(() => {
    if (autoSaveTimeoutRef.current) {
      clearTimeout(autoSaveTimeoutRef.current);
    }

    autoSaveTimeoutRef.current = setTimeout(() => {
      saveDocument(true);
    }, 2000);
  }, [saveDocument]);


  return (
    <div className="App">
      <header className="app-header">
        <h1>CoScribe - Collaborative Writing Tool</h1>
        <div className="connection-status">
          WebSocket:{" "}
          <span className={connected ? "connected" : "disconnected"}>
            {connected ? "🔗 Connected (Real-time)" : "🔌 Disconnected"}
          </span>
        </div>
      </header>

      <main className="app-main">
        <div className="document-info">
          <h2>{document.title}</h2>
          <p>Version: {document.version}</p>
          <p>Content length: {document.content.length} characters</p>
          <div className="save-status">
            {isSaving ? (
              <span style={{ color: "#ff9800" }}>💾 Saving...</span>
            ) : lastSaved ? (
              <span style={{ color: "#4caf50" }}>
                ✅ Saved at {lastSaved.toLocaleTimeString()}
              </span>
            ) : (
              <span style={{ color: "#666" }}>💭 Auto-save enabled</span>
            )}
          </div>
          <button onClick={() => saveDocument(false)} disabled={isSaving}>
            {isSaving ? "Saving..." : "Save Now"}
          </button>
        </div>

        <div className="document-editor">
          <textarea
            value={document.content}
            onChange={handleTextChange}
            placeholder="Start writing..."
            rows={20}
            cols={80}
          />
        </div>
      </main>

      <footer className="app-footer">
        <p>CoScribe v1.0.0 - Real-time collaborative writing</p>
      </footer>
    </div>
  );
}

export default App;
