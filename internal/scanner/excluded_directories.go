package scanner

var excludedDirectories = map[string]bool{
	".git":          true,
	".hg":           true,
	".idea":         true,
	".mypy_cache":   true,
	".pytest_cache": true,
	".ruff_cache":   true,
	".svn":          true,
	".tox":          true,
	".venv":         true,
	"venv":          true,
	".vscode":       true,
	"__pycache__":   true,
	"build":         true,
	"dist":          true,
	"node_modules":  true,
	"vendor":        true,
}
