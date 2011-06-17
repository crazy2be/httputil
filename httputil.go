// Basic http functionality that really aught to be in the http library, but isn't. 
package httputil

import (
	"strings"
	"http"
	"os"
)

// Searches for the cookie given by key in the request r, returning the value of the first found match. Can be inefficient if there are many cookies, as it does no sorting. Returns nil if no cookie was found. Case-insensitive.
func FindCookie(r *http.Request, key string) *http.Cookie {
	cookiearray := r.Cookie
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
		http.Error(w, "Unable to open file " + name, 404)
		return
	}
	if finfo.IsDirectory() {
		http.Error(w, "Access Denied To Folder Listing", 403)
		return
	}
	http.ServeFile(w, r, name)
	return
}
