package scanner

// supportedExtensions is the list of source files we currently inspect for env usage.
var supportedExtensions = map[string]bool{
	".go": true,
	".js": true,
	".py": true,
	".ts": true,
}
