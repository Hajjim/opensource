package main

import (  
    "fmt"
    "io/ioutil"
   "log"
   "net/http"
    "time"
    "github.com/Jeffail/gabs"
)

type Match struct{
    Status string
    Matchday float64
    Team1 string
    Team2 string
    Results [2]float64

}

func main() {
  dateNow := getDate()
 
  var matches []Match
    url := "http://api.football-data.org/v1/competitions/445/fixtures?timeFrameStart="+dateNow+"&timeFrameEnd="+dateNow

    jsonParsed, _ := gabs.ParseJSON( getContentFromAPI(url) )
   
   fmt.Println(getNbrsResults (jsonParsed))
   getData(jsonParsed,matches)

  
    
}

func getDate() string{
    now := time.Now()
    mask := "2006-01-02"
return now.Format(mask)
}

func getContentFromAPI(url string) []byte{
   
    resp, err := http.Get(url)
     if(err != nil){
         log.Fatal(err)
     }

     defer resp.Body.Close()
     bytes, err := ioutil.ReadAll(resp.Body)
     if(err != nil){
        log.Fatal(err)
    }
    return bytes
}

func getNbrsResults(jsonData *gabs.Container) float64{
    nbrResults:= jsonData.Path("count").Data().(float64) //mettre float car si je met int, il n'accepte pas
    return nbrResults
}


//Recupation des données des matches
    func getData(jsonData *gabs.Container, m []Match){
      
        children, _ := jsonData.S("fixtures").Children()
        for _, child := range children {
            /*m[i].Status = child.Path("status").Data().(string) Utiliser la fonction append() pour remplir le slice plutôt que faire m[i].x*/
            m = append(m,Match {child.Path("status").Data().(string), 
                child.Path("matchday").Data().(float64),child.Path("homeTeamName").Data().(string),
                child.Path("awayTeamName").Data().(string),
                [2]float64{child.Path("result.goalsHomeTeam").Data().(float64),child.Path("result.goalsAwayTeam").Data().(float64)}})//fin de la ligne
              
        }
        fmt.Println(m)
    }

    
    
