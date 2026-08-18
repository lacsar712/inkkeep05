package filedump

import (
	"bufio"
	"os"
	"path/filepath"
)

func WriteAll(path, body string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	if _, err := w.WriteString(body); err != nil {
		return err
	}
	// Flush the in-memory buffer to the file before the deferred Close runs;
	// otherwise short writes stay buffered and the on-disk content is empty or truncated.
	return w.Flush()
}
