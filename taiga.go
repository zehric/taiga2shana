package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var TaigaListFilename string
var TaigaDBFilename string

func findAnimeList() {
	TaigaListFilename = filepath.Join(os.Getenv("APPDATA"), "Taiga", "data", "user")
	dir, err := ioutil.ReadDir(TaigaListFilename)
	if err != nil {
		panic(err)
	}
	var profileName string
	if len(dir) == 0 {
		fmt.Println("INFO: Taiga profile doesn't exist")
	} else if len(dir) > 1 {
		fmt.Println("INFO: Found multiple Taiga profiles")
		var selections []string
		for _, f := range dir {
			selections = append(selections, f.Name())
		}
		profileName = dir[GetUserSelection(selections)].Name()
	} else {
		profileName = dir[0].Name()
	}
	TaigaListFilename = filepath.Join(TaigaListFilename, profileName, "anime.xml")
}

func findDatabase() {
	TaigaDBFilename = filepath.Join(os.Getenv("APPDATA"), "Taiga", "data", "db", "anime.xml")
}

type UserAnime struct {
	XMLName xml.Name `xml:"anime"`
	Id      int      `xml:"id"`
	Status  int      `xml:"status"`
}

type DBAnime struct {
	XMLName xml.Name `xml:"anime"`
	Id      int      `xml:"id"`
	Title   string   `xml:"title"`
}

func ReadTaigaList() []string {
	findAnimeList()
	findDatabase()

	userAnimeXml, err := os.Open(TaigaListFilename)
	if err != nil {
		return nil
	}
	decoder := xml.NewDecoder(userAnimeXml)
	var watchingIds []int
	var userElem UserAnime
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "anime" {
				decoder.DecodeElement(&userElem, &se)
				if userElem.Status == 1 {
					watchingIds = append(watchingIds, userElem.Id)
				}
			}
		}
	}
	userAnimeXml.Close()

	dbAnimeXml, err := os.Open(TaigaDBFilename)
	if err != nil {
		return nil
	}
	decoder = xml.NewDecoder(dbAnimeXml)
	var names []string
	var dbElem DBAnime
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "anime" {
				decoder.DecodeElement(&dbElem, &se)
				for _, id := range watchingIds {
					if dbElem.Id == id {
						names = append(names, dbElem.Title)
						continue
					}
				}
			}
		}
	}
	dbAnimeXml.Close()
	return names
}
