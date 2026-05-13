package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/azghr/forge/atomicfile"
	"github.com/azghr/forge/envconfig"
	"github.com/azghr/forge/flagsub"
	"github.com/azghr/forge/pathsafe"
	"github.com/azghr/forge/stopwatch"
	"github.com/azghr/forge/stringutil"
	"github.com/azghr/forge/validator"
)

type Config struct {
	Port    int    `env:"PORT,default=8080"`
	DataDir string `env:"DATA_DIR,default=./data"`
}

type Task struct {
	Title       string `json:"title" validate:"nonzero"`
	Description string `json:"description"`
}

type StoredTask struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Store struct {
	mu    sync.RWMutex
	Tasks map[string]StoredTask `json:"tasks"`
	path  string
}

func main() {
	var (
		port     int
		serveCmd *flagsub.Sub
	)
	serveCmd = flagsub.AddSub("serve", "Start the HTTP server", func() {
		serveCmd.Flags.Parse(os.Args[2:])

		var cfg Config
		if err := envconfig.Load(&cfg); err != nil {
			log.Fatalf("config: %v", err)
		}
		if port > 0 {
			cfg.Port = port
		}

		store, err := newStore(cfg.DataDir)
		if err != nil {
			log.Fatalf("store: %v", err)
		}
		startServer(cfg.Port, store)
	})
	serveCmd.Flags.IntVar(&port, "port", 0, "override PORT env")

	flagsub.AddSub("version", "Print version", func() {
		fmt.Println("forge-server v0.1.0")
	})

	flagsub.Parse()
}

func newStore(dataDir string) (*Store, error) {
	safeDir, err := pathsafe.SafeJoin(".", dataDir)
	if err != nil {
		return nil, fmt.Errorf("unsafe data dir: %w", err)
	}
	if err := os.MkdirAll(safeDir, 0755); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}
	storePath, err := pathsafe.SafeJoin(safeDir, "tasks.json")
	if err != nil {
		return nil, fmt.Errorf("unsafe store path: %w", err)
	}
	s := &Store{
		Tasks: make(map[string]StoredTask),
		path:  storePath,
	}
	data, err := os.ReadFile(storePath)
	if err == nil {
		json.Unmarshal(data, s)
	}
	return s, nil
}

func (s *Store) save() error {
	s.mu.RLock()
	data, err := json.MarshalIndent(s, "", "  ")
	s.mu.RUnlock()
	if err != nil {
		return err
	}
	return atomicfile.WriteFile(s.path, bytes.NewReader(data))
}

func startServer(port int, store *Store) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", withTiming(handleList(store)))
	mux.HandleFunc("POST /tasks", withTiming(handleCreate(store)))
	mux.HandleFunc("GET /tasks/{slug}", withTiming(handleGet(store)))

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func withTiming(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sw stopwatch.Stopwatch
		sw.Start()
		next(w, r)
		sw.Stop()
		log.Printf("%s %s %s", r.Method, r.URL.Path, sw.Elapsed())
	}
}

func handleList(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store.mu.RLock()
		tasks := make([]StoredTask, 0, len(store.Tasks))
		for _, t := range store.Tasks {
			tasks = append(tasks, t)
		}
		store.mu.RUnlock()
		writeJSON(w, http.StatusOK, tasks)
	}
}

func handleCreate(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if err := validator.ValidateStruct(&task); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		slug := stringutil.Slug(task.Title, stringutil.WithSeparator("-"), stringutil.WithMaxLength(40))
		stored := StoredTask{Slug: slug, Title: task.Title, Description: task.Description}

		store.mu.Lock()
		store.Tasks[slug] = stored
		store.mu.Unlock()

		if err := store.save(); err != nil {
			log.Printf("save: %v", err)
		}
		writeJSON(w, http.StatusCreated, stored)
	}
}

func handleGet(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		store.mu.RLock()
		task, ok := store.Tasks[slug]
		store.mu.RUnlock()
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, task)
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
