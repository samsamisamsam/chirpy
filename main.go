package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.mwareHits(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (c *apiConfig) mwareHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := fmt.Sprintf(
		`<html>
<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`,
		int(c.fileserverHits.Load()))
	w.Write([]byte(html))
}

func (c *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	c.fileserverHits.Store(0)
}

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpJSON struct {
		Body []byte `json:"body"`
	}
	type chirpError struct {
		Error []byte `json:"error"`
	}
	type chirpValid struct {
		Valid []byte `json:"valid"`
	}

	chirp := chirpJSON{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		respBody := chirpError{
			Error: []byte("Error: Something went wrong"),
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			respBody = chirpError{
				Error: []byte("Error: Something went wrong"),
			}
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
	} else if len(chirp.Body) > 140 {
		respBody := chirpError{
			Error: []byte("Error: chirp is too long"),
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			respBody = chirpError{
				Error: []byte("Error: Something went wrong"),
			}
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
	} else {
		respBody := chirpValid{
			Valid: []byte("valid"),
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			respBody := chirpError{
				Error: []byte("Something went wrong"),
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write(respBody.Error)
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(dat)
	}
}
