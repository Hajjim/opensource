package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
//	"time"

	//"github.com/Jeffail/gabs"
)

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
	//l'url de la stib pour la ligne 6437
	url := "https://opendata-api.stib-mivb.be/OperationMonitoring/1.0/PassingTimeByPoint/6437%2C6309"
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
	fmt.Println(s.Points[0].PointId)
	/*for _, Points := range s {
	
		for _, PassingTimes := range Points {
		fmt.Println(LineId)
		}
	}*/
}

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
	//passage sous forme de string
	//string_body := string(bytes)
	//désallocation des ressources
	resp.Body.Close()
	return bytes
}
