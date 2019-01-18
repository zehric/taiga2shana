package main

import (
	"flag"
	"fmt"
	"os"
)

var username string
var password string
var filename string

func main() {
	flag.StringVar(&username, "user", "", "username for ShanaProject; if not specified,"+
		" will prompt user")
	flag.StringVar(&password, "pass", "", "password for ShanaProject; if not specified,"+
		" will prompt user")
	flag.StringVar(&filename, "list", "", "custom anime list location if user does not have Taiga")
	flag.Parse()
	Login()

	var names []string

	if filename == "" {
		names = ReadTaigaList()
		if names == nil {
			fmt.Println("I can't seem to find a valid Taiga installation on your computer. Please provide an " +
				"anime list file manually with the -list option. The anime list should be a newline separated list " +
				"of anime names or search terms.")
			os.Exit(1)
		}
	} else {
		names = ReadCustomList(filename)
	}
	ids := GetAnimeIds(names)

	follows := GetFollows()

	AddAnime(ids, follows)
}
