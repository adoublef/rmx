package template

import (
	"html/template"
	"io"
	"io/fs"
	"maps"
)

type FS struct {
	fsys  fs.FS
	funcs template.FuncMap
}

// ParseFile reads from the file system fs
func (fsys *FS) ParseFiles(patterns ...string) (Template, error) {
	return template.New(patterns[0]).Funcs(fsys.funcs).ParseFS(fsys.fsys, patterns...)
}

// Funcs adds the elements of the argument map to the template's function map.
func (fsys *FS) Funcs(funcs ...map[string]any) *FS {
	for _, f := range funcs {
		maps.Copy(fsys.funcs, f)
	}
	return fsys
}

func NewFS(fsys fs.FS) *FS {
	return &FS{fsys, template.FuncMap{}}
}

type Template interface {
	Execute(wr io.Writer, data any) error
}
