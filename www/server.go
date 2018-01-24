/*Copyright <2018> <Belataris Dhoulkifli, Hajji Mohamed, Goossens Cédric>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package main

import (
	"fmt"

	/*IL FAUT DONNER LE CHEMIN D'ACCES AU PACKAGE POUR EVITER DE LE RECOPIER POUR
	CHAQUE MISE A JOUR !
	->Nous avons choisis un mauvais emplacement et donc on doit chaque fois modifier
	nos chemins d'accès, donc il faut garder celui-ci github.com/isib pour tous
	cela evitera de modifier chaque fois nos chemins d'accès ou encore de copier
	le package dans le GOPATH*/

	//nos packages internes :
	"time"

	cedPack "github.com/isib/ISIBOpen-source/www/src/cedricPackage"
	hajjiPack "github.com/isib/ISIBOpen-source/www/src/hajjiPackage"
	BelPack "github.com/isib/ISIBOpen-source/www/src/BelatarisPackage"
	//DONT FORGET : Installer mon package gofeed avec "go get -v github.com/mmcdole/gofeed"
	"github.com/mmcdole/gofeed"

	"html/template"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

//variables globales des structures des données à envoyer dans le html
var dTI DataToInsert
var cD CedData
var hD HajjiData
var Bel BelData
var eri EricData
var abdou AbdouData

func main() {
	
	//lancement de nos modules dans un thread à part
	//!! quand notre module termine sa maj, il faut relancer la func createHTML()
	go footModule()
	go stibModule()
	go rssModule()

	//boucle infinie pour que quand une erreur liée à une requête qui dépasse le timeout
	//survient, relance le programme
	startListening()
}

func startListening() {
	//lancement du serveur web
	http.Handle("/", http.FileServer(http.Dir("./HTTP")))
	http.Handle("/socket", websocket.Handler(socket))
	http.ListenAndServe(":8000", nil)

}

//structure envoyée dans le HTML
type DataToInsert struct {
	CedData
	HajjiData
	BelData
	EricData
	AbdouData
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

//structure pour mes données items de mon fluxrss
type HajjiData struct {
	Items []*gofeed.Item
}

//Structure relative à mes données sur le foot
type BelData struct {
	MatchesJournee  BelPack.Journee //Ensemble des matches de la journée
	JourneeActuelle string
}

//@Eric tes données ici
type EricData struct {
}

//@Abdoulilo tes données ici
type AbdouData struct {
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

//lance le module fluxrss
func rssModule() {
	hD.Items = hajjiPack.GetFeedRss() //mon loadRssData en quelque sorte
	for range time.Tick(time.Second * 60) {
		hD.Items = hajjiPack.GetFeedRss()
	}
}

//Lance le module qui va récupérer les données des matches
func footModule() {
	//récupération de mes données une première fois
	BelPack.GetFootballData()
	//Affectation des données récupérer dans la structure qui va être envoyé dans le html
	Bel = BelData{
		MatchesJournee:  BelPack.Jou,
		JourneeActuelle: BelPack.CurrentMatchday}
	createHTML()
	//Rafraichissement des données toutes les 5 minutes (des requêtes trop fréquentes font planter l'API)
	for range time.Tick(time.Minute * 5) {
		BelPack.GetFootballData()

		Bel = BelData{
			MatchesJournee:  BelPack.Jou,
			JourneeActuelle: BelPack.CurrentMatchday}
		createHTML()
		fmt.Println("updated")
	}

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
		CedData:   cD,
		HajjiData: hD,
		BelData:   Bel,
	}

	//chargement du html de sortie
	Output, err := os.Create("HTTP/index.html")
	if err != nil {
		log.Println("Erreur lors de la création du fichier", err)
	}
	//création du HTML de sortie en insériant la structure contenant les variables
	tpl.Execute(Output, &dTI)
}

  type message struct {
	// the json tag means this will serialize as a lowercased field
	Message string `json:"message"`
  }
  //fonction parcourue lorsqu'un client envoi un message au serveur
  func socket(ws *websocket.Conn) {
	for {
	  // allocate our container struct
	  var m message
  
	  // si il recoit un message
	  if err := websocket.JSON.Receive(ws, &m); err != nil {
		log.Println(err)
		break
	  }
  
	  log.Println("Received message:", m.Message)
  
	  // alors envoi une réponse à celui qui a envoyé le message
	  //envoi la structure contenant les données de la stib
	  m2 := cD;
	  if err := websocket.JSON.Send(ws, m2); err != nil {
		log.Println(err)
		break
	  }
	}
  }
  