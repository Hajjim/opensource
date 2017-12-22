package rssPackage

import ( 
	"fmt"
	"github.com/mmcdole/gofeed"
	"html/template"
	"os"
	"log"
//	."hajjiPackage"
)

func StartModule() {
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
