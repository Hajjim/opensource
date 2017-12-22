package StibModule

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
	//	"time"
	//"github.com/Jeffail/gabs"
)

//variables avec les id de l'arrêt congrès direction Schaerbeek et Fort-Jaco
var CongresFor int = 6308
var CongresScha int = 6357

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
				minutsBeforeArrival := dateToMinuts(time.ExpectedArrivalTime)
				fmt.Println("Congrès direction Schaerbeek ligne " + strconv.Itoa(time.LineId) + " arrive dans " + strconv.Itoa(minutsBeforeArrival) + " minutes")
				if time.LineId == 92 {
					HorraireFor92 = append(HorraireFor92, strconv.Itoa(minutsBeforeArrival))
				}
				if time.LineId == 93 {
					HorraireFor93 = append(HorraireFor93, strconv.Itoa(minutsBeforeArrival))
				}
			}
		}
		if p.PointId == CongresFor {
			for _, time := range p.PassingTimes {
				minutsBeforeArrival := dateToMinuts(time.ExpectedArrivalTime)
				fmt.Println("Congrès direction Fort-Jaco ligne " + strconv.Itoa(time.LineId) + " arrive dans " + strconv.Itoa(minutsBeforeArrival) + " minutes")
				if time.LineId == 92 {
					HorraireScha92 = append(HorraireScha92, strconv.Itoa(minutsBeforeArrival))
				}
				if time.LineId == 93 {
					HorraireScha93 = append(HorraireScha93, strconv.Itoa(minutsBeforeArrival))
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
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	//création d'un client
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
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
