package conf

import (
	"os"
	"path/filepath"
)

func ExecDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}

	dir := filepath.Dir(exe)
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		return filepath.Dir(exe)
	}

	return dir
}

func DefaultLogDir() string {
	return filepath.Join(ExecDir(), "logs")
}

func DefaultConfigCandidates() []string {
	dir := ExecDir()

	return []string{
		filepath.Join(dir, "config.enc"),
		filepath.Join(dir, "config.yaml"),
	}
}
