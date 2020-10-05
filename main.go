package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const downloadUrl = "https://songscribe.himadri.eu/"

type indexReturn struct {
	CurrentVersion string `json:"current_version"`
	DownloadUrl    string `json:"download_url"`
}

var indexJson string

func main() {
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := initHandler()
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func initHandler() error {
	resp, err := http.Get(downloadUrl + "download/version")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("HTTP request failed: %d", resp.StatusCode))
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	retValue := &indexReturn{
		CurrentVersion: strings.TrimSpace(string(bodyBytes)),
		DownloadUrl:    downloadUrl,
	}
	jsonBytes, err := json.Marshal(retValue)
	if err != nil {
		return err
	}
	indexJson = string(jsonBytes)
	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, indexJson)
}
