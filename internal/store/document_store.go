package store

import (
	"database/sql"
	"fmt"
	"time"

	"coscribe/internal/database"
	"coscribe/internal/document"
)

type DocumentStore struct {
	db *sql.DB
}

func NewDocumentStore() *DocumentStore {
	return &DocumentStore{
		db: database.DB,
	}
}

func (s *DocumentStore) GetDocument(id string) (*document.Document, error) {
	query := `
		SELECT id, title, content, version, created_at, updated_at
		FROM documents WHERE id = $1
	`
	
	var doc document.Document
	var createdAt, updatedAt time.Time
	
	err := s.db.QueryRow(query, id).Scan(
		&doc.ID,
		&doc.Title,
		&doc.Content,
		&doc.Version,
		&createdAt,
		&updatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return s.CreateDocument(id, "Untitled Document")
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	
	doc.Lines = document.ContentToLines(doc.Content)
	
	return &doc, nil
}

func (s *DocumentStore) CreateDocument(id, title string) (*document.Document, error) {
	query := `
		INSERT INTO documents (id, title, content, version)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, content, version, created_at, updated_at
	`
	
	var doc document.Document
	var createdAt, updatedAt time.Time
	
	err := s.db.QueryRow(query, id, title, "", 0).Scan(
		&doc.ID,
		&doc.Title,
		&doc.Content,
		&doc.Version,
		&createdAt,
		&updatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}
	
	doc.Lines = []string{""}
	return &doc, nil
}

func (s *DocumentStore) UpdateDocument(doc *document.Document) error {
	query := `
		UPDATE documents 
		SET title = $2, content = $3, version = $4, updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := s.db.Exec(query, doc.ID, doc.Title, doc.Content, doc.Version)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	
	return nil
}

func (s *DocumentStore) SaveEdit(docID string, edit *document.Edit) error {
	query := `
		INSERT INTO document_edits (document_id, edit_type, line_no, column_pos, content, length, version, user_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := s.db.Exec(
		query,
		docID,
		edit.Type,
		edit.LineNo,
		edit.Column,
		edit.Content,
		edit.Length,
		edit.Version,
		edit.User,
	)
	
	if err != nil {
		return fmt.Errorf("failed to save edit: %w", err)
	}
	
	return nil
}

func (s *DocumentStore) ListDocuments() ([]*document.DocumentInfo, error) {
	query := `
		SELECT id, title, version, 
		       LENGTH(content) - LENGTH(REPLACE(content, E'\n', '')) + 1 as lines,
		       updated_at
		FROM documents 
		ORDER BY updated_at DESC
	`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}
	defer rows.Close()
	
	var docs []*document.DocumentInfo
	for rows.Next() {
		var doc document.DocumentInfo
		var updatedAt time.Time
		
		err := rows.Scan(&doc.ID, &doc.Title, &doc.Version, &doc.Lines, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		
		docs = append(docs, &doc)
	}
	
	return docs, nil
}