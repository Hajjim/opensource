package main

//!!! pour importer nos packages persos, ils doivent se trouver dans le gopath ou le goroot.
//il faut donc faire une copie des packages dans /Www/src et les mettre dans le src du gopath
import (
	cedPack "cedricPackage"
	hajjiPack "hajjiPackage"
	"net/http"
	"time"
)

func main() {
	//lancement du module stib dans un thread à part
	go stibModule()
	hajjiPack.StartModule()
	//lancement du serveur web
	http.Handle("/", http.FileServer(http.Dir("./HTTP")))
	http.ListenAndServe(":8000", nil)
}

func stibModule() {
	//renouvelle les valeurs des variables toutes les 20 secondes
	for range time.Tick(time.Second * 21) {
		cedPack.GetVariablesFromServer()
	}
}w
