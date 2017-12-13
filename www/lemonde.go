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
type Item struct {

// FAIRE XNLNAME ETC

}

type Channel struct {}

func main() {
	xmlStr, err := getContent("http://www.lemonde.fr/rss/une.xml")
	//if err != nil {
	//	
