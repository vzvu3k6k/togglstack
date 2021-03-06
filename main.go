package main

import (
	"errors"
	"fmt"
	"github.com/gedex/go-toggl/toggl"
	"github.com/typester/go-pit"
	"math"
	"os"
	"strconv"
	"strings"
)

func getTogglToken() (string, error) {
	profile, err := pit.Get("toggl.com", pit.Requires{"token": ""})
	if err != nil {
		return "", err
	}
	token := (*profile)["token"]
	if token == "" {
		return "", errors.New("Can't get toggl.com token from pit")
	}
	return token, nil
}

func getCurrentTimeEntry(c *toggl.Client) (*toggl.TimeEntry, error) {
	u := "time_entries/current"
	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	data := new(toggl.TimeEntryResponse)
	_, err = c.Do(req, data)

	return data.Data, err
}

func parsePop(args []string) ([]string, bool, int) {
	if len(args) == 0 || args[0] != "pop" {
		return args, false, 0
	}
	args = args[1:]
	var pop_num = 1
	if len(args) >= 1 {
		if args[0] == "all" {
			pop_num = math.MaxUint8
			args = args[1:]
		} else {
			_pop_num, err := strconv.ParseUint(args[0], 10, 8)
			if err == nil {
				pop_num = int(_pop_num)
				args = args[1:]
			}
		}
	}
	return args, true, pop_num
}

func parsePush(args []string) ([]string, bool, string) {
	if len(args) < 2 || args[0] != "push" {
		return args, false, ""
	}
	return args[2:], true, args[1]
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("(pop (\\d+|all)?)? (push ...)?") // TODO: better description
		return
	}
	args, pop_matched, pop_num := parsePop(args)
	args, push_matched, new_item := parsePush(args)
	if len(args) > 0 {
		fmt.Printf("Unknown arguments (%s)\n", strings.Join(args, " "))
		return
	}

	token, err := getTogglToken()
	if err != nil {
		fmt.Println(err)
		return
	}
	client := toggl.NewClient(token)
	current_entry, err := getCurrentTimeEntry(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	if (pop_matched && pop_num > 0) || push_matched {
		if current_entry != nil {
			client.TimeEntries.Stop(current_entry.ID)
			fmt.Printf("Stop \"%s\"\n", current_entry.Description)
		} else {
			if (pop_matched) {
				fmt.Println("There is no time entry running.")
			}
		}
	}

	const separator = ": "
	var stack []string
	if current_entry != nil {
		stack = strings.Split(current_entry.Description, separator)
		if pop_num > 0 {
			if pop_num > len(stack) {
				pop_num = len(stack)
			}
			stack = stack[0 : len(stack) - pop_num]
		}
	}
	if push_matched {
		stack = append(stack, new_item)
	}
	if len(stack) > 0 {
		new_description := strings.Join(stack, separator)
		te := &toggl.TimeEntry{
			CreatedWith: "togglstack",
			Description: new_description,
			WorkspaceID: 0,
		}
		_, err := client.TimeEntries.Start(te)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Start \"%s\"\n", new_description)
	} else {
		fmt.Println("Empty")
	}
}
