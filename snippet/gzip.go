package snippet

import (
	"compress/gzip"
	"io"
	"os"
)

func gzFileReader(fname string) (io.ReadCloser, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	// Use io.Pipe and a goroutine to create reader
	r, w := io.Pipe()
	go func() {
		defer f.Close()

		// Copy file through gzip to pipe writer.
		gzw := gzip.NewWriter(w)
		_, err := io.Copy(gzw, f)
		if err != nil {
			w.CloseWithError(err)
			return
		}

		// Flush the gzip writer.
		w.CloseWithError(gzw.Close())
	}()
	return r, nil
}
