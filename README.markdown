HTTP Utility Package
====================

Getting Started
---------------

Install:

    goinstall github.com/crazy2be/httputil

Import:

    import "github.com/crazy2be/httputil"

Use:

    httputil.ServeFileOnly(rw, r, path)

What to Use it For
------------------
HTTP functionality that really aught to be in the http package, but isn't for whatever reason. Designed to supplement the existing http package.

More will be added as I find myself needing extra functionality.

Functions
---------

### func ServeFileOnly(w http.ResponseWriter, r *http.Request, name string)
Serves a file, but refuses to list directory contents. Attempting to access a directory will result in a 403 (access denied) error.

