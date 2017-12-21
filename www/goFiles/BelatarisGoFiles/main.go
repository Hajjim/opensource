package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Jeffail/gabs"
)

type Match struct {
	Status    string
	Team1     string
	Team2     string
	Resultat1 float64
	Resultat2 float64
	Matchday  float64
	Date      string
}
type Journee struct {
	Matches []Match
}

func (j *Journee) AddMatch(jsonData *gabs.Container) {

	children, _ := jsonData.S("fixtures").Children()
	for _, child := range children {
		s := child.Path("status").Data().(string)
		switch s {
		case "TIMED": //Cas particulier car quand le match n'est pas encore joué les résultats sont à null et donc cela créer des problèmes avec le type float qui ne peut pas contenir la valeur null
			j.Matches = append(j.Matches, Match{child.Path("status").Data().(string), child.Path("homeTeamName").Data().(string),
				child.Path("awayTeamName").Data().(string), 0, 0, child.Path("matchday").Data().(float64), child.Path("date").Data().(string)}) //Ajout des match dans le tableau
		default:
			j.Matches = append(j.Matches, Match{child.Path("status").Data().(string), child.Path("homeTeamName").Data().(string),
				child.Path("awayTeamName").Data().(string), child.Path("result.goalsHomeTeam").Data().(float64), child.Path("result.goalsAwayTeam").Data().(float64), child.Path("matchday").Data().(float64), child.Path("date").Data().(string)})
		}

	}
}

func main() {
	//lesJournees := make([]Journee, 38)
	//Anglais := Championnat{Matchdays}
	var jou Journee
	tpl, err := template.ParseFiles("tpl.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	currentMatchday := getCurrentMatchday()
	/*url := "http://api.football-data.org/v1/competitions/445/fixtures?timeFrameStart=" + dateBeforeNow + "&timeFrameEnd=" + dateNow*/
	url := "http://api.football-data.org/v1/competitions/445/fixtures?matchday=" + currentMatchday
	jsonParsed, _ := gabs.ParseJSON(getDataAPI(url)) //Récupération des données JSON

	fmt.Println(currentMatchday)
	fmt.Println(getNbrsDeMatch(jsonParsed))
	jou.AddMatch(jsonParsed) //Ajout de mes données JSON dans mon objet journee
	fmt.Println(jou.Matches[0].Date)
	Output, err := os.Create("MyPage.html")
	if err != nil {
		log.Println("Erreur lors de la création du fichier", err)
	}
	tpl.Execute(Output, jou.Matches)
}

func getDataAPI(url string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Auth-Token", "286d522f41bf4602b69e9eed00870841")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func getNbrsDeMatch(jsonData *gabs.Container) float64 {
	return jsonData.Path("count").Data().(float64) //mettre float car si je met int, il n'accepte pas

}
func getCurrentMatchday() string {
	url := "http://api.football-data.org/v1/competitions/445/leagueTable"
	jsonParsed2, _ := gabs.ParseJSON(getDataAPI(url))

	return fmt.Sprint(jsonParsed2.Path("matchday").Data().(float64))
}
