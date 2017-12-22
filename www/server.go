package main

//!!! pour importer nos packages persos, ils doivent se trouver dans le gopath ou le goroot.
//il faut donc faire une copie des packages dans /Www/src et les mettre dans le src du gopath
import (
	cedPack "cedricPackage"
	hajjiPack "hajjiPackage"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./HTTP")))
	http.ListenAndServe(":8000", nil)
	cedPack.StartModule()
	hajjiPack.StartModule()
}
