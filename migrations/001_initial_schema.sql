-- CoScribe Database Schema
-- Initial migration for document collaborative editing

-- Documents table
CREATE TABLE documents (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(500) NOT NULL DEFAULT 'Untitled Document',
    content TEXT NOT NULL DEFAULT '',
    version INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Document edits history table
CREATE TABLE document_edits (
    id SERIAL PRIMARY KEY,
    document_id VARCHAR(255) NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    edit_type VARCHAR(50) NOT NULL, -- 'insert', 'delete', 'replace'
    line_no INTEGER NOT NULL,
    column_pos INTEGER NOT NULL,
    content TEXT,
    length INTEGER,
    version INTEGER NOT NULL,
    user_name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Users table (for future authentication)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Document collaborators table
CREATE TABLE document_collaborators (
    id SERIAL PRIMARY KEY,
    document_id VARCHAR(255) NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'editor', -- 'owner', 'editor', 'viewer'
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_documents_updated_at ON documents(updated_at);
CREATE INDEX idx_document_edits_document_id ON document_edits(document_id);
CREATE INDEX idx_document_edits_version ON document_edits(document_id, version);
CREATE INDEX idx_document_collaborators_document_id ON document_collaborators(document_id);
CREATE INDEX idx_document_collaborators_user_id ON document_collaborators(user_id);

-- Insert sample data
INSERT INTO documents (id, title, content, version) VALUES 
('sample-doc', 'Sample Document', 'Welcome to CoScribe!\nThis is a collaborative writing tool.', 1),
('test-doc', 'Test Document', '', 0);