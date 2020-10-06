package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const downloadDomain = "https://songscribe.himadri.eu/download"

type versionFormat struct {
	DownloadFile string `json:"downloadFile"`
	BuildVersion string `json:"buildVersion"`
}

type indexReturn struct {
	CurrentVersion string `json:"currentVersion"`
	DownloadUrl    string `json:"downloadUrl"`
}

var platforms = []string{"mac", "windows"}
var platformJsonMap = map[string]string{}

func main() {
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := initHandler(); err != nil {
		log.Fatal(err)
		return
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func initHandler() error {
	for _, platform := range platforms {
		indexJson, err := createIndexJson(platform)
		if err != nil {
			return err
		}
		platformJsonMap[platform] = *indexJson
	}
	return nil
}

func createIndexJson(platform string) (*string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/version-%s.json", downloadDomain, platform))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("HTTP request failed: %d", resp.StatusCode))
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var versionFile versionFormat
	if err := json.Unmarshal(bodyBytes, &versionFile); err != nil {
		return nil, err
	}
	retValue := &indexReturn{
		CurrentVersion: versionFile.BuildVersion,
		DownloadUrl:    fmt.Sprintf("%s/%s", downloadDomain, versionFile.DownloadFile),
	}
	jsonBytes, err := json.Marshal(retValue)
	if err != nil {
		return nil, err
	}
	indexJson := string(jsonBytes)
	return &indexJson, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	platform := r.URL.Query().Get("platform")
	indexJson := platformJsonMap[platform]
	if indexJson == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, indexJson)
}
