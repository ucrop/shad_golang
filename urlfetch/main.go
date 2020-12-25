package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	argsWithProg := os.Args
	argsWithProg = argsWithProg[1:]

	for _, url := range argsWithProg {
		err := DownloadFile(url)
		if err != nil {
			fmt.Println("The requested URL was not found.")
			os.Exit(1)
		}
	}
}

func DownloadFile(url string) (err error) {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	// check err
	fmt.Printf("%s\n", b)
	return err
}
