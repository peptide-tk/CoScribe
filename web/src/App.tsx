import React, { useState, useRef, useCallback, useEffect } from "react";
import "./App.css";
import { useWebSocket } from "./hooks/useWebSocket";
import { DocumentState, SaveDocumentResponse } from "./types";
import { Layout, DocumentInfo, DocumentEditor, ConnectionStatus } from "./components";

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
    <Layout>
      <ConnectionStatus connected={connected} />
      <DocumentInfo
        document={document}
        isSaving={isSaving}
        lastSaved={lastSaved}
        onSave={() => saveDocument(false)}
      />
      <DocumentEditor content={document.content} onChange={handleTextChange} />
    </Layout>
  );
}

export default App;
