# links_scraper
Webscraper that gathers all links and all those links' links on any pages you tell it to

## Setup
1) Ensure your go environment path variables are set
2) run `go get github.com/mgsterling11/links_scraper`
3) cd into src/github.com/mgsterling11/links_scraper
4) run `go build`
5) run `go run main.go http://github.com`
6) you can scrape links from multiple initial urls by separating urls with a space. run: `go run main.go http://github.com http://facebook.com`

### About the scraper
This scraper runs in two phases. First, it fires a goroutine for each url you pass in as an argument. Values are returned through channels and collected into a slice; once all concurrent goroutines are wrapped up, the collected urls are then each passed again to a new round of goroutines, collecting all hrefs on the next urls' pages.

Once both concurrent rounds of goroutines wrap up, the number of links found on the user-supplied pages are printed, as well as the number of links found in the proceeding pages.  

This is currently two levels deep; future iteration will allow user to determine how many levels deep they'd like to go!
