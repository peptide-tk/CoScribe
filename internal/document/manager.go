package document

import (
	"sync"
	"time"
)

type Manager struct {
	documents map[string]*Document
	mu        sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		documents: make(map[string]*Document),
	}
}

func (m *Manager) GetDocument(id string) *Document {
	m.mu.Lock()
	defer m.mu.Unlock()

	if doc, exists := m.documents[id]; exists {
		return doc
	}

	doc := NewDocument(id, "Untitled Document")
	m.documents[id] = doc
	return doc
}

func (m *Manager) GetDocumentInfo(id string) *DocumentInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if doc, exists := m.documents[id]; exists {
		return &DocumentInfo{
			ID:      doc.ID,
			Title:   doc.Title,
			Version: doc.GetVersion(),
			Lines:   len(doc.GetLines()),
		}
	}

	return nil
}

func (m *Manager) ListDocuments() []*DocumentInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*DocumentInfo
	for _, doc := range m.documents {
		result = append(result, &DocumentInfo{
			ID:      doc.ID,
			Title:   doc.Title,
			Version: doc.GetVersion(),
			Lines:   len(doc.GetLines()),
		})
	}

	return result
}

func (m *Manager) ApplyEdit(docID string, edit *Edit) error {
	doc := m.GetDocument(docID)
	edit.Time = time.Now()
	return doc.ApplyEdit(edit)
}

type DocumentInfo struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Version int    `json:"version"`
	Lines   int    `json:"lines"`
}