package football

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"BelatarisPackage/github.com/Jeffail/gabs"
)

type Match struct {
	Status    string
	Team1     string
	Team2     string
	Resultat1 float64
	Resultat2 float64
	Matchday  float64
	Date      string
	Image1    string //Url vers le logo du club 1
	Image2    string //Url vers le logo du club 2
}
type Journee struct {
	Matches []Match
}

func (j *Journee) AddMatch(jsonData *gabs.Container) {

	children, _ := jsonData.S("fixtures").Children()
	for _, child := range children {
		s := child.Path("status").Data().(string)
		href1 := child.Path("_links.homeTeam.href").Data().(string) //Récupération du lien url vers les infos de l'équipe1
		href2 := child.Path("_links.awayTeam.href").Data().(string) //Récupération du lien url vers les infos de l'équipe2
		jsonParsedTeamInfo1, _ := gabs.ParseJSON(GetDataAPI(href1))
		jsonParsedTeamInfo2, _ := gabs.ParseJSON(GetDataAPI(href2))
		imgUrl1 := jsonParsedTeamInfo1.Path("crestUrl").Data().(string)
		imgUrl2 := jsonParsedTeamInfo2.Path("crestUrl").Data().(string)
		switch s {
		case "TIMED": //Cas particulier car quand le match n'est pas encore joué les résultats sont à null et donc cela créer des problèmes avec le type float qui ne peut pas contenir la valeur null
			j.Matches = append(j.Matches, Match{child.Path("status").Data().(string), ///Ajout des match dans le tableau
				child.Path("homeTeamName").Data().(string),
				child.Path("awayTeamName").Data().(string),
				0, 0,
				child.Path("matchday").Data().(float64),
				child.Path("date").Data().(string), imgUrl1, imgUrl2})

		case "SCHEDULED": //2ème Cas particulier car quand le match n'est pas encore joué les résultats sont à null et dans certains cas ces matchs sont réferencés en tant que SCHEDULED au lieu de TIMED
			j.Matches = append(j.Matches, Match{child.Path("status").Data().(string), ///Ajout des match dans le tableau
				child.Path("homeTeamName").Data().(string),
				child.Path("awayTeamName").Data().(string),
				0, 0,
				child.Path("matchday").Data().(float64),
				child.Path("date").Data().(string), imgUrl1, imgUrl2})
		default: //Matches terminés ou en cours
			j.Matches = append(j.Matches, Match{child.Path("status").Data().(string), child.Path("homeTeamName").Data().(string),
				child.Path("awayTeamName").Data().(string),
				child.Path("result.goalsHomeTeam").Data().(float64),
				child.Path("result.goalsAwayTeam").Data().(float64),
				child.Path("matchday").Data().(float64),
				child.Path("date").Data().(string), imgUrl1, imgUrl2})
		}

	}
}

var Jou Journee
var CurrentMatchday string

func GetFootballData() {

	CurrentMatchday = GetCurrentMatchday()
	/*url := "http://api.football-data.org/v1/competitions/445/fixtures?timeFrameStart=" + dateBeforeNow + "&timeFrameEnd=" + dateNow*/
	url := "http://api.football-data.org/v1/competitions/445/fixtures?matchday=" + CurrentMatchday
	jsonParsed, _ := gabs.ParseJSON(GetDataAPI(url)) //Récupération des données JSON
	Jou.AddMatch(jsonParsed)                         //Ajout de mes données JSON dans mon objet journee
}
func GetDataAPI(url string) []byte {
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

func GetNbrsDeMatch(jsonData *gabs.Container) float64 {
	return jsonData.Path("count").Data().(float64) //mettre float car si je met int, il n'accepte pas

}
func GetCurrentMatchday() string {
	url := "http://api.football-data.org/v1/competitions/445/leagueTable"
	jsonParsed2, _ := gabs.ParseJSON(GetDataAPI(url))

	return fmt.Sprint(jsonParsed2.Path("matchday").Data().(float64))
}
