/*The MIT License (MIT)

Copyright (c) 2018 cedricgoo

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

package StibModule

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	//	"time"
	//"github.com/Jeffail/gabs"
)

//variables avec les id de l'arrêt congrès direction Schaerbeek et Fort-Jaco
var CongresFor int = 6308
var CongresScha int = 6357

//slices qui récupèrent les minutes avant l'arrivée des trams
var HorraireFor92 []string
var HorraireFor93 []string

var HorraireScha92 []string
var HorraireScha93 []string

//structure du json de la stib qu'on remplit avec le Unmarshal plus bas
type StibData struct {
	Points []struct {
		PointId      int
		PassingTimes []struct {
			ExpectedArrivalTime string
			LineId              int
		}
	}
}

//fonction lancée pour recevoir les valeurs du serveur et les intégrer dans le html
func GetVariablesFromServer() {

	//l'url de la stib pour la ligne ArretId (conversion int to string avec import strconv)
	url := "https://opendata-api.stib-mivb.be/OperationMonitoring/1.0/PassingTimeByPoint/" + strconv.Itoa(CongresScha) + "%2C" + strconv.Itoa(CongresFor)
	//récupération du json sous forme de byte
	json_resp := getDataAPI(url)

	//pour afficher la totalité sous forme de string :
	//string_body := string(json_resp)
	//fmt.Println(string_body)

	//définition de la variable qui va receuillir les données
	var s StibData
	//envoi des données dans la structure
	json.Unmarshal(json_resp, &s)
	//print de test d'affichage du point ID du premier point
	//fmt.Println(s.Points[0].PointId)

	//clear les variables à l'itération suivante pour les remplacer par les nouvelles valeurs
	HorraireFor92 = HorraireFor92[:0]
	HorraireFor93 = HorraireFor93[:0]
	HorraireScha92 = HorraireScha92[:0]
	HorraireScha93 = HorraireScha93[:0]

	//foreach Point p in stibData.Points
	for _, p := range s.Points {
		//si c'est bien le pointID qui nous intéresse
		if p.PointId == CongresScha {
			//fmt.Println("Congrès direction Scha")
			for _, time := range p.PassingTimes {
				//récupération des minutes avant d'arriver
				minutsBeforeArrival := dateToMinuts(time.ExpectedArrivalTime)
				fmt.Println("Congrès direction Schaerbeek ligne " + strconv.Itoa(time.LineId) + " arrive dans " + strconv.Itoa(minutsBeforeArrival) + " minutes")
				//envoi de ces minutes dans le slice correspondant à la ligne
				if time.LineId == 92 {
					HorraireScha92 = append(HorraireScha92, strconv.Itoa(minutsBeforeArrival))
				}
				if time.LineId == 93 {
					HorraireScha93 = append(HorraireScha93, strconv.Itoa(minutsBeforeArrival))
				}
			}
		}
		if p.PointId == CongresFor {
			for _, time := range p.PassingTimes {
				minutsBeforeArrival := dateToMinuts(time.ExpectedArrivalTime)
				fmt.Println("Congrès direction Fort-Jaco ligne " + strconv.Itoa(time.LineId) + " arrive dans " + strconv.Itoa(minutsBeforeArrival) + " minutes")
				if time.LineId == 92 {
					HorraireFor92 = append(HorraireFor92, strconv.Itoa(minutsBeforeArrival))
				}
				if time.LineId == 93 {
					HorraireFor93 = append(HorraireFor93, strconv.Itoa(minutsBeforeArrival))
				}
			}
		}

	}
	//fmt.Println(HorraireFor92[0], " et ensuite ", HorraireFor92[1])
	//fmt.Println(HorraireFor93[0], " et ensuite ", HorraireFor93[1])
	//fmt.Println(HorraireScha92[0], " et ensuite ", HorraireScha92[1])
	//fmt.Println(HorraireScha93[0], " et ensuite ", HorraireScha93[1])
}

//func servant à récupérer le json de la stib sous forme de byte
func getDataAPI(url string) []byte {

	//Tansporteur
	netTransport := &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
	}

	//création d'un client
	client := &http.Client{
		//erreur récurente avec les Timeouts des clients sur cette version de Go
		//si on le met à 0 il est annulé et tant qu'il n'y a pas de problème, l'erreur
		//est évitée... à améliorer
		Timeout:   time.Second * 0,
		Transport: netTransport,
	}
	//création d'une requête à envoyer au serveur
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("fatal à la création de la requete ", err)
	}
	//ajout des headers dans la requête pour authentifier la connexion à l'open data de la stib
	req.Header.Set("Accept", "application/json")
	//req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer 01dc5ca7d2c53c40771fcce562bb0377")
	//création d'une réponse su base de la réception de la requête
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(" fatal à la réception ", err)

	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("fatal dans le io ", err)
	}
	//vérification que toutes les données sont bien récupérées avant de libérer les ressources
	io.Copy(ioutil.Discard, resp.Body)
	//libération des ressources
	resp.Body.Close()
	//sauvegarde de son contenu sous forme de byte

	return bytes
}

//fonction qui prend en entrée le format de date en string de la stib et
//la compare à la date actuelle pour retourner un nombre de minutes avant arrivée
func dateToMinuts(arrivalTime string) int {
	//valeur en minute qui sera retournée
	minut := 0
	//récupère l'heure actuelle en string
	actualTimeInStr := time.Now().Local()
	//fmt.Println("Heure actuelle : ", actualTimeInStr)

	//récupère l'heure d'arrivée en string
	arrivalTimeInStr := arrivalTime
	//fmt.Println("Heure d'arrivée supposée : ", arrivalTimeInStr)

	//format de la date reçuepar la stib
	layOut := "2006-01-02T15:04:05"
	//convertion du string en type Time pour récupérer les minutes facilement
	timeStamp, err := time.Parse(layOut, arrivalTimeInStr)
	if err != nil {
		fmt.Println(err)

	}
	//si la minute actuelle est plus grande que la minute d'arrivée (heure passée)
	// alors on ajoute 60 à la minute d'arrivée lors de la soustraction pour ajouter l'heure perdue
	if actualTimeInStr.Minute() > timeStamp.Minute() {
		minut = timeStamp.Minute() - actualTimeInStr.Minute() + 60
	} else {
		minut = timeStamp.Minute() - actualTimeInStr.Minute()
	}

	//fmt.Println("Le tram arrive dans ", minut, " minutes")
	return minut
}
