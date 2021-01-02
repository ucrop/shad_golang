package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	argsWithProg := os.Args
	argsWithProg = argsWithProg[1:]
	var wg sync.WaitGroup
	time0 := time.Now()

	for _, url := range argsWithProg {
		wg.Add(1)
		go func(url string, wg *sync.WaitGroup) {
			defer wg.Done()
			resp, err := http.Get(url)
			duration := time.Since(time0)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			fmt.Printf("%s\t%d\t%s\n", duration, len(body), url)
		}(url, &wg)
	}
	wg.Wait()

	totalDur := time.Since(time0)
	fmt.Printf("%s elapsed\n", totalDur)
}
