package handler

import "strings"

// Handle @@##
func Handle(book, entryStr string) (count uint64) {
	for {
		i := strings.Index(book, entryStr)
		if i == -1 {
			return
		}
		count++
		book = book[i+len(entryStr):]
	}
}
