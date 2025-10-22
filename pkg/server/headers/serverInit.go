package headers

import (
	"log"
	"net/http"
)

func ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// We want to handle requests to /test separately
	if req.URL.Path == "/test" {
		_, err := w.Write([]byte("hello"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	// Handle other requests as needed, or return a 404
	http.NotFound(w, req)
}
