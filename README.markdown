httputil Package
================

    import "github.com/crazy2be/httputil"

What It Is
----------
Basic http functionality that really aught to be in the http library, but isn't.

Functions
---------

### func FindCookie(r *http.Request, key string) *http.Cookie
Searches for the cookie given by key in the request r, returning the value of the first found match. Can be inefficient if there are many cookies, as it does no sorting. Returns nil if no cookie was found. Case-insensitive.

### func ServeFileOnly(w http.ResponseWriter, r *http.Request, name string)
Serves a file, but refuses to list directory contents. Attempting to access a directory will result in a 403 (access denied) error.
