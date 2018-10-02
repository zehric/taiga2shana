package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const searchPath = basePath + "ajax/search/title/?term="

type Anime struct {
	Id    int
	Value string
}

func searchAnime(name string) (body []Anime) {
	if len(name) == 0 {
		return
	}

	escaped := url.QueryEscape(name)
	resp, err := http.Get(searchPath + escaped)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == http.StatusOK {
		bytes, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err := json.Unmarshal(bytes, &body); err != nil {
			panic(err)
		}
		if len(body) == 0 && len(name) > 0 {
			name = name[:strings.LastIndex(name, " ")]
			name = strings.TrimSpace(name)
			return searchAnime(name)
		}
	} else {
		fmt.Printf("ERROR: Status code was %d on search\n", resp.StatusCode)
		os.Exit(1)
	}
	return
}

func GetAnimeIds(names []string) (ids []Anime) {
	println("#### Searching ShanaProject for anime in your list ####")
	for _, name := range names {
		body := searchAnime(name)
		if len(body) == 0 {
			fmt.Printf("WARNING: Could not find match on ShanaProject for %s\n", name)
		} else if len(body) > 1 {
			fmt.Printf("INFO: Found multiple matches on ShanaProject for %s\n", name)
			var selections []string
			for _, anime := range body {
				selections = append(selections, anime.Value)
			}
			i := GetUserSelection(selections)
			ids = append(ids, body[i])
		} else {
			if !strings.EqualFold(body[0].Value, name) {
				fmt.Printf("WARNING: matched anime title is not the same as requested title:\n%s\n%s\n",
					body[0].Value, name)
			}
			ids = append(ids, body[0])
		}
	}
	return
}
