package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const searchPath = basePath + "ajax/search/title/?term="

type Anime struct {
	Id    int
	Value string
	LowQ  bool
}

func searchAnime(name string) (body []Anime) {
	if len(name) < 3 {
		return
	}

	escaped := url.QueryEscape(name)
	resp, err := http.Get(searchPath + escaped)
	if err != nil {
		PrintAndExit(err.Error())
	}
	if resp.StatusCode == http.StatusOK {
		bytes, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err := json.Unmarshal(bytes, &body); err != nil {
			PrintAndExit(err.Error())
		}
		if len(body) == 0 && len(name) > 0 {
			idx := strings.LastIndexAny(name, " `~!@#$%^&*()_+-=[]{}\\|;:'\",<.>/?")
			if idx == -1 {
				return
			}
			name = name[:idx]
			name = strings.TrimSpace(name)
			return searchAnime(name)
		}
	} else {
		PrintAndExit(fmt.Sprintf("ERROR: Status code was %d on search", resp.StatusCode))
	}
	return
}

func searchResolutions(japanese string) string {
	for _, res := range Resolutions {
		if strings.Contains(res.Name, japanese) || strings.Contains(japanese, res.Name) {
			return res.Res
		}
	}
	return ""
}

func GetAnimeIds(names []DBAnime) (ids map[int]Anime) {
	ids = make(map[int]Anime)
	fmt.Println()
	fmt.Println("#### Searching ShanaProject for anime in your list ####")
	for _, name := range names {
		body := searchAnime(name.Title)
		idx := -1
		if len(body) == 0 {
			fmt.Printf("WARNING: Could not find match on ShanaProject for %s\n", name.Title)
			continue
		} else if len(body) > 1 {
			fmt.Printf("INFO: Found multiple matches on ShanaProject for %s\n", name.Title)
			var selections []string
			for _, anime := range body {
				selections = append(selections, anime.Value)
			}
			idx = GetUserSelection(selections)
		} else {
			if !strings.EqualFold(body[0].Value, name.Title) {
				fmt.Printf("WARNING: matched anime title is not the same as requested title:\n%s\n%s\n",
					body[0].Value, name.Title)
			}
			idx = 0
		}
		if idx >= 0 {
			body[idx].LowQ = false
			if searchResolutions(name.Japanese) == "HV1280" {
				body[idx].LowQ = true
			}
			ids[body[idx].Id] = body[idx]
		}
	}
	return
}
