// Basic http functionality that really aught to be in the http library, but isn't. 
package httputil

import (
	"strings"
	"bytes"
	"http"
	"fmt"
	"log"
	"os"
	"io"
)

type HttpResponseWriter struct {
	conn         io.Writer
	buf          *bytes.Buffer
	headers      http.Header
	statusCode   int
	wroteHeaders bool
}

func NewHttpResponseWriter(conn io.Writer) *HttpResponseWriter {
	return &HttpResponseWriter{conn, nil, make(map[string][]string), 0, false}
}

func (h *HttpResponseWriter) Header() http.Header {
	return h.headers
}

// BUG: Assumes http 1.1 and that status is always "OK" (although not necessarily 200).
func (h *HttpResponseWriter) WriteHeader(code int) {
	fmt.Fprintln(h.conn, "HTTP/1.1", code, http.StatusText(code))
	log.Println(h.headers)
	log.Println(h.headers.Get("Content-Length"))
	if h.headers.Get("Content-Length") == "" {
		//h.statusCode = code
		if h.buf == nil {
			h.buf = bytes.NewBuffer([]byte(""))
		}
		log.Println("No content-length")
		return
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

func (h *HttpResponseWriter) Write(buf []byte) (int, os.Error) {
	if !h.wroteHeaders {
		h.WriteHeader(200)
	}
	if h.buf != nil {
		log.Println("BUffer length:", len(buf))
		return h.buf.Write(buf)
	}
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
	if finfo.IsDirectory() {
		http.Error(w, "Access Denied To Folder Listing", 403)
		return
	}
	http.ServeFile(w, r, name)
	return
}
