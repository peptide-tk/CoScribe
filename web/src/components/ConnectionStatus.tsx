import React from "react";

interface ConnectionStatusProps {
  connected: boolean;
}

export const ConnectionStatus: React.FC<ConnectionStatusProps> = ({
  connected,
}) => {
  return (
    <div className="connection-status">
      WebSocket:{" "}
      <span className={connected ? "connected" : "disconnected"}>
        {connected ? "Connected" : "Disconnected"}
      </span>
    </div>
  );
};
