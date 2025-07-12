import React, { useState, useEffect } from "react";
import "./App.css";

interface DocumentState {
  id: string;
  title: string;
  content: string;
  version: number;
}

function App() {
  const [, setWs] = useState<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const [document, setDocument] = useState<DocumentState>({
    id: "sample-doc",
    title: "Loading...",
    content: "",
    version: 0,
  });

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
        </div>

        <div className="document-editor">
          <textarea
            value={document.content}
            readOnly
            placeholder="仮内容"
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
