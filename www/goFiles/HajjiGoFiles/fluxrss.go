package main

import ( 
	"fmt"
	"github.com/mmcdole/gofeed"
)

func main() {

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://www.lemonde.fr/rss/une.xml" )
	fmt.Println(feed.Title)

}
