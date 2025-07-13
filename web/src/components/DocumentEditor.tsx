import React from "react";

interface DocumentEditorProps {
  content: string;
  onChange: (event: React.ChangeEvent<HTMLTextAreaElement>) => void;
  placeholder?: string;
  rows?: number;
  cols?: number;
}

export const DocumentEditor: React.FC<DocumentEditorProps> = ({
  content,
  onChange,
  placeholder = "Start writing...",
  rows = 20,
  cols = 80,
}) => {
  return (
    <div className="document-editor">
      <textarea
        value={content}
        onChange={onChange}
        placeholder={placeholder}
        rows={rows}
        cols={cols}
      />
    </div>
  );
};