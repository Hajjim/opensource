package main

//!!! nos packages sont des packages internes !
import (
	//!!! IL FAUT DONNER LE CHEMIN D'ACCES AU PACKAGE POUR EVITER DE LE RECOPIER 
	//POUR CHAQUE MISE A JOUR 
	//Nous avons choisis un mauvais emplacement et donc on doit chaque fois modifier
	//nos chemins d'accès car j'ai Hajjim avant ISIBOpen-source par exemple
	cedPack "github.com/Hajjim/ISIBOpen-source/www/src/cedricPackage"
	"html/template"
	"log"
	"os"
	//!!! installer le package gofeed de git 'go get -v github.com/mmcdole/gofeed pour ce package

	"net/http"
	"time"
)

//variables globales des structures des données à envoyer dans le html
var dTI DataToInsert
var cD CedData

func main() {
	//lancement du module stib dans un thread à part
	//!! quand notre module termine sa maj, il faut relancer la func createHTML()
	go stibModule()

	//boucle infinie pour que quand une erreur liée à une requête qui dépasse le timeout
	//survient, relance le programme
	startListening()
}

func startListening() {
	//hajjiPack.StartModule()
	//lancement du serveur web
	http.Handle("/", http.FileServer(http.Dir("./HTTP")))
	http.ListenAndServe(":8000", nil)
	/*server := &http.Server{}
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{IP: net.IPv4(192, 168, 11, 2), Port: 8000})
	if err != nil {
		log.Fatal("error creating listener")
	}
	server.Serve(listener)*/
}

//structure envoyée dans le HTML
type DataToInsert struct {
	CedData
	BelData
	HajData
}

//structure liée à la Stib
type CedData struct {
	For9201  string
	For9202  string
	For9301  string
	For9302  string
	Scha9201 string
	Scha9202 string
	Scha9301 string
	Scha9302 string
}

//vos données ici
type BelData struct {
}
type HajData struct {
}

//lance le module de la stib
func stibModule() {
	//lance une fois avant de boucler sinon met 41 sec à démarrer
	loadStibData()
	//renouvelle les valeurs des variables toutes les 41 secondes
	//la stib refresh toutes les 20 secondes mais des requêtes trop fréquentes entrainent
	//des erreurs dans les réponses
	for range time.Tick(time.Second * 41) {
		loadStibData()
	}
}

func loadStibData() {
	//actionne la réception des données
	cedPack.GetVariablesFromServer()
	//enregistre les données dans la structure
	cD = CedData{
		For9201:  cedPack.HorraireFor92[0],
		For9202:  cedPack.HorraireFor92[1],
		For9301:  cedPack.HorraireFor93[0],
		For9302:  cedPack.HorraireFor93[1],
		Scha9201: cedPack.HorraireScha92[0],
		Scha9202: cedPack.HorraireScha92[1],
		Scha9301: cedPack.HorraireScha93[0],
		Scha9302: cedPack.HorraireScha93[1],
	}
	//quand les données sont récupérées, création du html
	createHTML()
}

//crée un fichier index.html sur base du template et des variables à envoyer
func createHTML() {
	//chargement du template
	tpl, err := template.ParseFiles("./HTTP/templates/template.html")
	if err != nil {
		log.Fatalln(err)
	}

	//remplissage des valeurs récupérées par les modules
	dTI = DataToInsert{
		CedData: cD,
	}

	//chargement du html de sortie
	Output, err := os.Create("HTTP/index.html")
	if err != nil {
		log.Println("Erreur lors de la création du fichier", err)
	}
	//création du HTML de sortie en insériant la structure contenant les variables
	tpl.Execute(Output, &dTI)
}
