package templates

import "embed"

var (
	//go:embed layouts *.html
	FS embed.FS
)
