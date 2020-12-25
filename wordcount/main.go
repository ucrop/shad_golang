// +build !solution

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func main() {

	argsWithProg := os.Args
	argsWithProg = argsWithProg[1:]

	countWords := make(map[string]int)

	for _, file := range argsWithProg {
		dat, err := ioutil.ReadFile(string(file))
		check(err)
		words := strings.Fields(string(dat))

		for _, word := range words {
			countWords[word]++
		}
	}

	for k, v := range countWords {
		if v < 2 {
			continue
		}
		fmt.Print(v)
		fmt.Print("\t")
		fmt.Println(k)
	}

}
