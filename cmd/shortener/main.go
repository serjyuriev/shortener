package main

import (
	"log"

	"github.com/serjyuriev/shortener/internal/pkg/server"
)

var (
	// win-flags:
	buildVersion string = "N/A" // -X main.buildVersion=v1.0.0
	buildDate    string = "N/A" // -X 'main.buildDate=$(Get-Date)'
	buildCommit  string = "N/A" // -X 'main.buildCommit=$((git show --oneline -s).split(" ")[0])'
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	log.Printf("Build version: %s\n", getBuildFlag(buildVersion))
	log.Printf("Build date: %s\n", getBuildFlag(buildDate))
	log.Printf("Build commit: %s\n", getBuildFlag(buildCommit))

	server, err := server.NewServer()
	if err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
	log.Fatal(server.Start())
}

func getBuildFlag(buildFlag string) string {
	if buildFlag == "" {
		buildFlag = "N/A"
	}
	return buildFlag
}
