import React from "react";

interface FooterProps {
  version?: string;
}

export const Footer: React.FC<FooterProps> = ({ version = "1.0.0" }) => {
  return (
    <footer className="app-footer">
      <p>CoScribe v{version} - Real-time collaborative writing</p>
    </footer>
  );
};