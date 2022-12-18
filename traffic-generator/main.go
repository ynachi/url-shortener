package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

// shorten issues a shortening request for the URL https://www.golang.org/{ID}.
// ID is a number used to simulate different URLs.
func shorten(url string, ID int) {
	url = fmt.Sprintf("%s/url/create?url=http://www.golang.com/%d", url, ID)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)
}

// worker instantiates a worker to issue encoding requests to the url shortening server
func worker(url string, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range ch {
		shorten(url, j)
	}
}

// main run the traffic generator.
// arg 1: base url, arg 2: number of requests, arg 3: number of workers
func main() {
	var wg sync.WaitGroup
	urlBase := flag.String("api", "https://www.golang.com", "api endpoint")
	numReq := flag.Int("requests", 2000, "total number of requests to issue")
	numWorkers := flag.Int("workers", 50, "number of concurrent workers")
	flag.Parse()
	wg.Add(*numWorkers)
	jobs := make(chan int, *numReq)

	// create workers
	for i := 0; i < *numWorkers; i++ {
		go worker(*urlBase, jobs, &wg)
	}
	// fill working queue (provide jobs)
	for i := 0; i < *numReq; i++ {
		jobs <- i
	}
	close(jobs)
	fmt.Println("Waiting for all the requests to complete")
	wg.Wait()
	fmt.Println("Done")
}
