// MODULE pour le flux rss de Le Monde 

//execution du paquet main 
package main
//j'importe tout ce qui me faut pour le xml
import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
)
//mes items se trouvant dans channel
type Item struct {
	XMLName     xml.Name `xml:"item"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	PubDate     string   `xml:"pubDate"`
	Title       string   `xml:"title"`
}
//mon channel se trouvant dans rss
type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Items   []Item   `xml:"channel>item"`
}
//ma fonction principale
func main() {
	//Ioutil.Readfile lit le fichier nommé et retourne le contenu
	file, err := ioutil.ReadFile("/go-work/src/github.com/Hajjim/ISIBOpen-source/www/fluxWiki.rss")
 	//Lorsqu'un appel est réussi, err renvoie une valeur non nul, on écrit
	if err != nil {
 		log.Println(err)
 		return
 	}
	v := Rss{} //:= pour les variables déclarées dans les func
	//J'utilise la fonction Unmarshal de Golang pour qu'il analyse les données codées en XML et stocke le résultat dans la valeur pointée par v
	err = xml.Unmarshal(file, &v)
	//en cas d'erreur
	if err != nil {// nil est le zero pour pointeur/interface/etc
        log.Println(err.Error())
        os.Exit(1) // je sors
    }
	//j'écris tout dans les log avec ma range d'items
	// ''_'' est le blank identifier que j'utilise car je m'en fou du nombre
	for _, item := range v.Items {
		//j'écris titre = ... etc
		log.Println("Title: ", item.Title)
		log.Println("Description: ", item.Description)
		log.Println("Link: ", item.Link)
		log.Println("PubDate: ", item.PubDate)
	}
}
