package main

import (
	BelPack "BelatarisPackage"
	"fmt"
	"html/template"
	"log"
	"os"
)

type BelData struct {
	MatchesJournee  BelPack.Journee
	JourneeActuelle string
}

var Bel BelData

func main() {

	tpl, err := template.ParseFiles("tpl.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	BelPack.GetFootballData()
	Bel = BelData{
		MatchesJournee:  BelPack.Jou,
		JourneeActuelle: BelPack.CurrentMatchday}

	fmt.Println(Bel.JourneeActuelle)
	//fmt.Println(GetNbrsDeMatch(jsonParsed))

	fmt.Println(Bel.MatchesJournee.Matches[0].Date)
	fmt.Println(Bel.MatchesJournee.Matches[0].Image1)
	fmt.Println(Bel.MatchesJournee.Matches[0].Image2)
	Output, err := os.Create("MyPage.html")
	if err != nil {
		log.Println("Erreur lors de la création du fichier", err)
	}
	tpl.Execute(Output, Bel)
}
