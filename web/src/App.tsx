import React, { useState, useEffect } from "react";
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

  useEffect(() => {
    const websocket = new WebSocket(
      "ws://localhost:8080/ws/document?doc=sample-doc"
    );

    websocket.onopen = () => {
      console.log("WebSocket connected");
      setConnected(true);
      setWs(websocket);
    };

    websocket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log("Received:", data);

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
        setDocument((prev) => ({
          ...prev,
          version: data.version,
        }));
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
      console.log("WebSocket disconnected");
      setConnected(false);
      setWs(null);
    };

    return () => {
      websocket.close();
    };
  }, []);

  const handleTextChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newContent = event.target.value;

    setDocument((prev) => ({
      ...prev,
      content: newContent,
    }));

    console.log("Text changed to:", newContent);
  };

  const saveDocument = async () => {
    console.log("Save button clicked!");
    console.log("Document content:", document.content);

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
      } else {
        console.error("Failed to save document:", result);
      }
    } catch (error) {
      console.error("Save error:", error);
    }
  };

  return (
    <div className="App">
      <header className="app-header">
        <h1>CoScribe - Collaborative Writing Tool</h1>
        <div className="connection-status">
          Status:{" "}
          <span className={connected ? "connected" : "disconnected"}>
            {connected ? "Connected" : "Disconnected"}
          </span>
        </div>
      </header>

      <main className="app-main">
        <div className="document-info">
          <h2>{document.title}</h2>
          <p>Version: {document.version}</p>
          <p>Content length: {document.content.length} characters</p>
          <button onClick={saveDocument}>Save Document (HTTP)</button>
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
