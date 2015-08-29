# togglstack

Togglstack is a cli client for toggl with stack-based timetracking, written in Go.

## Usage

```
$ togglstack push "Write togglstack"
# Working on "Write togglstack".
$ togglstack push "Find a toggl API library in Go"
# Working on "Write togglstack: Find a toggl API library in Go".
$ togglstack pop
# Working on "Write togglstack".
$ togglstack push "Option parser"
# Working on "Write togglstack: Option parser".
$ togglstack push "Implement roughly"
# Working on "Write togglstack: Option parser: Implement roughly".
$ togglstack pop 2 push "Write readme"
# Working on "Write togglstack: Write readme".
$ togglstack pop all
# Empty
````

## License

CC0
