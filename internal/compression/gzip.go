// Package compression provides structures and methods for handling gzip compression.
package compression

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// Compressor is a simple struct to bind gzip middleware method.
type Compressor struct {
}

// CompressWriter wraps an http.ResponseWriter, providing gzip compression.
type CompressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewCompressWriter creates a new CompressWriter for given http.ResponseWriter.
func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header returns the header map that will be sent by WriteHeader.
func (c *CompressWriter) Header() http.Header {
	return c.w.Header()
}

// Write writes the data to the connection as part of an HTTP reply with gzip compression.
func (c *CompressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader sends an HTTP response header with the provided status code.
func (c *CompressWriter) WriteHeader(statusCode int) {
	if statusCode >= 199 && statusCode <= 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close closes the gzip writer.
func (c *CompressWriter) Close() error {
	return c.zw.Close()
}

type CompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewCompressReader CompressReader wraps an io.ReadCloser, providing gzip decompression.
func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &CompressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read reads the uncompressed form of the compressed data from the reader.
func (c CompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closes the gzip reader.
func (c *CompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// GzipMiddleware is a middleware function for compressing and decompressing HTTP traffic.
func (c *Compressor) GzipMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ow := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := NewCompressWriter(w)
			ow = cw
			defer cw.Close()
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}
		h.ServeHTTP(ow, r)
	}
	return http.HandlerFunc(fn)
}
