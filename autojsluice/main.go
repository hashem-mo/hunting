package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/BishopFox/jsluice"
	"github.com/akamensky/argparse"
)





func parseArgs()(string){

	parser := argparse.NewParser("autoJsluice", "A tool to help using Jsluice")
	var list *string = parser.String("l", "list", &argparse.Options{
		Required: false,
		Help: "A file contains URLs for js files",
	Default: ""})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	return *list
	}



func main() {
	var input = parseArgs()
	fmt.Println(input)
	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// create a slice of strings


	// loop over the file and scan each line
	for scanner.Scan() {
		line := scanner.Text()
		
		req, err := http.Get(strings.TrimSpace(line))

		if err != nil {
			fmt.Println(err)
			continue 
		}
		if req.StatusCode != 200{
			fmt.Printf("%s      %s", line, req.StatusCode)
			continue
		}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
		}
		analyzer := jsluice.NewAnalyzer(body)
		for _, match := range analyzer.GetURLs() {
			fmt.Println(match.URL)
		}


	}




	}