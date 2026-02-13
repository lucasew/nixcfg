package source

import (
	"bytes"
	"context"
	"io"
	"workspaced/pkg/template"
)

// TemplateFile representa um template renderizado sob demanda
type TemplateFile struct {
	BasicFile
	SourceFile File
	Engine     *template.Engine
	Data       interface{}
	Context    context.Context
}

func (f *TemplateFile) Reader() (io.ReadCloser, error) {
	// Read source content
	srcReader, err := f.SourceFile.Reader()
	if err != nil {
		return nil, err
	}
	defer srcReader.Close()

	srcContent, err := io.ReadAll(srcReader)
	if err != nil {
		return nil, err
	}

	// Render
	rendered, err := f.Engine.Render(f.Context, string(srcContent), f.Data)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(rendered)), nil
}

// ConcatenatedFile representa múltiplos arquivos unidos (DotD)
type ConcatenatedFile struct {
	BasicFile
	Components []File
}

func (f *ConcatenatedFile) Reader() (io.ReadCloser, error) {
	readers := []io.Reader{}
	for i, c := range f.Components {
		r, err := c.Reader()
		if err != nil {
			return nil, err
		}
		// Note: we might leak readers if we don't close them.
		// io.MultiReader doesn't close underlying readers.
		// However, most readers here are BufferFiles or StaticFiles.
		// For safety, let's wrap them.
		readers = append(readers, &autoCloserReader{r})

		// Add newline between components if not the last one
		if i < len(f.Components)-1 {
			readers = append(readers, bytes.NewReader([]byte("\n")))
		}
	}
	return io.NopCloser(io.MultiReader(readers...)), nil
}

type autoCloserReader struct {
	inner io.ReadCloser
}

func (r *autoCloserReader) Read(p []byte) (n int, err error) {
	n, err = r.inner.Read(p)
	if err == io.EOF {
		r.inner.Close()
	}
	return n, err
}
