/*The MIT License (MIT)

Copyright (c) 2018 Dbelataris

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.*/

package Football

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/isib/ISIBOpen-source/www/src/BelatarisPackage/Jeffail/gabs"
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
	j.Matches = nil //Vide le tableau pour se débarraser des anciennes données des matches et pour ne garder que les données les plus récentes
	children, _ := jsonData.S("fixtures").Children()
	for _, child := range children {
		stringDate := ConvertToLocalTime(child.Path("date").Data().(string)) // Conversion de l'heure de l'API vers notre heure locale
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
				stringDate, imgUrl1, imgUrl2})

		case "SCHEDULED": //2ème Cas particulier car quand le match n'est pas encore joué les résultats sont à null et dans certains cas ces matchs sont réferencés en tant que SCHEDULED au lieu de TIMED
			j.Matches = append(j.Matches, Match{child.Path("status").Data().(string), ///Ajout des match dans le tableau
				child.Path("homeTeamName").Data().(string),
				child.Path("awayTeamName").Data().(string),
				0, 0,
				child.Path("matchday").Data().(float64),
				stringDate, imgUrl1, imgUrl2})
		default: //Matches terminés ou en cours
			j.Matches = append(j.Matches, Match{child.Path("status").Data().(string), child.Path("homeTeamName").Data().(string),
				child.Path("awayTeamName").Data().(string),
				child.Path("result.goalsHomeTeam").Data().(float64),
				child.Path("result.goalsAwayTeam").Data().(float64),
				child.Path("matchday").Data().(float64),
				stringDate, imgUrl1, imgUrl2})
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
func ConvertToLocalTime(timeToConvert string) string {
	mask := "2006-01-02T15:04:05Z"                 //Format de la date fournit par l'API
	timeDate, _ := time.Parse(mask, timeToConvert) //COnversion de la date String en objet time
	timeDate = timeDate.Local()                    //Conversion de l'heure reçu par l'API en l'heure local de Belgique
	ConvertedTime := timeDate.Format("02-01-2006T15:04:05")
	return ConvertedTime
}
