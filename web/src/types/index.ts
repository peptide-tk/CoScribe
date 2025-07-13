export interface DocumentState {
  id: string;
  title: string;
  content: string;
  version: number;
}

export interface BaseMessage {
  type: string;
  document: string;
  user: string;
  time: string;
}

export interface DocumentUpdateMessage extends BaseMessage {
  type: "document_update";
  content: string;
}

export interface DocumentStateMessage extends BaseMessage {
  type: "document_state";
  content: string;
  version: number;
}

export interface DocumentUpdatedMessage extends BaseMessage {
  type: "document_updated";
  content: string;
  version: number;
}

export interface ErrorMessage extends BaseMessage {
  type: "error";
  content: string;
}

export interface RequestDocumentMessage extends BaseMessage {
  type: "request_document";
}

export type WebSocketMessage =
  | DocumentUpdateMessage
  | DocumentStateMessage
  | DocumentUpdatedMessage
  | ErrorMessage
  | RequestDocumentMessage;

export interface UseWebSocketOptions {
  documentId: string;
  userId: string;
  onDocumentUpdate?: (document: DocumentState) => void;
  onError?: (error: string) => void;
}

export interface UseWebSocketReturn {
  ws: WebSocket | null;
  connected: boolean;
  sendRealtimeUpdate: (content: string) => void;
  reconnect: () => void;
}

export interface SaveDocumentResponse {
  success: boolean;
  version: number;
  content: string;
}
