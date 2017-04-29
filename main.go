package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
)

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	return
}

func retrieveUrls(url string, urlChannel chan string, channelFinished chan bool) {
	response, err := http.Get(url)
	z := html.NewTokenizer(response.Body)

	defer func() {
		channelFinished <- true
	}()

	defer response.Body.Close()

	if err != nil {
		fmt.Println("Failed to gather links!")
		return
	}

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			isHref, hrefValue := getHref(t)
			if !isHref {
				continue
			}

			isHttpLink := strings.Index(hrefValue, "http") == 0

			if isHttpLink {
				urlChannel <- hrefValue
			}
		}
	}
}

func fireURlWorkers(urls []string) []string {
	var foundUrls []string
	urlChannel := make(chan string)
	channelFinished := make(chan bool)

	for _, url := range urls {
		go retrieveUrls(url, urlChannel, channelFinished)
	}

	for i := 0; i < len(urls); {
		select {
		case url := <-urlChannel:
			foundUrls = append(foundUrls, url)
		case <-channelFinished:
			i++
		}
	}

	defer close(urlChannel)
	return foundUrls
}

func printUrls(urls []string) {
	for _, url := range urls {
		fmt.Println(" - " + url)
	}
}

func main() {
	urls := fireURlWorkers(os.Args[1:])
	originalUrlsCount := len(urls)
	printUrls(urls)

	urls = fireURlWorkers(urls)
	printUrls(urls)

	fmt.Println("\nFound", originalUrlsCount, "unique urls from original url(s) provided\n")
	fmt.Println("\nFound", len(urls), "additional unique urls found on linked pages\n")
}
