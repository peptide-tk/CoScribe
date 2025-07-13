package document

import (
	"sync"
	"time"
)

type Store interface {
	GetDocument(id string) (*Document, error)
	CreateDocument(id, title string) (*Document, error)
	UpdateDocument(doc *Document) error
	SaveEdit(docID string, edit *Edit) error
	ListDocuments() ([]*DocumentInfo, error)
}

type Manager struct {
	documents map[string]*Document
	store     Store
	mu        sync.RWMutex
}

func NewManager(store Store) *Manager {
	return &Manager{
		documents: make(map[string]*Document),
		store:     store,
	}
}

func (m *Manager) GetDocument(id string) *Document {
	m.mu.Lock()
	defer m.mu.Unlock()

	if doc, exists := m.documents[id]; exists {
		return doc
	}

	if m.store != nil {
		doc, err := m.store.GetDocument(id)
		if err == nil {
			m.documents[id] = doc
			return doc
		}
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
	
	err := doc.ApplyEdit(edit)
	if err != nil {
		return err
	}
	
	if m.store != nil {
		if err := m.store.SaveEdit(docID, edit); err != nil {
			// TODO: ログ機能を実装後に適切にログ出力
		}
		
		// 文書を更新
		if err := m.store.UpdateDocument(doc); err != nil {
			// TODO: ログ機能を実装後に適切にログ出力
		}
	}
	
	return nil
}

func (m *Manager) SaveDocument(docID string) error {
	doc := m.GetDocument(docID)
	if m.store != nil {
		return m.store.UpdateDocument(doc)
	}
	return nil
}

type DocumentInfo struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Version int    `json:"version"`
	Lines   int    `json:"lines"`
}