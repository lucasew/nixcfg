package template

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"text/template"
)

// Engine é o motor de renderização de templates do workspaced
type Engine struct {
	funcMap template.FuncMap
}

// Option é uma função de configuração para Engine
type Option func(*Engine)

// NewEngine cria uma nova engine de templates
func NewEngine(ctx context.Context, opts ...Option) *Engine {
	e := &Engine{
		funcMap: makeFuncMap(ctx),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// WithCustomFunc adiciona uma função customizada ao FuncMap
func WithCustomFunc(name string, fn interface{}) Option {
	return func(e *Engine) {
		e.funcMap[name] = fn
	}
}

// WithFuncMap substitui o FuncMap inteiro
func WithFuncMap(funcMap template.FuncMap) Option {
	return func(e *Engine) {
		e.funcMap = funcMap
	}
}

// Render renderiza um template string com os dados fornecidos
func (e *Engine) Render(ctx context.Context, tmpl string, data interface{}) ([]byte, error) {
	t, err := template.New("template").Funcs(e.funcMap).Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		if errors.Is(err, ErrFileSkipped) {
			return nil, ErrFileSkipped
		}
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// RenderFile renderiza um arquivo de template
func (e *Engine) RenderFile(ctx context.Context, path string, data interface{}) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	return e.Render(ctx, string(content), data)
}

// RenderMultiFile renderiza e retorna múltiplos arquivos
func (e *Engine) RenderMultiFile(ctx context.Context, tmpl string, data interface{}) ([]MultiFile, error) {
	rendered, err := e.Render(ctx, tmpl, data)
	if err != nil {
		return nil, err
	}

	files, isMulti := ParseMultiFile(rendered)
	if !isMulti {
		return nil, fmt.Errorf("template is not a multi-file template")
	}

	return files, nil
}
