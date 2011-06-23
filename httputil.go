// Basic http functionality that really aught to be in the http library, but isn't. 
package httputil

import (
	"strings"
	"http"
	"fmt"
	"os"
	"io"
)

type HttpResponseWriter struct {
	conn         io.Writer
	headers      http.Header
	wroteHeaders bool
}

func NewHttpResponseWriter(conn io.Writer) *HttpResponseWriter {
	return &HttpResponseWriter{conn, make(map[string][]string), false}
}

func (h *HttpResponseWriter) Header() http.Header {
	return h.headers
}

// BUG: Assumes http 1.1 and that status is always "OK" (although not necessarily 200).
func (h *HttpResponseWriter) WriteHeader(code int) {
	fmt.Fprintln(h.conn, "HTTP/1.1", code, "OK")
	for name, value := range h.headers {
		valuestr := ""
		for _, singlevalue := range value {
			valuestr += singlevalue + "; "
		}
		fmt.Fprintf(h.conn, "%s: %s\n", name, valuestr)
	}
	fmt.Fprintln(h.conn)
	h.wroteHeaders = true
}

func (h *HttpResponseWriter) Write(buf []byte) (int, os.Error) {
	if !h.wroteHeaders {
		h.WriteHeader(200)
	}
	return h.conn.Write(buf)
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
