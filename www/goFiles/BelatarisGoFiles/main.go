package main

import (
	. "BelatarisPackage"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/Jeffail/gabs"
)

func main() {
	var jou Journee
	tpl, err := template.ParseFiles("tpl.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	currentMatchday := GetCurrentMatchday()
	/*url := "http://api.football-data.org/v1/competitions/445/fixtures?timeFrameStart=" + dateBeforeNow + "&timeFrameEnd=" + dateNow*/
	url := "http://api.football-data.org/v1/competitions/445/fixtures?matchday=" + currentMatchday
	jsonParsed, _ := gabs.ParseJSON(GetDataAPI(url)) //Récupération des données JSON

	fmt.Println(currentMatchday)
	fmt.Println(GetNbrsDeMatch(jsonParsed))
	jou.AddMatch(jsonParsed) //Ajout de mes données JSON dans mon objet journee
	fmt.Println(jou.Matches[0].Date)
	Output, err := os.Create("MyPage.html")
	if err != nil {
		log.Println("Erreur lors de la création du fichier", err)
	}
	tpl.Execute(Output, jou.Matches)
}
