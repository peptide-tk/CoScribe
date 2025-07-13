package document

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Document struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Version int       `json:"version"`
	Lines   []string  `json:"lines"`
	mu      sync.RWMutex
}

type Edit struct {
	Type    string    `json:"type"`
	LineNo  int       `json:"line_no"`
	Column  int       `json:"column"`
	Content string    `json:"content"`
	Length  int       `json:"length"`
	Version int       `json:"version"` 
	Time    time.Time `json:"time"`
	User    string    `json:"user"`
}

func NewDocument(id, title string) *Document {
	return &Document{
		ID:      id,
		Title:   title,
		Content: "",
		Version: 0,
		Lines:   []string{""},
	}
}

func (d *Document) GetContent() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Content
}

func (d *Document) GetLines() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return append([]string{}, d.Lines...)
}

func (d *Document) GetVersion() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Version
}

func (d *Document) SetContent(content string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Content = content
	d.Lines = ContentToLines(content)
	d.Version++
}

func (d *Document) ApplyEdit(edit *Edit) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if edit.Version != d.Version {
		return &VersionConflictError{
			Expected: d.Version,
			Actual:   edit.Version,
		}
	}

	if edit.LineNo < 0 || edit.LineNo >= len(d.Lines) {
		return &InvalidLineError{
			LineNo:    edit.LineNo,
			MaxLines:  len(d.Lines),
		}
	}

	switch edit.Type {
	case "insert":
		d.applyInsert(edit)
	case "delete":
		d.applyDelete(edit)
	case "replace":
		d.applyReplace(edit)
	default:
		return &InvalidEditTypeError{Type: edit.Type}
	}

	d.Version++
	d.updateContent()
	
	return nil
}

func (d *Document) applyInsert(edit *Edit) {
	line := d.Lines[edit.LineNo]

	if strings.Contains(edit.Content, "\n") {
		newLines := strings.Split(edit.Content, "\n")
		
		beforeText := line[:edit.Column]
		afterText := line[edit.Column:]
		
		result := make([]string, 0)
		result = append(result, d.Lines[:edit.LineNo]...)
		result = append(result, beforeText+newLines[0])
		
		for i := 1; i < len(newLines)-1; i++ {
			result = append(result, newLines[i])
		}
		
		result = append(result, newLines[len(newLines)-1]+afterText)
		result = append(result, d.Lines[edit.LineNo+1:]...)
		
		d.Lines = result
	} else {
		newLine := line[:edit.Column] + edit.Content + line[edit.Column:]
		d.Lines[edit.LineNo] = newLine
	}
}

func (d *Document) applyDelete(edit *Edit) {
	line := d.Lines[edit.LineNo]
	
	endPos := edit.Column + edit.Length
	if endPos > len(line) {
		endPos = len(line)
	}
	
	newLine := line[:edit.Column] + line[endPos:]
	d.Lines[edit.LineNo] = newLine
}

func (d *Document) applyReplace(edit *Edit) {
	line := d.Lines[edit.LineNo]
	
	endPos := edit.Column + edit.Length
	if endPos > len(line) {
		endPos = len(line)
	}
	
	newLine := line[:edit.Column] + edit.Content + line[endPos:]
	d.Lines[edit.LineNo] = newLine
}

func (d *Document) updateContent() {
	d.Content = strings.Join(d.Lines, "\n")
}

func ContentToLines(content string) []string {
	if content == "" {
		return []string{""}
	}
	return strings.Split(content, "\n")
}

type VersionConflictError struct {
	Expected int
	Actual   int
}

func (e *VersionConflictError) Error() string {
	return fmt.Sprintf("version conflict: expected %d, got %d", e.Expected, e.Actual)
}

type InvalidLineError struct {
	LineNo   int
	MaxLines int
}

func (e *InvalidLineError) Error() string {
	return fmt.Sprintf("invalid line number: %d (max: %d)", e.LineNo, e.MaxLines-1)
}

type InvalidEditTypeError struct {
	Type string
}

func (e *InvalidEditTypeError) Error() string {
	return "invalid edit type: " + e.Type
}