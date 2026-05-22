package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bienvenue sur le forum !")
}

func main() {
	http.HandleFunc("/", homeHandler)

	fmt.Println("Serveur lancé sur http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erreur serveur :", err)
	}
}