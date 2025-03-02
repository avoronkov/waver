package static

import (
	"embed"
	"io/fs"
	"log/slog"
	"path/filepath"
	"sort"
)

//go:embed samples
var Files embed.FS

func ListFiles(dir string) (files []string) {
	subdir := "samples"
	if dir != "" {
		subdir = filepath.Join(subdir, dir)
	}
	subdirlen := len(subdir) + 1
	err := fs.WalkDir(Files, subdir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, path[subdirlen:])
		return nil
	})
	if err != nil {
		slog.Error("WalkDir failed", "error", err)
	}
	sort.Strings(files)
	return files
}
