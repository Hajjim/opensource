package main //je déclare package comme en  java

import ( //je commence par import les librairies fournient par golang
        "encoding/xml" //pour l'encodage avec xml
        "fmt"          //pour l'utilisation de print
        "io/ioutil"    //pour lire les fichiers
        "log"          //pour les log
        "net/http"     //pour requete internet
        "os"           //pour communiquer avec le sys exploit
)

func getContent(url string) ([]byte, error) {
//donc sur cette fonction, je vais recuperer le contenu par rapport à l'url que je passerai sur le m$
        resp, err := http.Get(url) //http.Get renvoie deux valeur : une reponse si bien passé et une$
        //si j'avais mis ''_'' alors ça ne prendrait pas en compte
        //:= quand on initialise une variable (peu import le type) 
        //donc ma fonction pour effectuer une requete get à l'url
        if err != nil { //plus tard supprime l'objet apres que la fonction est finis (return)
                return nil, fmt.Errorf("GET error: %v", err)//%v remplacer par error
        }//donc si y a une erreur, je renvoie l'erreur et []byte va etre nil
        defer resp.Body.Close() //supprime corps de la reponse pour pas etre ouvert après le return

        if resp.StatusCode != http.StatusOK { //Vérifie status code (si différent de 200 alors renvo$
                return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
        }

        data, err := ioutil.ReadAll(resp.Body) //lis le contenu de body + err redefinis
        if err != nil { //erreur au niveau de lecture
		return nil, fmt.Errorf("Read body: %v", err)
	}

	return data, nil//renvoie nil car si err different de nil alors erreur
}

type Enclosure struct {
	Url    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}
/*
type Enclosure struct { // traiter à part car les données trier dans la balise particulière de xml         //struct = representation d'un objet auquel on donne methode
        Url     string  `xml:"url,attr"` 
        Length  int64   `xml:"length,attr"`
        Type    string  `xml:"type,attr"`
//nom type TAG en golang c'est plus simple car represente un truc dans un autre dans un autre type d$
*/
//je recupere chaque item de lemonde RSS :

type Item struct {
	XMLName     xml.Name  `xml:"item"`
	Description string    `xml:"description"`
	Link        string    `xml:"link"`
	PubDate     string    `xml:"pubDate"`
	Title       string    `xml:"title"`
	Enclosure   Enclosure `xml:"enclosure"`
	Guid        string    `xml:"guid"`
}

type Channel struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Desc    string `xml:"description"`
	Items   []Item `xml:"item"`
	PubDate string `xml:"pubDate"`
}

/*
 Represent a xml file
 Each line is a balise
*/
type Rss struct {
	Channel Channel `xml:"channel"`
}

func main() {
	xmlStr, err := getContent("http://www.lemonde.fr/rss/une.xml")
	if err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		log.Println("Received XML:")
		log.Println(xmlStr)
	}
	v := Rss{}
	err = xml.Unmarshal(xmlStr, &v)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	for _, item := range v.Channel.Items {
		log.Println("Title: ", item.Title)
		log.Println("Description: ", item.Description)
		log.Println("Link: ", item.Link)
		log.Println("PubDate: ", item.PubDate)
	}
}
