import React, { useState, useEffect, useRef, useCallback } from "react";
import "./App.css";

interface DocumentState {
  id: string;
  title: string;
  content: string;
  version: number;
}

function App() {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
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

  useEffect(() => {
    let reconnectTimer: NodeJS.Timeout;

    const connectWebSocket = () => {
      const websocket = new WebSocket(
        "ws://localhost:8080/ws/document?doc=sample-doc"
      );

      websocket.onopen = () => {
        setConnected(true);
        setWs(websocket);
      };

      websocket.onmessage = (event) => {
        const data = JSON.parse(event.data);

        if (data.type === "document_state") {
          setDocument({
            id: data.document,
            title: "Sample Document",
            content: data.content,
            version: data.version,
          });
        } else if (data.type === "edit") {
          setDocument((prev) => ({
            ...prev,
            version: data.version,
          }));
        } else if (data.type === "document_updated") {
          if (data.content !== undefined && data.user !== userId) {
            setDocument((prev) => ({
              ...prev,
              content: data.content,
              version: data.version,
            }));
          } else {
            setDocument((prev) => ({
              ...prev,
              version: data.version,
            }));
          }
        } else if (data.type === "error") {
          console.error("Server error:", data.content);
          if (ws && connected) {
            ws.send(
              JSON.stringify({
                type: "request_document",
                document: document.id,
                user: userId,
              })
            );
          }
        }
      };

      websocket.onerror = (error) => {
        console.error("WebSocket error:", error);
      };

      websocket.onclose = () => {
        console.log("🔌 WebSocket disconnected");
        setConnected(false);
        setWs(null);

        reconnectTimer = setTimeout(() => {
          connectWebSocket();
        }, 3000);
      };

      return websocket;
    };

    const websocket = connectWebSocket();

    return () => {
      if (reconnectTimer) {
        clearTimeout(reconnectTimer);
      }
      websocket.close();
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

        const result = await response.json();

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

  const sendRealtimeUpdate = useCallback(
    (content: string) => {
      if (ws && connected) {
        const message = {
          type: "document_update",
          document: document.id,
          content: content,
          user: userId,
          time: new Date().toISOString(),
        };
        try {
          ws.send(JSON.stringify(message));
        } catch (error) {
          console.error("Failed to send message:", error);
        }
      } else {
        console.log("Cannot send: WebSocket not connected", {
          ws: !!ws,
          connected,
        });
      }
    },
    [ws, connected, document.id, userId]
  );

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
