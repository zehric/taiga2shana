package main

import (
	"fmt"
	"github.com/headzoo/surf/browser"
	"github.com/howeyc/gopass"
	"gopkg.in/headzoo/surf.v1"
	"os"
	"strconv"
)

const basePath = "https://www.shanaproject.com/"

var bow *browser.Browser

func Login() {
	bow = surf.NewBrowser()
	bow.Open(basePath + "login/")
	fm, err := bow.Form("form.form")
	if err != nil {
		panic(err)
	}

	if username == "" {
		fmt.Printf("ShanaProject Username: ")
		fmt.Scanln(&username)
	}
	fm.Input("username", username)

	if password == "" {
		fmt.Printf("ShanaProject Password: ")
		pwbytes, _ := gopass.GetPasswd()
		password = string(pwbytes)
	}
	fm.Input("password", password)
	if err := fm.Submit(); err != nil {
		panic(err)
	}

	bow.Open(basePath + "follows/list/")
	if bow.Url().Path != "/follows/list/" {
		fmt.Println("ERROR: Incorrect password.")
		os.Exit(1)
	}
}

func AddAnime(ids []Anime) {
	fmt.Println()
	fmt.Println("#### Begining addition to ShanaProject follows ####")
	for _, anime := range ids {
		bow.Open(basePath + "follows/add/")
		fmt.Printf("INFO: add %s\n", anime.Value)
		fm, err := bow.Form("form.form")
		if err != nil {
			panic(err)
		}

		/* TODO: make these customizable from the command line */
		fm.Input("title", strconv.Itoa(anime.Id))
		fm.SelectByOptionLabel("subber_tag", "Don't Care")
		fm.SelectByOptionLabel("quality_preference", "1080p Only (HD)")
		fm.SelectByOptionLabel("profile_preference", "Don't Care")
		fm.SelectByOptionLabel("source_preference", "Any")
		fm.SelectByOptionLabel("back_date", "Retroactively match all existing releases")
		fm.UnCheck("get_any_subber")
		fm.UnCheck("get_any_quality")

		if err = fm.Submit(); err != nil {
			panic(err)
		}
	}
}
