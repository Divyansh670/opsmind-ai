package api

import "net/http"

// allowedOrigins lists every frontend origin permitted to call this API
var allowedOrigins = map[string]bool{
	"http://localhost:5173": true, // Vite dev server
	"http://localhost:3000": true, // Dockerized frontend (Nginx)
}

// CORSMiddleware allows known frontend origins to call our API
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
