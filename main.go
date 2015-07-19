package main

import (
	"fmt"
	"github.com/jason0x43/go-toggl"
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

func main() {
	const separator = ": "

	session := toggl.OpenSession(os.Getenv("TOGGL_TOKEN"))
	current_entry := getCurrentTimeEntry(session)

	if os.Args[1] == "push" {
		if current_entry.IsRunning() {
			session.StopTimeEntry(current_entry)
		}
		new_description := current_entry.Description
		if new_description != "" {
			new_description += separator
		}
		new_description += os.Args[2]
		session.StartTimeEntry(new_description)
		fmt.Printf("Start %s\n", new_description)
		return
	}

	if os.Args[1] == "pop" {
		if !current_entry.IsRunning() {
			fmt.Println("There is no time entry running.")
			return
		}
		session.StopTimeEntry(current_entry)

		items := strings.Split(current_entry.Description, separator)

		var pop_num uint64 = 1
		if len(os.Args) >= 2 {
			pop_num, _ = strconv.ParseUint(os.Args[2], 10, 64)
		}
		if uint64(len(items)) <= pop_num || os.Args[2] == "all" {
			fmt.Println("Done.")
			return
		}
		new_description := strings.Join(items[0:-pop_num], separator)
		session.StartTimeEntry(new_description)
		fmt.Printf("Start %s\n", new_description)
		return
	}
}
