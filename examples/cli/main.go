package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/azghr/forge/envconfig"
	"github.com/azghr/forge/flagsub"
	"github.com/azghr/forge/option"
	"github.com/azghr/forge/regexcache"
	"github.com/azghr/forge/shellquote"
	"github.com/azghr/forge/sliceutil"
	"github.com/azghr/forge/tablewriter"
)

type Config struct {
	TasksFile string `env:"TASKS_FILE,default=tasks.json"`
}

type Task struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Done    bool   `json:"done"`
	DueDate string `json:"due_date,omitempty"`
}

type TaskStore struct {
	mu    sync.Mutex
	Tasks []Task `json:"tasks"`
	Next  int    `json:"next"`
	path  string
}

var reCache = regexcache.New()

// findTask returns the task with the given id, if it exists.
func (s *TaskStore) findTask(id int) option.Option[Task] {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.Tasks {
		if t.ID == id {
			return option.Some(t)
		}
	}
	return option.None[Task]()
}

func main() {
	var cfg Config
	if err := envconfig.Load(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}

	store := loadStore(cfg.TasksFile)

	var addCmd *flagsub.Sub
	addCmd = flagsub.AddSub("add", "Add a task", func() {
		addCmd.Flags.Parse(os.Args[2:])
		due := addCmd.Flags.Lookup("due").Value.String()
		if addCmd.Flags.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "usage: cli add [--due YYYY-MM-DD] <title>")
			fmt.Fprintln(os.Stderr, "  flags must come before the title")
			os.Exit(1)
		}
		title := strings.Join(addCmd.Flags.Args(), " ")
		store.add(title, due)
		fmt.Printf("Added task: %s\n", title)
	})
	addCmd.Flags.String("due", "", "due date (YYYY-MM-DD)")

	flagsub.AddSub("list", "List active tasks", func() {
		store.mu.Lock()
		tasks := append([]Task(nil), store.Tasks...)
		store.mu.Unlock()
		tasks = sliceutil.Filter(tasks, func(t Task) bool { return !t.Done })
		if len(tasks) == 0 {
			fmt.Println("No active tasks.")
			return
		}
		tw := tablewriter.New([]string{"ID", "Title", "Due"})
		for _, t := range tasks {
			tw.Append(fmt.Sprintf("%d", t.ID), t.Title, t.DueDate)
		}
		tw.Write(os.Stdout)
	})

	flagsub.AddSub("done", "Mark a task as complete", func() {
		idStr := strings.Join(os.Args[2:], " ")
		var id int
		n, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil || n != 1 {
			fmt.Fprintf(os.Stderr, "usage: cli done <id>\n")
			os.Exit(1)
		}
		store.complete(id)
	})

	flagsub.AddSub("rm", "Delete a task", func() {
		idStr := strings.Join(os.Args[2:], " ")
		var id int
		n, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil || n != 1 {
			fmt.Fprintf(os.Stderr, "usage: cli rm <id>\n")
			os.Exit(1)
		}
		store.remove(id)
	})

	flagsub.AddSub("search", "Search tasks by regex", func() {
		pattern := strings.Join(os.Args[2:], " ")
		re, err := reCache.Compile(pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid regex: %v\n", err)
			os.Exit(1)
		}
		store.mu.Lock()
		tasks := append([]Task(nil), store.Tasks...)
		store.mu.Unlock()
		matched := sliceutil.Filter(tasks, func(t Task) bool {
			return re.MatchString(t.Title)
		})
		tw := tablewriter.New([]string{"ID", "Title", "Done"})
		for _, t := range matched {
			doneStr := ""
			if t.Done {
				doneStr = "✓"
			}
			tw.Append(
				fmt.Sprintf("%d", t.ID),
				shellquote.Quote(t.Title),
				doneStr,
			)
		}
		if tw.Len() == 0 {
			fmt.Println("No matches.")
			return
		}
		tw.Write(os.Stdout)
	})

	// Demonstrate option.Option: look up a task by id from args, if provided.
	flagsub.AddSub("show", "Show a task by ID", func() {
		idStr := strings.Join(os.Args[2:], " ")
		var id int
		n, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil || n != 1 {
			fmt.Fprintf(os.Stderr, "usage: cli show <id>\n")
			os.Exit(1)
		}
		taskOpt := store.findTask(id)
		if !taskOpt.IsSome() {
			fmt.Fprintf(os.Stderr, "task %d not found\n", id)
			os.Exit(1)
		}
		t := taskOpt.Must()
		fmt.Printf("ID: %d\nTitle: %s\nDone: %v\nDue: %s\n", t.ID, t.Title, t.Done, t.DueDate)
	})

	flagsub.Parse()
}

func loadStore(path string) *TaskStore {
	data, err := os.ReadFile(path)
	if err != nil {
		return &TaskStore{Next: 1, path: path}
	}
	var s TaskStore
	if err := json.Unmarshal(data, &s); err != nil {
		fmt.Fprintf(os.Stderr, "corrupt store: %v\n", err)
		os.Exit(1)
	}
	s.path = path
	return &s
}

func (s *TaskStore) save() {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "save: %v\n", err)
		return
	}
	os.WriteFile(s.path, data, 0644)
}

func (s *TaskStore) add(title, due string) {
	s.mu.Lock()
	s.Tasks = append(s.Tasks, Task{
		ID:      s.Next,
		Title:   title,
		DueDate: due,
	})
	s.Next++
	s.mu.Unlock()
	s.save()
}

func (s *TaskStore) complete(id int) {
	taskOpt := s.findTask(id)
	if !taskOpt.IsSome() {
		fmt.Fprintf(os.Stderr, "task %d not found\n", id)
		return
	}
	s.mu.Lock()
	for i := range s.Tasks {
		if s.Tasks[i].ID == id {
			s.Tasks[i].Done = true
			s.mu.Unlock()
			s.save()
			fmt.Printf("Task %d marked done.\n", id)
			return
		}
	}
	s.mu.Unlock()
}

func (s *TaskStore) remove(id int) {
	taskOpt := s.findTask(id)
	if !taskOpt.IsSome() {
		fmt.Fprintf(os.Stderr, "task %d not found\n", id)
		return
	}
	s.mu.Lock()
	for i, t := range s.Tasks {
		if t.ID == id {
			s.Tasks = append(s.Tasks[:i], s.Tasks[i+1:]...)
			s.mu.Unlock()
			s.save()
			fmt.Printf("Task %d deleted.\n", id)
			return
		}
	}
	s.mu.Unlock()
}
