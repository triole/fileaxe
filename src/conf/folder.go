package conf

import "path/filepath"

// FolderArg turns the folder argument into the absolute path to process. With
// follow set, symlinks are resolved as well, so a symlinked folder is processed
// like its target: filepath.Walk does not descend into a symlinked root on its
// own. Resolution is off by default, which keeps the walk's symlink handling as
// it is inside the tree.
func FolderArg(folder string, follow bool) (string, error) {
	abs, err := filepath.Abs(folder)
	if err != nil {
		return "", err
	}
	if !follow {
		return abs, nil
	}
	return filepath.EvalSymlinks(abs)
}
