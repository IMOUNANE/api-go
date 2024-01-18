package main

import (
	"bufio"
	"sort"
	"strings"
	"text/tabwriter"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type formatedData = map[string]map[string]int

type index struct {
	Path string
	Version string
	Timestamp string
}

type sortableFormatedData struct {
	forge string
	modules int
	versions int
}

const (
	baseURL = "https://index.golang.org/index"
	moduleIndexName = "Modules"
	versionsIndexName = "Versions"
)

var (
	client *http.Client
	indexes []index
)


func groupIndexesByVersions(indexes []index) formatedData {
	groupedIndexes := make(formatedData)
	for _, ix := range indexes {
		if repo, ok := groupedIndexes[ix.Path]; ok {
			repo["Versions"]++
		} else {
			groupedIndexes[ix.Path] = map[string]int{moduleIndexName: 1, versionsIndexName: 1}
		}
	}

	return groupedIndexes
}

func groupFormatedIndexesByVersions(data formatedData) formatedData {
	groupedData := make(formatedData)

	for repo, meta := range data {
		forgeFormated := strings.Split(repo, "/")[0]

		if forge, ok := groupedData[forgeFormated]; ok {
			forge[moduleIndexName] += meta[moduleIndexName]
			forge[versionsIndexName] += meta[versionsIndexName]
		} else {
			groupedData[forgeFormated] = meta
		}
	}

	return groupedData
}

func sortFormatedData(data formatedData) []sortableFormatedData {
	var slice []sortableFormatedData

	for f, meta := range data {
		slice = append(slice, sortableFormatedData{
			forge:    f,
			modules:  meta[moduleIndexName],
			versions: meta[versionsIndexName],
		})
	}

	sort.Slice(slice, func(i, j int) bool {
		return slice[i].versions > slice[j].versions
	})

	return slice
}

func render(data formatedData) {
	sortedData := sortFormatedData(data)
	
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 20, 4, 0, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w,"Forge\t Modules\t Versions\n")
	totalModules := 0
	totalVersions := 0

	for _, d := range sortedData {
		fmt.Fprintf(w,"%s\t %d\t %d\n", d.forge, d.modules, d.versions)
		totalModules += d.modules
		totalVersions += d.versions
	}

	fmt.Fprintf(w, "_total_\t %d\t %d\t", totalModules, totalVersions)
}

func getIndexGolang() (*http.Response, error) {
	client = &http.Client{}

	request, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Disable-Module-Fetch", "true")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func main() {
	response, err := getIndexGolang()
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode < 400 {
		var indexes []index

		scanner := bufio.NewScanner(response.Body)
		for scanner.Scan() {
			var indexData index
			if err := json.Unmarshal(scanner.Bytes(), &indexData); err != nil {
				log.Printf(err.Error())
			}
			indexes = append(indexes, indexData)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		groupedByIndex := groupIndexesByVersions(indexes)
		groupedByVersion := groupFormatedIndexesByVersions(groupedByIndex)

		render(groupedByVersion)
	} else {
		log.Printf("Erreur de statut de la réponse : %d", response.StatusCode)
	}
}