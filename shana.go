package main

import (
	"fmt"
	"github.com/headzoo/surf/browser"
	"github.com/howeyc/gopass"
	"gopkg.in/headzoo/surf.v1"
	"os"
	"strconv"
	"strings"
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

func GetFollows() (follows map[int]Anime) {
	follows = make(map[int]Anime)
	fmt.Println()
	fmt.Println("#### Finding currently followed anime ####")
	bow.Open(basePath + "follows/list/")
	links := bow.Links()
	for _, link := range links {
		url := strings.Replace(link.URL.Path, "/series/", "", -1)
		url = strings.Replace(url, "/", "", -1)
		num, err := strconv.Atoi(url)
		if len(url) > 0 && err == nil {
			follows[num] = Anime{
				Id:    num,
				Value: link.Text,
			}
			fmt.Printf("INFO: Found %s\n", link.Text)
		}
	}
	return
}

func AddAnime(ids map[int]Anime, follows map[int]Anime) {
	fmt.Println()
	fmt.Println("#### Beginning addition to ShanaProject follows ####")
	for _, anime := range ids {
		_, ok := follows[anime.Id]
		if ok {
			fmt.Printf("INFO: %s already in follows, skipping\n", anime.Value)
			continue
		}
		bow.Open(basePath + "follows/add/")
		fmt.Printf("INFO: add %s\n", anime.Value)
		fm, err := bow.Form("form.form")
		if err != nil {
			panic(err)
		}

		/* TODO: make these customizable from the command line */
		fm.Input("title", strconv.Itoa(anime.Id))
		fm.SelectByOptionLabel("subber_tag", "Don't Care")
		if anime.LowQ {
			fm.SelectByOptionLabel("quality_preference", "720p Only (HD)")
		} else {
			fm.SelectByOptionLabel("quality_preference", "1080p Only (HD)")
		}
		fm.SelectByOptionLabel("profile_preference", "Prefer Hi10P")
		fm.SelectByOptionLabel("source_preference", "Any")
		fm.SelectByOptionLabel("back_date", "Retroactively match all existing releases")
		fm.UnCheck("get_any_subber")
		fm.UnCheck("get_any_quality")

		if err = fm.Submit(); err != nil {
			panic(err)
		}
	}
}

func RemoveAnime(ids map[int]Anime, autoremove bool) {
	fmt.Println()
	fmt.Println("#### Beginning removal of ShanaProject follows ####")
	bow.Open(basePath + "follows/list/")
	links := bow.Links()
	var lastAnime Anime
	for _, link := range links {
		url := strings.Replace(link.URL.Path, "/series/", "", -1)
		url = strings.Replace(url, "/", "", -1)
		num, err := strconv.Atoi(url)
		if len(url) > 0 && err == nil {
			lastAnime.Id = num
			lastAnime.Value = link.Text
		}
		url = strings.Replace(link.URL.Path, "/follows/series/", "", -1)
		url = strings.Replace(url, "/", "", -1)
		num, err = strconv.Atoi(url)
		if len(url) > 0 && err == nil {
			_, ok := ids[lastAnime.Id]
			if !ok {
				var input string
				if !autoremove {
					fmt.Printf("INFO: %s is in follows but not in anime list, remove? y/[N] ",
						lastAnime.Value)
					fmt.Scanln(&input)
				}
				if autoremove || strings.ToLower(input) == "y" {
					var token string
					for _, cookie := range bow.SiteCookies() {
						if cookie.Name == "csrftoken" {
							token = cookie.Value
						}
					}
					bow.AddRequestHeader("X-CSRFToken", token)
					body := fmt.Sprintf("id=%d", num)
					bow.Post("http://www.shanaproject.com/ajax/delete_follow/",
						"application/x-www-form-urlencoded", strings.NewReader(body))
					fmt.Printf("INFO: removed %s\n", lastAnime.Value)
				}
			}
		}
	}
}
