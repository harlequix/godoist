# godoist

A Go client library for the [Todoist API v1](https://developer.todoist.com/api/v1/).

## Install

```sh
go get github.com/harlequix/godoist
```

## Usage

### Basic Usage

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/harlequix/godoist"
)

func main() {
	td := godoist.NewTodoist(os.Getenv("TODOIST_TOKEN"))

	// Fetch all tasks and projects (uses REST API by default)
	if err := td.Sync(); err != nil {
		log.Fatal(err)
	}

	// List projects
	for _, p := range td.Projects.All() {
		fmt.Printf("Project: %s (%s)\n", p.Name, p.ID)
	}

	// List tasks
	for _, t := range td.Tasks.All() {
		fmt.Printf("  [%s] %s (priority: %s)\n", t.ID, t.Content, t.Priority)
	}

	// Create a task
	task, err := td.Tasks.Create("Buy groceries")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created: %s (%s)\n", task.Content, task.ID)

	// Update a task
	task.Update("content", "Buy groceries and snacks")

	// Add a label
	task.AddLabel("errands")

	// Complete a task
	task.Close()

	// Working with context (metadata stored in comments)
	// Add context to a task
	task.SetContext(map[string]interface{}{
		"source": "email",
		"urgency": "high",
		"tags": []string{"work", "important"},
	})

	// Get context from a task
	context, err := task.GetContext()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Task context: %v\n", context)

	// Update specific fields in context
	task.UpdateContext(map[string]interface{}{
		"urgency": "medium",
	})

	// Delete a specific field from context
	task.DeleteContextField("tags")

	// Delete all context
	task.DeleteContext()

	// Working with comments
	comments, err := task.GetComments()
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range comments {
		fmt.Printf("Comment: %s\n", c.Content)
	}
}
```

### Using Sync API

For better performance, you can use the Todoist Sync API endpoint instead of REST API calls:

```go
config := &godoist.Config{
	Token:      os.Getenv("TODOIST_TOKEN"),
	UseSyncAPI: true,
}
td := godoist.NewTodoistWithConfig(config)

// This will now use /sync endpoint instead of separate REST calls
if err := td.Sync(); err != nil {
	log.Fatal(err)
}
```

### Configuration Options

You can configure the client using a `Config` struct:

```go
config := &godoist.Config{
	Token:      "your-token",
	ApiURL:     "https://api.todoist.com/api/v1",  // default
	Timeout:    30,                                 // default
	Debug:      false,                              // default
	UseSyncAPI: true,                               // use /sync endpoint (default: false)
}
td := godoist.NewTodoistWithConfig(config)
```

Configuration can also be loaded from files (YAML/TOML) and environment variables:

```go
config, err := godoist.BuildConfig(
	[]string{"config.yaml", "config.toml"}, // config files
	"TODOIST_",                              // env prefix
	&godoist.Config{UseSyncAPI: true},      // overrides
)
if err != nil {
	log.Fatal(err)
}
td := godoist.NewTodoistWithConfig(config)
```

## License

MIT
