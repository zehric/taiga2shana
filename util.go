package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

func GetUserSelection(selections []string) (num int) {
	for i, match := range selections {
		fmt.Printf("%d: %v\n", i, match)
	}
	var err error
	for {
		fmt.Printf("Please enter the number corresponding to the desired entry (%d-%d) or -1 if none match: ",
			0, len(selections)-1)
		var selection string
		fmt.Scanln(&selection)
		num, err = strconv.Atoi(selection)
		if err == nil && num < len(selections) {
			return
		}
	}
}

func ReadCustomList(filename string) (names []DBAnime) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	var dbElem DBAnime

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dbElem.Title = strings.TrimSpace(scanner.Text())
		dbElem.Japanese = ""
		names = append(names, dbElem)
	}

	if err := scanner.Err(); err != nil {
		PrintAndExit(err.Error())
	}
	return
}

func isDoubleClickRun() bool {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	lp := kernel32.NewProc("GetConsoleProcessList")
	if lp != nil {
		var pids [2]uint32
		var maxCount uint32 = 2
		ret, _, _ := lp.Call(uintptr(unsafe.Pointer(&pids)), uintptr(maxCount))
		if ret > 1 {
			return false
		}
	}
	return true
}

func PrintAndExit(message string) {
	fmt.Println(message)
	if runtime.GOOS == "windows" && isDoubleClickRun() {
		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
	os.Exit(1)
}
