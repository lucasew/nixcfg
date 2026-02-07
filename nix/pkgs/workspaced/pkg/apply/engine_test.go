package apply

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
)

type MockFileSystem struct {
	Files map[string][]byte
	Infos map[string]os.FileInfo
	Links map[string]string
}

func NewMockFS() *MockFileSystem {
	return &MockFileSystem{
		Files: make(map[string][]byte),
		Infos: make(map[string]os.FileInfo),
		Links: make(map[string]string),
	}
}

type MockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (m MockFileInfo) Name() string       { return m.name }
func (m MockFileInfo) Size() int64        { return m.size }
func (m MockFileInfo) Mode() os.FileMode  { return m.mode }
func (m MockFileInfo) ModTime() time.Time { return m.modTime }
func (m MockFileInfo) IsDir() bool        { return m.isDir }
func (m MockFileInfo) Sys() interface{}   { return nil }

func (fs *MockFileSystem) Lstat(name string) (os.FileInfo, error) {
	if info, ok := fs.Infos[name]; ok {
		return info, nil
	}
	return nil, os.ErrNotExist
}

func (fs *MockFileSystem) Readlink(name string) (string, error) {
	if link, ok := fs.Links[name]; ok {
		return link, nil
	}
	// Check if it exists but is not a link
	if _, ok := fs.Infos[name]; ok {
		return "", errors.New("not a link")
	}
	return "", os.ErrNotExist
}

func (fs *MockFileSystem) ReadFile(name string) ([]byte, error) {
	if data, ok := fs.Files[name]; ok {
		return data, nil
	}
	return nil, os.ErrNotExist
}

func (fs *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	// Simple mock: just assume it works
	return nil
}

func (fs *MockFileSystem) Remove(name string) error {
	delete(fs.Files, name)
	delete(fs.Infos, name)
	delete(fs.Links, name)
	return nil
}

func (fs *MockFileSystem) Rename(oldpath, newpath string) error {
	if _, ok := fs.Infos[oldpath]; !ok {
		return os.ErrNotExist
	}
	if val, ok := fs.Files[oldpath]; ok {
		fs.Files[newpath] = val
		delete(fs.Files, oldpath)
	}
	if val, ok := fs.Infos[oldpath]; ok {
		fs.Infos[newpath] = val
		delete(fs.Infos, oldpath)
	}
	if val, ok := fs.Links[oldpath]; ok {
		fs.Links[newpath] = val
		delete(fs.Links, oldpath)
	}
	return nil
}

func (fs *MockFileSystem) Symlink(oldname, newname string) error {
	fs.Links[newname] = oldname
	fs.Infos[newname] = MockFileInfo{name: newname, mode: os.ModeSymlink}
	return nil
}

func (fs *MockFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	fs.Files[name] = data
	fs.Infos[name] = MockFileInfo{name: name, size: int64(len(data)), mode: perm}
	return nil
}

func (fs *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	// For mock, Stat follows link if it exists? Or just returns Lstat?
	// os.Stat follows links.
	if link, ok := fs.Links[name]; ok {
		return fs.Stat(link)
	}
	return fs.Lstat(name)
}

func TestPlan_CreateFile(t *testing.T) {
	fs := NewMockFS()
	engine := NewEngine(fs)
	ctx := context.Background()

	// Setup: Source file exists
	srcPath := "/source/file"
	if err := fs.WriteFile(srcPath, []byte("content"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	desired := []DesiredState{
		{Target: "/target/file", Source: srcPath, Mode: 0644},
	}
	state := &State{Files: make(map[string]ManagedInfo)}

	actions, err := engine.Plan(ctx, desired, state)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}

	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Type != ActionCreate {
		t.Errorf("expected ActionCreate, got %s", actions[0].Type)
	}
}

func TestPlan_NoopFile(t *testing.T) {
	fs := NewMockFS()
	engine := NewEngine(fs)
	ctx := context.Background()

	srcPath := "/source/file"
	dstPath := "/target/file"
	content := []byte("content")
	if err := fs.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if err := fs.WriteFile(dstPath, content, 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	desired := []DesiredState{
		{Target: dstPath, Source: srcPath, Mode: 0644},
	}
	state := &State{Files: map[string]ManagedInfo{
		dstPath: {Source: srcPath},
	}}

	actions, err := engine.Plan(ctx, desired, state)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}

	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Type != ActionNoop {
		t.Errorf("expected ActionNoop, got %s", actions[0].Type)
	}
}

func TestPlan_UpdateFile_ContentMismatch(t *testing.T) {
	fs := NewMockFS()
	engine := NewEngine(fs)
	ctx := context.Background()

	srcPath := "/source/file"
	dstPath := "/target/file"
	if err := fs.WriteFile(srcPath, []byte("new content"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if err := fs.WriteFile(dstPath, []byte("old content"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	desired := []DesiredState{
		{Target: dstPath, Source: srcPath, Mode: 0644},
	}
	state := &State{Files: map[string]ManagedInfo{
		dstPath: {Source: srcPath},
	}}

	actions, err := engine.Plan(ctx, desired, state)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}

	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Type != ActionUpdate {
		t.Errorf("expected ActionUpdate, got %s", actions[0].Type)
	}
}

func TestPlan_DeleteOrphan(t *testing.T) {
	fs := NewMockFS()
	engine := NewEngine(fs)
	ctx := context.Background()

	dstPath := "/target/orphaned"
	if err := fs.WriteFile(dstPath, []byte("content"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	desired := []DesiredState{}
	state := &State{Files: map[string]ManagedInfo{
		dstPath: {Source: "/source/old"},
	}}

	actions, err := engine.Plan(ctx, desired, state)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}

	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Type != ActionDelete {
		t.Errorf("expected ActionDelete, got %s", actions[0].Type)
	}
}
