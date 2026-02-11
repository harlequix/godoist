# godoist

A Go client library for the [Todoist REST API v2](https://developer.todoist.com/rest/v2).

## Install

```sh
go get github.com/harlequix/godoist
```

## Usage

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

	// Fetch all tasks and projects
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
}
```

## License

MIT
