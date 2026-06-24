package main

import (
	"flag"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	mux       = http.NewServeMux()
	DS        *Datastore
	DbHost    string
	DbUser    string
	DbPass    string
	DbPort    int
	DbName    string
	DbSSLMode string
)

type Content struct {
	Title    string
	Hostname string
	Entries  []*LogEntry
}

func init() {
	flag.StringVar(&DbHost, "db-host", "", "database host")
	flag.IntVar(&DbPort, "db-port", 0, "database port")
	flag.StringVar(&DbUser, "db-user", "", "database username")
	flag.StringVar(&DbPass, "db-pass", "", "database password")
	flag.StringVar(&DbName, "db-name", "", "database name")
	flag.StringVar(&DbSSLMode, "db-ssl", "", "database ssl mode")
}

func loadTemplate(filename string) (*template.Template, error) {
	return template.ParseFiles(filename)
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := loadTemplate("templates/index.html.tmpl")
	if err != nil {
		slog.Error("error loading template", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	t.Execute(w, &Content{
		Title:    os.Getenv("TITLE"),
		Hostname: hostname,
	})
}

func dbIndex(w http.ResponseWriter, r *http.Request) {
	DS.Add(&LogEntry{
		Date: time.Now(),
		Addr: r.RemoteAddr,
		Path: r.URL.Path,
	})

	t, err := loadTemplate("templates/db.html.tmpl")
	if err != nil {
		slog.Error("error loading template", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	entries, err := DS.GetAll()
	if err != nil {
		slog.Error("error retrieving data", "error", err)
	}

	t.Execute(w, &Content{
		Title:    os.Getenv("TITLE"),
		Hostname: hostname,
		Entries:  entries,
	})
}

func main() {
	flag.Parse()

	if DbHost == "" {
		DbHost = os.Getenv("DB_HOST")
	}
	if DbPort == 0 {
		if port := os.Getenv("DB_PORT"); port != "" {
			p, err := strconv.Atoi(port)
			if err != nil {
				log.Fatalf("unable to parse database port: %s", err)
			}
			DbPort = p
		}
	}
	if DbUser == "" {
		DbUser = os.Getenv("DB_USER")
	}
	if DbPass == "" {
		DbPass = os.Getenv("DB_PASS")
	}
	if DbName == "" {
		DbName = os.Getenv("DB_NAME")
	}
	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" && DbSSLMode == "" {
		DbSSLMode = sslMode
	}

	ds, err := NewDatastore(DbHost, DbPort, DbUser, DbPass, DbName, DbSSLMode)
	if err != nil {
		log.Fatalf("unable to connect to datastore: %s", err)
	}
	DS = ds
	defer DS.Close()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(os.DirFS("./static"))))
	mux.HandleFunc("GET /", index)
	mux.HandleFunc("GET /db", dbIndex)

	slog.Info("listening", "addr", ":8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("error serving: %s", err)
	}
}
