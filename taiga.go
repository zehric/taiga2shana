package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var taigaInstallDir string
var taigaListFilename string
var taigaDBFilename string

func findInstallDir() {
	taigaInstallDir = filepath.Join(os.Getenv("APPDATA"), "Taiga", "asdf")
	for _, err := os.Stat(taigaInstallDir); os.IsNotExist(err); _, err = os.Stat(taigaInstallDir) {
		fmt.Printf("INFO: Path does not exist: %s\n", taigaInstallDir)
		fmt.Printf("Please specify the install directory of Taiga: ")
		fmt.Scanln(&taigaInstallDir)
	}
}

func findAnimeList() {
	taigaListFilename = filepath.Join(taigaInstallDir, "data", "user")
	dir, err := ioutil.ReadDir(taigaListFilename)
	if err != nil {
		fmt.Printf("ERROR: Path does not exist: %s\n", taigaListFilename)
		fmt.Println("Make sure your Taiga install directory contains Taiga.exe and a folder called data")
		os.Exit(1)
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
	taigaListFilename = filepath.Join(taigaListFilename, profileName, "anime.xml")
}

func findDatabase() {
	taigaDBFilename = filepath.Join(taigaInstallDir, "data", "db", "anime.xml")
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
	findInstallDir()
	findAnimeList()
	findDatabase()

	userAnimeXml, err := os.Open(taigaListFilename)
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

	dbAnimeXml, err := os.Open(taigaDBFilename)
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
