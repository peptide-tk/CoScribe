import { useState, useEffect, useRef, useCallback } from "react";

interface DocumentState {
  id: string;
  title: string;
  content: string;
  version: number;
}

interface UseWebSocketOptions {
  documentId: string;
  userId: string;
  onDocumentUpdate?: (document: DocumentState) => void;
  onError?: (error: string) => void;
}

interface UseWebSocketReturn {
  ws: WebSocket | null;
  connected: boolean;
  sendRealtimeUpdate: (content: string) => void;
  reconnect: () => void;
}

export const useWebSocket = ({
  documentId,
  userId,
  onDocumentUpdate,
  onError,
}: UseWebSocketOptions): UseWebSocketReturn => {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const reconnectTimerRef = useRef<NodeJS.Timeout | null>(null);

  const onDocumentUpdateRef = useRef(onDocumentUpdate);
  const onErrorRef = useRef(onError);

  useEffect(() => {
    onDocumentUpdateRef.current = onDocumentUpdate;
    onErrorRef.current = onError;
  });

  const connectWebSocket = useCallback(() => {
    const websocket = new WebSocket(
      `ws://localhost:8080/ws/document?doc=${documentId}`
    );

    websocket.onopen = () => {
      setConnected(true);
      setWs(websocket);
    };

    websocket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log("Received:", data);

      if (data.type === "document_state") {
        const document: DocumentState = {
          id: data.document,
          title: "Sample Document",
          content: data.content,
          version: data.version,
        };
        onDocumentUpdateRef.current?.(document);
      } else if (data.type === "document_updated") {
        if (data.content !== undefined && data.user !== userId) {
          const document: DocumentState = {
            id: data.document,
            title: "Sample Document",
            content: data.content,
            version: data.version,
          };
          onDocumentUpdateRef.current?.(document);
        }
      } else if (data.type === "error") {
        console.error("Server error:", data.content);
        onErrorRef.current?.(data.content);
        if (websocket && websocket.readyState === WebSocket.OPEN) {
          websocket.send(
            JSON.stringify({
              type: "request_document",
              document: documentId,
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
      setConnected(false);
      setWs(null);

      reconnectTimerRef.current = setTimeout(() => {
        connectWebSocket();
      }, 3000);
    };

    return websocket;
  }, [documentId, userId]);

  const sendRealtimeUpdate = useCallback(
    (content: string) => {
      if (ws && connected) {
        const message = {
          type: "document_update",
          document: documentId,
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
    [ws, connected, documentId, userId]
  );

  const reconnect = useCallback(() => {
    if (reconnectTimerRef.current) {
      clearTimeout(reconnectTimerRef.current);
    }
    if (ws) {
      ws.close();
    }
    connectWebSocket();
  }, [connectWebSocket, ws]);

  useEffect(() => {
    const websocket = connectWebSocket();

    return () => {
      if (reconnectTimerRef.current) {
        clearTimeout(reconnectTimerRef.current);
      }
      websocket.close();
    };
  }, [connectWebSocket]);

  return {
    ws,
    connected,
    sendRealtimeUpdate,
    reconnect,
  };
};
