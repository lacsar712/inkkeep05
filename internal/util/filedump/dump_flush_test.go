package filedump

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestWriteAllFlushesBuffer is a regression test for the bug where WriteAll
// wrote through a bufio.Writer but never flushed it before closing the file.
// As a result on-disk content was empty (small bodies, entirely buffered) or
// truncated to a multiple of the buffer size (large bodies, partially
// auto-flushed). Both cases must now round-trip exactly.
func TestWriteAllFlushesBuffer(t *testing.T) {
	dir := t.TempDir()

	cases := map[string]string{
		// Fits entirely in the default 4096-byte buffer: would be lost without Flush.
		"small": "ink",
		// Larger than the default buffer: without Flush the tail stays buffered
		// and the file is truncated to a multiple of 4096.
		"large": strings.Repeat("INK", 4096),
	}

	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(dir, name+".txt")
			if err := WriteAll(path, body); err != nil {
				t.Fatalf("WriteAll: %v", err)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("ReadFile: %v", err)
			}
			if string(got) != body {
				t.Fatalf("len(got)=%d, want len=%d (content mismatch; buffer not flushed before close)", len(got), len(body))
			}
		})
	}
}
