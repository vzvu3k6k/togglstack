# togglstack

Togglstack is a cli client for toggl with stack-based timetracking, written in Go.

## Usage

```
$ export TOGGL_TOKEN="..."
$ togglstack push "Write togglstack"
# Working on "Write togglstack".
$ togglstack push "Find a toggl API library in Go"
# Working on "Write togglstack: Find a toggl API library in Go".
$ togglstack pop
# Working on "Write togglstack".
$ togglstack push "Option parser"
# Working on "Write togglstack: Option parser".
````
