package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// Global variables are allowed (and encouraged) for this project.

// Entry represents a single user's credentials for a site.
type Entry struct {
	User     string
	Password string
}

// EntrySlice is a slice of Entry structs.
type EntrySlice []Entry

// Global passwordMap to store the site to EntrySlice mapping.
var passwordMap map[string]EntrySlice

// _______________________________________________________________________
// initialize before main()
// _______________________________________________________________________
func init() {
	passwordMap = make(map[string]EntrySlice)
	pmRead()
}

// _______________________________________________________________________
// find the matching entry slice
// _______________________________________________________________________
func findEntrySlice(site string) (EntrySlice, bool) {
	entries, exists := passwordMap[site]
	return entries, exists
}

// _______________________________________________________________________
// set the entrySlice for site
// _______________________________________________________________________
func setEntrySlice(site string, entrySlice EntrySlice) {
	passwordMap[site] = entrySlice
	pmWrite() // Save changes to the file after updating the entry slice.
}

// _______________________________________________________________________
// find
// _______________________________________________________________________
func find(user string, entrySlice EntrySlice) (int, bool) {
	for i, entry := range entrySlice {
		if entry.User == user {
			return i, true
		}
	}
	return -1, false // Return -1 if the user is not found.
}

// _______________________________________________________________________
// print the list in columns
// _______________________________________________________________________
func pmList() {
	w := tabwriter.NewWriter(os.Stdout, 8, 8, 1, ' ', 0)
	for site, entries := range passwordMap {
		for _, entry := range entries {
			fmt.Fprintf(w, "%s\t%s\t%s\t\n", site, entry.User, entry.Password)
		}
	}
	w.Flush()
}

// _______________________________________________________________________
//
//	add an entry if the site, user is not already found
//
// _______________________________________________________________________
func pmAdd(site, user, password string) {
	entrySlice, exists := findEntrySlice(site)
	if exists {
		if _, found := find(user, entrySlice); found {
			fmt.Println("add: duplicate entry")
			return
		}
	}
	setEntrySlice(site, append(entrySlice, Entry{User: user, Password: password}))
}

// _______________________________________________________________________
// remove by site and user
// _______________________________________________________________________
func pmRemove(site, user string) {
	entrySlice, exists := findEntrySlice(site)
	if !exists {
		fmt.Println("remove: site not found")
		return
	}

	userIndex, found := find(user, entrySlice)
	if !found {
		fmt.Println("remove: user not found")
		return
	}

	// Remove the entry at index userIndex for site.
	setEntrySlice(site, append(entrySlice[:userIndex], entrySlice[userIndex+1:]...))
}

// _______________________________________________________________________
// remove the whole site if there is a single user at that site
// _______________________________________________________________________
func pmRemoveSite(site string) {
	entrySlice, exists := findEntrySlice(site)
	if !exists {
		fmt.Println("remove: site not found")
		return
	}

	if len(entrySlice) > 1 {
		fmt.Println("attempted to remove multiple users")
		return
	}

	delete(passwordMap, site)
	pmWrite() // Save changes to the file after removing the site.
}

// _______________________________________________________________________
// read the passwordVault
// _______________________________________________________________________
func pmRead() {
	file, err := os.Open("passwordVault")
	if err != nil {
		if os.IsNotExist(err) {
			passwordMap = make(map[string]EntrySlice)
			return
		}
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 3 {
			site, user, password := parts[0], parts[1], parts[2]
			passwordMap[site] = append(passwordMap[site], Entry{User: user, Password: password})
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
}

// _______________________________________________________________________
// write the passwordVault
// _______________________________________________________________________
func pmWrite() {
	file, err := os.OpenFile("passwordVault", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file for writing:", err)
		os.Exit(1)
	}
	defer file.Close()

	for site, entries := range passwordMap {
		for _, entry := range entries {
			if _, err := fmt.Fprintf(file, "%s %s %s\n", site, entry.User, entry.Password); err != nil {
				fmt.Println("Error writing to file:", err)
				os.Exit(1)
			}
		}
	}
}

// _______________________________________________________________________
// do forever loop reading the following commands
//
//	  l
//	  a s u p
//	  r s
//	  r s u
//	  x
//	where l,a,r,x are list, add, remove, and exit
//	and s,u,p are site, user, and password
//
// _______________________________________________________________________
func loop() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "l":
			pmList()
		case "a":
			if len(parts) != 4 {
				fmt.Println("Usage: a site user password")
				continue
			}
			pmAdd(parts[1], parts[2], parts[3])
		case "r":
			if len(parts) == 2 {
				pmRemoveSite(parts[1])
			} else if len(parts) == 3 {
				pmRemove(parts[1], parts[2])
			} else {
				fmt.Println("Usage: r site [user]")
			}
		case "x":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid command")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

// _______________________________________________________________________
//
//	let her rip
//
// _______________________________________________________________________
func main() {
	loop()
}
