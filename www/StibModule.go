package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"strconv"
//	"time"

	//"github.com/Jeffail/gabs"
)
//variables avec les id de l'arrêt congrès direction Schaerbeek et Fort-Jaco
var CongresFor int = 6308
var CongresScha int = 6357

var HorraireFor [4]string
var HorraireScha [4]string

//structure du json de la stib qu'on remplit avec le Unmarshal plus bas
type StibData struct {
	Points []struct {
		PointId int
		PassingTimes []struct {
			ExpectedArrivalTime string
			LineId int
		}
	}
}

func main() {
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

	//foreach Point p in stibData.Points
	for _, p := range s.Points {
		//si c'est bien le pointID qui nous intéresse
		if (p.PointId == CongresScha){
			//fmt.Println("Congrès direction Scha")
			for _, time := range p.PassingTimes{
				fmt.Println("Congrès direction Schaerbeek ligne "+ strconv.Itoa(time.LineId) +" : expected time arrival =" + time.ExpectedArrivalTime)
			}
		}
		if (p.PointId == CongresFor){
			for _, time := range p.PassingTimes{
				fmt.Println("Congrès direction Fort-Jaco ligne "+ strconv.Itoa(time.LineId) +" : expected time arrival =" + time.ExpectedArrivalTime)
			}
		}
		
	}
}

//func servant à récupérer le json de la stib sous forme de byte
func getDataAPI(url string) []byte {
	//création d'un client
	client := &http.Client{}
	//création d'une requête à envoyer au serveur
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	//ajout des headers dans la requête pour authentifier la connexion à l'open data de la stib
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer 01dc5ca7d2c53c40771fcce562bb0377")
	//création d'une réponse su base de la réception de la requête
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//sauvegarde de son contenu sous forme de byte
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//libération des ressources
	resp.Body.Close()
	return bytes
}
