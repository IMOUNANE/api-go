package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
)

// Définir une structure pour représenter les données JSON
type PackageInfo struct {
	Path      string `json:"Path"`
	Version   string `json:"Version"`
	Timestamp string `json:"Timestamp"`
}

// Variable globale pour stocker les informations des packages
var globalPackages []PackageInfo

func countInstances() map[string]int {
	// Créer une carte pour stocker le nombre d'instances pour chaque valeur unique de "Path"
	countMap := make(map[string]int)

	// Parcourir les informations des packages
	for _, pkg := range globalPackages {
		// Incrémenter le compteur pour la valeur de "Path"
		countMap[pkg.Path]++
	}

	return countMap
}

func getIndexGolang() {
	// URL à laquelle faire la requête GET
	url := "https://index.golang.org/index"

	// Faire la requête HTTP GET
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la requête GET :", err)
		return
	}
	defer response.Body.Close()

	// Scanner pour lire le corps de la réponse ligne par ligne
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		// Chaque ligne contient un objet JSON
		line := scanner.Text()

		// Décodez le JSON dans la structure de données
		var pkg PackageInfo
		err := json.Unmarshal([]byte(line), &pkg)
		if err != nil {
			fmt.Println("Erreur lors du décodage JSON :", err)
			continue
		}

		// Ajouter les informations du package à la liste globale
		globalPackages = append(globalPackages, pkg)
	}

	// Vérifier les erreurs de numérisation
	if err := scanner.Err(); err != nil {
		fmt.Println("Erreur lors de la lecture du corps de la réponse :", err)
		return
	}

	// Afficher le nombre d'instances pour chaque valeur unique de "Path"
	countMap := countInstances()
	for path, count := range countMap {
		fmt.Printf("Nombre d'instances pour %s : %d\n", path, count)
	}
}

func main() {
	getIndexGolang()
}
