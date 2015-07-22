package main

import (
	"fmt"
	"github.com/jason0x43/go-toggl"
	"math"
	"os"
	"strconv"
	"strings"
)

func getCurrentTimeEntry(session toggl.Session) toggl.TimeEntry {
	account, _ := session.GetAccount()
	var current_entry toggl.TimeEntry
	for _, v := range(account.Data.TimeEntries) {
		if v.IsRunning() {
			current_entry = v
			break
		}
	}
	return current_entry
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

	token := os.Getenv("TOGGL_TOKEN")
	if token == "" {
		fmt.Println("TOGGL_TOKEN is not found.")
		return
	}
	session := toggl.OpenSession(token)
	current_entry := getCurrentTimeEntry(session)

	if (pop_matched && pop_num > 0) || push_matched {
		if current_entry.IsRunning() {
			session.StopTimeEntry(current_entry)
			fmt.Printf("Stop \"%s\"\n", current_entry.Description)
		} else {
			if (pop_matched) {
				fmt.Println("There is no time entry running.")
			}
		}
	}

	const separator = ": "
	var stack []string
	if current_entry.Description != "" {
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
		session.StartTimeEntry(new_description)
		fmt.Printf("Start \"%s\"\n", new_description)
	} else {
		fmt.Println("Empty")
	}
}
