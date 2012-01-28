// Basic http functionality that really aught to be in the http library, but isn't. 
package httputil

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type HttpResponseWriter struct {
	conn         io.Writer
	buf          *bytes.Buffer
	headers      http.Header
	noContent    bool
	statusCode   int
	wroteHeaders bool
}

// Returns a new HttpResponseWriter, useful for calls to httputil.ServeFileOnly() where you only have a raw connection and no http.ResponseWriter. Note that this only works with files or other content where the length is known beforehand: if the content-length is not set explicitly, it defaults to zero, and no content is sent!
func NewHttpResponseWriter(conn io.Writer) *HttpResponseWriter {
	return &HttpResponseWriter{conn, nil, make(map[string][]string), false, 0, false}
}

func (h *HttpResponseWriter) Header() http.Header {
	return h.headers
}

// BUG: If content-length is not set explicitly, it is set to 0 and any calls to h.Write() do nothing. Supporting chunked encoding would be a pain...
func (h *HttpResponseWriter) WriteHeader(code int) {
	fmt.Fprintln(h.conn, "HTTP/1.1", code, http.StatusText(code))
	if h.headers.Get("Content-Length") == "" {
		h.headers.Set("Content-Length", "0")
		h.noContent = true
	}
	log.Println("Writing headers")
	for name, value := range h.headers {
		valuestr := ""
		for _, singlevalue := range value {
			valuestr += singlevalue + "; "
		}
		valuestr = valuestr[:len(valuestr)-2]
		fmt.Fprintf(h.conn, "%s: %s\n", name, valuestr)
	}
	fmt.Fprintln(h.conn)
	h.wroteHeaders = true
}

func (h *HttpResponseWriter) Write(buf []byte) (int, error) {
	if !h.wroteHeaders {
		h.WriteHeader(200)
	}
	if h.noContent {
		return 0, nil
	}
	// 	if h.buf != nil {
	// 		log.Println("BUffer length:", len(buf))
	// 		return h.buf.Write(buf)
	// 	}
	n, err := h.conn.Write(buf)
	log.Println(n, err)
	return n, err
}

func (h *HttpResponseWriter) Flush() {
	h.headers.Set("Content-Length", fmt.Sprintf("%d", h.buf.Len()))
	h.WriteHeader(200)
	h.conn.Write(h.buf.Bytes())
}

// Searches for the cookie given by key in the request r, returning the value of the first found match. Can be inefficient if there are many cookies, as it does no sorting. Returns nil if no cookie was found. Case-insensitive.
// DEPRECATED, use r.Cookie(name) instead.
func FindCookie(r *http.Request, key string) *http.Cookie {
	cookiearray := r.Cookies()
	for _, cookie := range cookiearray {
		if strings.ToLower(cookie.Name) == strings.ToLower(key) {
			return cookie
		}
	}
	return nil
}

// Serves a file, but refuses to list directory contents. Attempting to access a directory will result in a 403 (access denied) error.
func ServeFileOnly(w http.ResponseWriter, r *http.Request, name string) {
	finfo, e := os.Stat(name)
	if e != nil {
		http.Error(w, "Unable to open file "+name, 404)
		return
	}
	if finfo.IsDir() {
		http.Error(w, "Access Denied To Folder Listing", 403)
		return
	}
	http.ServeFile(w, r, name)
	return
}
