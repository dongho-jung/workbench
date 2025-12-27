// Package task provides task management functionality for TAW.
package task

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// QueueManager handles the quick task queue.
type QueueManager struct {
	queueDir string
}

// NewQueueManager creates a new queue manager.
func NewQueueManager(queueDir string) *QueueManager {
	return &QueueManager{
		queueDir: queueDir,
	}
}

// QueuedTask represents a task in the queue.
type QueuedTask struct {
	Number  int
	Path    string
	Content string
}

// Add adds a new task to the queue.
func (q *QueueManager) Add(content string) error {
	// Ensure queue directory exists
	if err := os.MkdirAll(q.queueDir, 0755); err != nil {
		return fmt.Errorf("failed to create queue directory: %w", err)
	}

	// Get next number
	nextNum, err := q.getNextNumber()
	if err != nil {
		return err
	}

	// Create task file
	filename := fmt.Sprintf("%03d.task", nextNum)
	path := filepath.Join(q.queueDir, filename)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write queue task: %w", err)
	}

	return nil
}

// List returns all queued tasks in order.
func (q *QueueManager) List() ([]QueuedTask, error) {
	entries, err := os.ReadDir(q.queueDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read queue directory: %w", err)
	}

	var tasks []QueuedTask
	pattern := regexp.MustCompile(`^(\d+)\.task$`)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		matches := pattern.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}

		var num int
		fmt.Sscanf(matches[1], "%d", &num)

		path := filepath.Join(q.queueDir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		tasks = append(tasks, QueuedTask{
			Number:  num,
			Path:    path,
			Content: string(content),
		})
	}

	// Sort by number
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Number < tasks[j].Number
	})

	return tasks, nil
}

// Pop removes and returns the first task in the queue.
func (q *QueueManager) Pop() (*QueuedTask, error) {
	tasks, err := q.List()
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, nil
	}

	first := &tasks[0]

	// Remove the task file
	if err := os.Remove(first.Path); err != nil {
		return nil, fmt.Errorf("failed to remove queue task: %w", err)
	}

	return first, nil
}

// Clear removes all tasks from the queue.
func (q *QueueManager) Clear() error {
	tasks, err := q.List()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if err := os.Remove(task.Path); err != nil {
			return fmt.Errorf("failed to remove queue task: %w", err)
		}
	}

	return nil
}

// Count returns the number of tasks in the queue.
func (q *QueueManager) Count() (int, error) {
	tasks, err := q.List()
	if err != nil {
		return 0, err
	}
	return len(tasks), nil
}

// getNextNumber returns the next available task number.
func (q *QueueManager) getNextNumber() (int, error) {
	tasks, err := q.List()
	if err != nil {
		return 1, err
	}

	if len(tasks) == 0 {
		return 1, nil
	}

	return tasks[len(tasks)-1].Number + 1, nil
}

// GenerateTaskName generates a task name from queue task content.
func GenerateTaskNameFromContent(content string, existingNames map[string]bool) string {
	// Get first line or first 30 chars
	lines := strings.Split(content, "\n")
	name := strings.TrimSpace(lines[0])

	if len(name) > 30 {
		name = name[:30]
	}

	// Sanitize: lowercase, replace non-alphanumeric with hyphens
	name = strings.ToLower(name)
	var sb strings.Builder
	lastHyphen := false

	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
			lastHyphen = false
		} else if !lastHyphen {
			sb.WriteRune('-')
			lastHyphen = true
		}
	}

	name = strings.Trim(sb.String(), "-")

	// Ensure minimum length
	if len(name) < 8 {
		name = "queue-task-" + name
	}

	// Handle duplicates
	baseName := name
	counter := 1
	for existingNames[name] {
		name = fmt.Sprintf("%s-%d", baseName, counter)
		counter++
	}

	return name
}
