package gitpath

// Parse returns git location parsed from the given string.
// Supports both local and remote paths.
func Parse(location string) (*GitPath, error) {
	if path, err := FromDir(location); err == nil {
		return path, nil
	}
	return FromURL(location)
}
