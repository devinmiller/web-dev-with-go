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

	layoutDir string

	templates map[string]*template.Template
}

func NewTemplateManager(fs embed.FS, dir string, layoutDir string, ext string) (tmpl *TemplateManager, err error) {
	tmpl = &TemplateManager{
		fs:        fs,
		dir:       dir,
		ext:       ext,
		layoutDir: layoutDir,
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

		fmt.Println(path)

		var info os.FileInfo
		if info, err = dir.Info(); err != nil {
			return err
		}

		// skip over directories
		if info.IsDir() {
			return
		}

		// skip files in the layout directory
		if filepath.Dir(path) == t.layoutDir {
			return
		}

		// get relative path
		var rel string
		if rel, err = filepath.Rel(t.dir, path); err != nil {
			return err
		}

		// trim extension from path
		rel = strings.TrimSuffix(rel, t.ext)

		// create a new template
		var nt = template.New(filepath.Base(path))

		// parse the template and layouts
		tmpl, err := nt.ParseFS(t.fs, path, t.layoutDir+"/*"+t.ext)

		if err != nil {
			return fmt.Errorf("error parsing template: %w", err)
		}

		t.templates[rel] = tmpl

		return err
	}

	if err = fs.WalkDir(t.fs, t.dir, walkDir); err != nil {
		return err
	}

	return
}

func (t *TemplateManager) Template(name string) (tmpl *template.Template, err error) {
	tmpl, exists := t.templates[name]

	if !exists {
		return nil, fmt.Errorf("template not found: %v", name)
	}

	return tmpl.Clone()
}

func (t *TemplateManager) Render(w io.Writer, name string, data interface{}) (err error) {

	tmpl, err := t.Template(name)

	if err != nil {
		return err
	}

	buf := bytes.Buffer{}

	if err = tmpl.ExecuteTemplate(&buf, tmpl.Name(), data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	_, err = w.Write(buf.Bytes())

	return
}
