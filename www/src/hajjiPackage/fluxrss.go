/*The MIT License (MIT)

Copyright (c) 2018 hajjim

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.**/

package RssModule //au lieu de main

import (
	//	"fmt"
	"github.com/mmcdole/gofeed"
)

func GetFeedRss() []*gofeed.Item {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://www.lemonde.fr/rss/une.xml")
	return feed.Items //je renvoie mes items pour server.go
}

/*
func main() {
//	htmlgo, err := template.ParseFiles("fluxgo.html")
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://www.lemonde.fr/rss/une.xml" )
	fmt.Println(feed.Title)
	fmt.Println(feed.Link)
	fmt.Println(feed.Items[0].Title)
        fmt.Println(feed.Items[0].Description)
        fmt.Println(feed.Items[0].Link)
//         fmt.Println(feed.Items[0].Enclosur)
    //objet rss vide (il parse par rapport au type de la variable qu'on lui donne)
        for _, item := range feed.Items { // for range prend tableau des item de en haut et
        //et il le parcourt (indice et la valeur) comme ci on avait for key qui va  stocker
        //l'indice et l'element (range parcourt tout le tableau)
                fmt.Println("Title: ", item.Title)
                fmt.Println("Description: ", item.Description)
                fmt.Println("Link: ", item.Link)
                fmt.Println("PubDate: ", item.Published)
		for _, enclosure := range item.Enclosures {
		fmt.Println("Enclosure", enclosure.URL)
//		Output, err := os.Create("Mypage.html")
		}
                }
//	htmlgo.Execute(Output, feed.Items)

}
*/
