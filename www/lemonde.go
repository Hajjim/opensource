package main //je déclare package comme en  java

import ( //je commence par import les librairies fournient par golang
	"encoding/xml" //pour l'encodage avec xml
	"fmt"	       //pour l'utilisation de print
	"io/ioutil"    //pour lire les fichiers
         	"net/http"     //pour requete internet
	"os"	       //pour communiquer avec le sys exploit
)

func getContent(url string) ([]byte, error) {
//donc sur cette fonction, je vais recuperer le contenu par rapport à l'url que je passerai sur le main 
	resp, err := http.Get(url) //http.Get renvoie deux valeur : une reponse si bien passé et une si pas bien passé
	//si j'avais mis ''_'' alors ça ne prendrait pas en compte
	//:= quand on initialise une variable (peu import le type) 
	//donc ma fonction pour effectuer une requete get à l'url
	if err != nil {
		return nil, fmt.Errorf("Get error: %v", err)//%v remplacer par error
	}//donc si y a une erreur, je renvoie l'erreur et []byte va etre nil
	defer resp.Body.Close() //supprime corps de la reponse pour pas etre ouvet après le return

	if resp.StatusCode !=  // FAUT VERIFIER LE STATUS CODE
// ?? docu 3 a lire

type Enclosure struct { // traiter à part car les données trier dans la balise particulière de xml
	//struct = representation d'un objet auquel on donne methode
	Url	string
	Length	string	'xml:"le
	Type	string

type Item struct {}

type Channel struct {}

func main() {
	xmlStr, err := getContent("http://www.lemonde.fr/rss/une.xml")
	//if err != nil {
	//	
