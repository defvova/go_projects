package main

import (
	"fmt"
	"sync"
	"webscraper/internal"
)

func main() {
	cr := internal.Crawler{
		ValidUrls:   make(map[string]map[string]struct{}),
		InvalidUrls: make(map[string]struct{}),
		Visited:     make(map[string]struct{}),
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	wg.Add(1)
	go cr.TraverseLinks("/", &wg, sem)
	wg.Wait()

	for u := range cr.InvalidUrls {
		baseUrl := cr.BuildUrl(u)
		fmt.Println(baseUrl.String())
	}
}
