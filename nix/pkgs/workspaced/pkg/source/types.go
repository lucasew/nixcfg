package source

import (
	"context"
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
	TypeStatic                     // Copiar arquivo estático (sem template)
	TypeTemplate                   // Renderizar template simples
	TypeMultiFile                  // Template que gera múltiplos arquivos
	TypeDotD                       // Diretório .d.tmpl (concatenação)
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

// File representa um arquivo descoberto por uma source
type File struct {
	SourceName string      // Nome da source de origem
	SourcePath string      // Caminho absoluto no source
	TargetPath string      // Caminho absoluto de destino no sistema
	Type       FileType    // Tipo de processamento
	Content    []byte      // Conteúdo (já renderizado se template)
	Mode       os.FileMode // Permissões (0 = symlink)
	Priority   int         // Priority da source (para resolver conflitos)
}

// Conflict representa múltiplos files querendo o mesmo target
type Conflict struct {
	TargetPath string
	Files      []File // Ordenados por priority (maior primeiro)
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
