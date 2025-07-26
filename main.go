package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"main.go/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	dbQueries := database.New(db)

	config := apiConfig{
		db: dbQueries,
	}

	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", healdeReadiness)
	mux.Handle("GET /admin/metrics", http.HandlerFunc(config.handleAdminMetrics))
	mux.Handle("POST /admin/reset", http.HandlerFunc(config.handleReset))
	mux.Handle("POST /api/validate_chirp", http.HandlerFunc(JsonHandler))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func healdeReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))

}

// func (c *apiConfig) handleHit(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Hits: " + strconv.Itoa(int(c.fileserverHits.Load()))))
// }

func (c *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func (c *apiConfig) handleAdminMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	html := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, c.fileserverHits.Load())

	w.Write([]byte(html))
}

func JsonHandler(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, `{"error": "Something went wrong"}`, http.StatusBadRequest)
		return
	}
	if len(params.Body) > 140 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(` {"error": "Chirp is too long"}`))
		return
	}
	params.Body = replacebadwords(params.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := fmt.Sprintf(`{"cleaned_body": "%s"}`, params.Body)
	w.Write([]byte(res))
}

func replacebadwords(s string) string {
	bad_words := []string{"kerfuffle", "sharbert", "fornax"}

	for _, word := range strings.Split(s, " ") {
		for _, bad := range bad_words {
			if strings.EqualFold(word, bad) {
				s = strings.ReplaceAll(s, word, "****")
			}
		}
	}

	return s
}
