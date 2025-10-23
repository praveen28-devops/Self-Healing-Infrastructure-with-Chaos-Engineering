package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	db *sql.DB
	// --- Metrics (unchanged) ---
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_app_http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path", "method", "status"},
	)
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "go_app_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
	databaseHealth = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_app_database_health",
			Help: "Database health status (1 = healthy, 0 = unhealthy).",
		},
	)

	// --- NEW: Prometheus metric for votes ---
	voteCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_app_votes_total",
			Help: "Total number of votes for each company.",
		},
		[]string{"company"},
	)
)

// Company holds the vote data
type Company struct {
	Name      string
	LogoClass string // For CSS
	Votes     int
}

// PageData holds all the data we need to render the HTML template
type PageData struct {
	Companies []Company
	Healthy   bool
}

// createVotesTable creates the 'votes' table if it doesn't exist
func createVotesTable() {
	// List of companies we want to feature
	companies := []string{"google", "microsoft", "apple", "amazon", "nvidia", "meta"}

	// Create the table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS votes (
		id SERIAL PRIMARY KEY,
		company_name VARCHAR(50) UNIQUE NOT NULL,
		vote_count INT DEFAULT 0
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create votes table: %v", err)
	}

	// Insert company names if they don't already exist.
	// This ensures our table is pre-populated with the companies we want.
	insertSQL := `INSERT INTO votes (company_name) VALUES ($1) ON CONFLICT (company_name) DO NOTHING;`
	for _, company := range companies {
		_, err := db.Exec(insertSQL, company)
		if err != nil {
			log.Printf("Failed to insert company %s: %v", company, err)
		}
	}
	log.Println("Votes table checked and populated.")
}

// initDB initializes the database connection and creates the table
func initDB() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set.")
	}

	var err error
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Printf("Failed to ping database on startup. Error: %v", err)
		databaseHealth.Set(0)
	} else {
		log.Println("Successfully connected to the database.")
		databaseHealth.Set(1)
		// --- NEW: Create the table on startup ---
		createVotesTable()
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
}

// checkDBHealth (unchanged)
func checkDBHealth() bool {
	if db == nil {
		return false
	}
	err := db.Ping()
	if err != nil {
		log.Printf("Database health check failed: %v", err)
		return false
	}
	return true
}

// --- UPDATED: homeHandler ---
// Serves the HTML voting page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Define company logo mappings
	logoMap := map[string]string{
		"google":    "logo-google",
		"microsoft": "logo-microsoft",
		"apple":     "logo-apple",
		"amazon":    "logo-amazon",
		"nvidia":    "logo-nvidia",
		"meta":      "logo-meta",
	}

	// 1. Fetch vote counts from the database
	rows, err := db.Query("SELECT company_name, vote_count FROM votes ORDER BY vote_count DESC")
	if err != nil {
		log.Printf("Error querying votes: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		httpRequestsTotal.WithLabelValues(r.URL.Path, r.Method, "500").Inc()
		return
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var c Company
		if err := rows.Scan(&c.Name, &c.Votes); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		c.LogoClass = logoMap[c.Name] // Add the logo class
		companies = append(companies, c)
	}

	// 2. Prepare data for the template
	pageData := PageData{
		Companies: companies,
		Healthy:   true, // If we got this far, the DB is healthy
	}

	// 3. Parse and execute the HTML template
	// We parse the file directly from the 'templates' folder
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		httpRequestsTotal.WithLabelValues(r.URL.Path, r.Method, "500").Inc()
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, pageData) // Pass the data to the template

	// 4. Record metrics
	duration := time.Since(start).Seconds()
	httpRequestsTotal.WithLabelValues(r.URL.Path, r.Method, "200").Inc()
	httpRequestDuration.WithLabelValues(r.URL.Path, r.Method).Observe(duration)
}

// --- NEW: voteHandler ---
// Handles the POST /vote request
func voteHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// 1. Decode the incoming JSON
	var requestBody struct {
		Company string `json:"company"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		httpRequestsTotal.WithLabelValues(r.URL.Path, r.Method, "400").Inc()
		return
	}

	// 2. Update the vote count in the database
	updateSQL := `
	UPDATE votes 
	SET vote_count = vote_count + 1 
	WHERE company_name = $1 
	RETURNING vote_count;`

	var updatedCount int
	err := db.QueryRow(updateSQL, requestBody.Company).Scan(&updatedCount)
	if err != nil {
		log.Printf("Error updating vote count: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		httpRequestsTotal.WithLabelValues(r.URL.Path, r.Method, "500").Inc()
		return
	}

	// 3. Increment Prometheus counter
	voteCounter.WithLabelValues(requestBody.Company).Inc()

	// 4. Send back the new count
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"company":   requestBody.Company,
		"new_count": updatedCount,
	})

	duration := time.Since(start).Seconds()
	httpRequestsTotal.WithLabelValues(r.URL.Path, r.Method, "200").Inc()
	httpRequestDuration.WithLabelValues(r.URL.Path, r.Method).Observe(duration)
}

// healthHandler (unchanged from your original)
func healthHandler(w http.ResponseWriter, r *http.Request) {
	isHealthy := checkDBHealth()
	w.Header().Set("Content-Type", "application/json")

	if isHealthy {
		databaseHealth.Set(1)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	} else {
		databaseHealth.Set(0)
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy"})
	}
}

// --- UPDATED: main ---
func main() {
	// Initialize database connection
	initDB()

	// --- NEW: Handle static files (CSS, JS) ---
	// Create a file server for the 'static' directory
	fs := http.FileServer(http.Dir("./static"))
	// Serve files under the /static/ URL path
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register other handlers
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/vote", voteHandler) // NEW
	http.HandleFunc("/health", healthHandler)
	http.Handle("/metrics", promhttp.Handler())

	// Start the web server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
