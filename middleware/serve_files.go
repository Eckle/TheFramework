package middleware

import (
	"net/http"
	"os"
)

var dir string = "assets"
var fs http.Handler = http.FileServer(http.Dir(dir))

func ServeFiles(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := dir + r.URL.Path
		info, err := os.Stat(path)
		if os.IsNotExist(err) || info.IsDir() || r.URL.Path == "/" {
			next.ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
