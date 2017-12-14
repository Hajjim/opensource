package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Jeffail/gabs"
)

type Match struct {
	Status    string
	Team1     string
	Team2     string
	Resultat1 float64
	Resultat2 float64
	Matchday  float64
}
type Journees struct {
	Matches []Match
}

func (j *Journees) AddMatch(jsonData *gabs.Container) {

	children, _ := jsonData.S("fixtures").Children()
	for _, child := range children {
		s := child.Path("status").Data().(string)
		switch s {
		case "TIMED":
			j.Matches = append(j.Matches, Match{child.Path("status").Data().(string), child.Path("homeTeamName").Data().(string),
				child.Path("awayTeamName").Data().(string), 0, 0, child.Path("matchday").Data().(float64)})
		default:
			j.Matches = append(j.Matches, Match{child.Path("status").Data().(string), child.Path("homeTeamName").Data().(string),
				child.Path("awayTeamName").Data().(string), child.Path("result.goalsHomeTeam").Data().(float64), child.Path("result.goalsAwayTeam").Data().(float64), child.Path("matchday").Data().(float64)})
		}

	}
}
func (j *Journees) display(matchday float64) {
	fmt.Printf("%v ème journée\n", matchday)
	for _, match := range j.Matches {

		/*defer*/ fmt.Printf("[%v]====> %v %v - %v %v   %v\n", match.Status, match.Team1, match.Resultat1, match.Resultat2, match.Team2, match.Matchday)
	}
}
func (j *Journees) Afficher() {
	var ref Match
	//var n int = 1
	var k Journees
	for i, match := range j.Matches {
		if i == 0 {
			ref = match
			k.Matches = append(k.Matches, ref)
		}
		if ref.Matchday != match.Matchday {
			//n++
			/*defer*/ k.display(ref.Matchday)
			ref = match
			k.Matches = []Match{}
			k.Matches = append(k.Matches, ref)
		} else if i != 0 && ref.Matchday == match.Matchday {
			k.Matches = append(k.Matches, match)
		}

	}
	/*defer*/ k.display(ref.Matchday)
	//println(n)
}

func main() {
	dateNow, dateBeforeNow := getDate()
	//fmt.Println(dateNow + dateBeforeNow)
	//lesJournees := make([]Journee, 38)
	//Anglais := Championnat{Matchdays}
	var jou Journees
	tpl, err := template.ParseFiles("tpl.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	currentMatchday := getCurrentMatchday()
	//url := "http://api.football-data.org/v1/competitions/445/fixtures?timeFrameStart=" + "2017-12-10" + "&timeFrameEnd=" + "2017-12-10"
	url := "http://api.football-data.org/v1/competitions/445/fixtures?timeFrameStart=" + dateBeforeNow + "&timeFrameEnd=" + dateNow
	jsonParsed, _ := gabs.ParseJSON(getDataAPI(url))

	fmt.Println(currentMatchday)
	fmt.Println(getNbrsDeMatch(jsonParsed))
	jou.AddMatch(jsonParsed)

	//jou.Afficher()
	Output, err := os.Create("MyPage.html")
	if err != nil {
		log.Println("Erreur lors de la création du fichier", err)
	}
	tpl.Execute(Output, jou.Matches)
}

func getDate() (string, string) {
	now := time.Now()
	beforeNow := now.AddDate(0, 0, -3)
	mask := "2006-01-02"
	return now.Format(mask), beforeNow.Format(mask)
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
func getCurrentMatchday() float64 {
	url := "http://api.football-data.org/v1/competitions/445/leagueTable"
	jsonParsed2, _ := gabs.ParseJSON(getDataAPI(url))

	return jsonParsed2.Path("matchday").Data().(float64)
}
