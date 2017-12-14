package main //je déclare package comme en  java

import ( //je commence par import les librairies fournient par golang
	"encoding/xml" //pour l'encodage avec xml
	"fmt"	       //pour l'utilisation de print
	"io/ioutil"    //pour lire les fichiers
        "log"	       //pour les log
	"net/http"     //pour requete internet
	"os"	       //pour communiquer avec le sys exploit
)

func getContent(url string) ([]byte, error) {
//donc sur cette fonction, je vais recuperer le contenu par rapport à l'url que je passerai sur le main 
	resp, err := http.Get(url) //http.Get renvoie deux valeur : une reponse si bien passé et une si pas bien passé
	//si j'avais mis ''_'' alors ça ne prendrait pas en compte
	//:= quand on initialise une variable (peu import le type) 
	//donc ma fonction pour effectuer une requete get à l'url
	if err != nil { //plus tard supprime l'objet apres que la fonction est finis (return)
		return nil, fmt.Errorf("Get error: %v", err)//%v remplacer par error
	}//donc si y a une erreur, je renvoie l'erreur et []byte va etre nil
	defer resp.Body.Close() //supprime corps de la reponse pour pas etre ouvert après le return

	if resp.StatusCode != //Vérifie status code (si différent de 200 alors renvoie erreur meme si c'est plus au moins ok)
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body) //lis le contenu de body + err redefinis
	if err != nil { //erreur au niveau de lecture
		return nil, fmt.Errorf("Read body: %v", err)
	}

	return data, nil //renvoie nil car si err different de nil alors erreur
}


type Enclosure struct { // traiter à part car les données trier dans la balise particulière de xml
	//struct = representation d'un objet auquel on donne methode
	Url	string	`xml:"url,attr"` 
	Length	string	`xml:"length,attr"`
	Type	string	`xml:"type,attr"`
//nom type TAG en golang c'est plus simple car represente un truc dans un autre dans un autre type de format

//je recupere chaque item de lemonde RSS :

type Item struct {
	XMLNAME		xml.Name	`xml:"item"`// creer  un objet en fct des tag donné
	Description	string		`xml:"description"`
	Link		string		`xml:"link"`
	PubDate		string		`xml:"pubDate"`
	Title		string		`xml:"title"`
	Enclosure	Enclosure	`xml:"enclosure"`//Enclosure on choppe attribut en haut
	Guid		string		`xml:"guid"`

}

//balise channel contient item 

type Channel struct {
	Title		string		`xml:"title"`
	Link		string		`xml:"link"`
	Desc		string		`xml:"description"`
	Items		[]Item		`xml:"item"`
	PubDate		string		`xml:"pubDate"`
}

// ma balise RSS

type RSS struct {
	Channel		Channel		`xml:"channel"`
}

func main() {
	xmlStr, err := getContent("http://www.lemonde.fr/rss/une.xml")
	if err != nil {
		log.Printf("Failed to get XML: %v", err)
	}
	v := RSS {} //objet rss vide (il parse par rapport au type de la variable qu'on lui donne)
	err = xml.Unmarshal(xmlStr, &v) //unmarshal traduire xml vers l'objet donné
	if err!= nil { //unmarshal donne un tableau de byte et la variable ou il doit stocker et il parse
		log.Println(err.Error())//  &v pointe vers v
		os.Exit(1) // si y a erreur je quitte
	}
	for _, item := range v.Channel.Items { // for range prend tableau des item de en haut et
	//et il le parcourt (l'indice et la valeur) comme ci on avait for key qui va  stocker
	//l'indice et l'element (range parcourt tout le tableau)
		log.Println("Title: ", item.Title)
		log.Println("Description: ", item.Description)
		log.Println("Link: ", item.Link)
		log.Println("PubDate: ", item.PubDate)
		}
}
