package views

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type TemplateManager struct {
	fs  embed.FS
	dir string
	ext string

	templates map[string]*template.Template
}

func NewTemplateManager(fs embed.FS, dir string, ext string) (tmpl *TemplateManager, err error) {
	tmpl = &TemplateManager{
		fs:  fs,
		dir: dir,
		ext: ext,
	}

	tmpl.templates = make(map[string]*template.Template)

	if err := tmpl.Load(); err != nil {
		tmpl = nil
	}

	return
}

// Dir returns absolute path to directory with views
func (t *TemplateManager) Dir() string {
	return t.dir
}

// Ext returns extension of views
func (t *TemplateManager) Ext() string {
	return t.ext
}

func (t *TemplateManager) Load() (err error) {
	fmt.Println("Loading templates...")

	var walkDir = func(path string, dir fs.DirEntry, err error) (_ error) {
		if err != nil {
			return err
		}

		var info os.FileInfo
		if info, err = dir.Info(); err != nil {
			return err
		}

		// skip over directories
		if info.IsDir() {
			return
		}

		// get relative path as name
		var name string
		if name, err = filepath.Rel(t.dir, path); err != nil {
			return err
		}

		// trim extension from name
		name = strings.TrimSuffix(name, t.ext)

		var (
			nt = template.New(name)
			b  []byte
		)

		if b, err = t.fs.ReadFile(path); err != nil {
			return err
		}

		tmpl, err := nt.Parse(string(b))

		if err != nil {
			return fmt.Errorf("error parsing template: %w", err)
		}

		t.templates[tmpl.Name()] = tmpl

		return err
	}

	if err = fs.WalkDir(t.fs, t.dir, walkDir); err != nil {
		return err
	}

	return
}

func (t *TemplateManager) Render(w io.Writer, name string, data interface{}) (err error) {

	tmpl, exists := t.templates[name]

	if !exists {
		return fmt.Errorf("template not found: %v", name)
	}

	buf := bytes.Buffer{}

	if err = tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	_, err = w.Write(buf.Bytes())

	return
}
