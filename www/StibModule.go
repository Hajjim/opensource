package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
//	"time"

	//"github.com/Jeffail/gabs"
)

/*type stibData struct {
	Points []struct {
		PointID int 'json:"pointID"'
		PassingTimes []struct {
			ExpectedArrivalTime string 'json:"expectedArrivalTime"'
			LineID int 'json:"lineId"'
		} 'json:"passingTimes"'
	} 'json:"points"'
}*/

type stibData struct {
	Points []Point
}
type Point struct{
	PointID int
	PassingTimes []PassingTime
}
type PassingTime struct {
	ExpectedArrivalTime string
	LineID int
}

func main() {
	//l'url de la stib pour la ligne 6437
	url := "https://opendata-api.stib-mivb.be/OperationMonitoring/1.0/PassingTimeByPoint/6437"
	//récupération du json sous forme de string
	json := getDataAPI(url)
	//affichage du json
	fmt.Println(json)
}



func getDataAPI(url string) string {
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
	//passage sous forme de string
	string_body := string(bytes)
	//désallocation des ressources
	resp.Body.Close()
	return string_body
}
