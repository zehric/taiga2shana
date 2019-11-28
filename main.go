package main

import (
	"flag"
)

var username string
var password string
var filename string
var autoremove bool

func main() {
	flag.StringVar(&username, "user", "", "username for ShanaProject; if not specified,"+
		" will prompt user")
	flag.StringVar(&password, "pass", "", "password for ShanaProject; if not specified,"+
		" will prompt user")
	flag.StringVar(&filename, "list", "", "custom anime list location if user does not have Taiga")
	flag.BoolVar(&autoremove, "autoremove", false, "automatically remove all anime in follows but not in animelist")
	flag.Parse()
	Login()

	var names []DBAnime

	if filename == "" {
		names = ReadTaigaList()
		if names != nil {
			GetResolutions() // fills Resolutions array
		} else {
			names = ReadCustomList("anime.txt")
		}
	} else {
		names = ReadCustomList(filename)
	}

	if names == nil {
		PrintAndExit("I can't seem to find a valid Taiga installation or anime list on your computer. Please provide an " +
			"anime list file manually with the -list option, or create a file called 'anime.txt' in the " +
			"same directory as this program. The anime list should be a newline separated list " +
			"of anime names or search terms.")
	}

	ids := GetAnimeIds(names)

	follows := GetFollows()

	AddAnime(ids, follows)

	RemoveAnime(ids, autoremove)
}
