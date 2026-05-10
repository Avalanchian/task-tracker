# task-tracker
An implementation of the roadmap.sh task tracker cli project using Go (https://roadmap.sh/projects/task-tracker). For educational purposes.

## Usage
Invoke the app along with one of the following commands:

### add
```
task-tracker add "an example description"
```

This adds a new task with the given description.

### list
```
task-tracker list [options]
```

Use list to view all tasks that are stored. Follow this command with optional flags like `-done`, `-inprogress`, or `-todo` to view only the tasks with the given status.

### mark
```
task-tracker mark <ID> [options]
```

Use mark to change the status of the task with the given ID from "todo" to "in progress" to "done". Pass an optional flag (`-done`, `-inprogress`, `-todo`) to set the desired status directly.

### delete
```
task-tracker delete <ID>
```

Pass an ID to delete to remove a task from the list.

### update
```
task-tracker update <ID> "a new description"
```

To change the description of a task, pass an ID and new description to to the update command.
