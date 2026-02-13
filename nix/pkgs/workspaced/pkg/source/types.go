package source

import (
	"bytes"
	"context"
	"io"
	"os"
)

// Source representa uma origem de arquivos/templates
type Source interface {
	// Name retorna identificador único da source
	Name() string

	// Priority define precedência em conflitos (maior = mais prioritário)
	Priority() int

	// Scan descobre e retorna todos os arquivos desta source
	Scan(ctx context.Context) ([]File, error)
}

// FileType indica o tipo de processamento necessário
type FileType int

const (
	TypeSymlink   FileType = iota // Criar symlink direto
	TypeStatic                    // Copiar arquivo estático (sem template)
	TypeTemplate                  // Renderizar template simples
	TypeMultiFile                 // Template que gera múltiplos arquivos
	TypeDotD                      // Diretório .d.tmpl (concatenação)
)

func (t FileType) String() string {
	switch t {
	case TypeSymlink:
		return "symlink"
	case TypeStatic:
		return "static"
	case TypeTemplate:
		return "template"
	case TypeMultiFile:
		return "multifile"
	case TypeDotD:
		return "dotd"
	default:
		return "unknown"
	}
}

// File representa um arquivo descoberto ou gerado no pipeline
type File interface {
	RelPath() string
	TargetBase() string
	Mode() os.FileMode
	Reader() (io.ReadCloser, error)
	SourceInfo() string
	Type() FileType
}

// BasicFile implementa os campos comuns de File
type BasicFile struct {
	RelPathStr    string
	TargetBaseDir string
	FileMode      os.FileMode
	Info          string
	FileType      FileType
}

func (f *BasicFile) RelPath() string    { return f.RelPathStr }
func (f *BasicFile) TargetBase() string { return f.TargetBaseDir }
func (f *BasicFile) Mode() os.FileMode  { return f.FileMode }
func (f *BasicFile) SourceInfo() string { return f.Info }
func (f *BasicFile) Type() FileType     { return f.FileType }

// StaticFile representa um arquivo real no disco
type StaticFile struct {
	BasicFile
	AbsPath string
}

func (f *StaticFile) Reader() (io.ReadCloser, error) {
	return os.Open(f.AbsPath)
}

// BufferFile representa um arquivo com conteúdo em memória
type BufferFile struct {
	BasicFile
	Content []byte
}

func (f *BufferFile) Reader() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(f.Content)), nil
}

// Conflict representa múltiplos files querendo o mesmo target
type Conflict struct {
	RelPath string // Caminho relativo em conflito
	Files   []File // Ordenados por priority (maior primeiro)
}

// ConflictResolution define estratégia de resolução de conflitos
type ConflictResolution int

const (
	// ResolveByPriority usa arquivo da source com maior priority
	ResolveByPriority ConflictResolution = iota

	// ResolveByError retorna erro parando o processo
	ResolveByError

	// ResolveBySkip ignora todos os arquivos em conflito
	ResolveBySkip
)

func (r ConflictResolution) String() string {
	switch r {
	case ResolveByPriority:
		return "priority"
	case ResolveByError:
		return "error"
	case ResolveBySkip:
		return "skip"
	default:
		return "unknown"
	}
}

// DesiredState representa estado desejado de um arquivo
type DesiredState struct {
	File File
}

// Provider gera estados desejados (interface legacy, mantida para compatibilidade)
type Provider interface {
	Name() string
	GetDesiredState(ctx context.Context) ([]DesiredState, error)
}
