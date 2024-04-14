package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/akamensky/argparse"
)

func ParseReq(fileName string) (string, string, map[string]string) {

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	defer file.Close()
	headers := map[string]string{}
	var method string
	var path string

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	firstLine := scanner.Text()

	method = strings.Split(firstLine, " ")[0]
	path = strings.Split(firstLine, " ")[1]

	var bodyline = false
	var reqbody strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if bodyline {
			reqbody.WriteString(strings.TrimSpace(line))
			continue
		}

		if strings.EqualFold(line, "") {
			bodyline = true
			continue
		}

		parts := strings.SplitN(line, ":", 2)

		headers[parts[0]] = strings.TrimSpace(parts[1])
	}
	baseUrl := fmt.Sprintf("https://%s%s", headers["Host"], path)
	
	BodyData = reqbody.String()
	
	return method, baseUrl, headers
}






func DoesSliceContains(mySlice []string, myStr string) bool { // This func is specially created to check if the slice has an item! Cuz Golang doesn't have a function to do this job!
	for _, Value := range mySlice {
		if Value == myStr {
			return true
		}
	}
	return false
}




func ParseArgs()(string, string, map[string]string, int, string, string, int, string, int){
	var headers = make(map[string]string)
	var method = "POST"
	var baseUrl string

	parser := argparse.NewParser("GraphQlFuzzer", "A tool to help fuzzing GraphQl APIs")
	var req *string = parser.String("r", "request", &argparse.Options{
		Required: false,
		Help: "A path to a file containing raw request put FUZZ where you want to fuzz weather query, mutation etc...",
	Default: ""})

	var proxy *string = parser.String("p", "proxy", &argparse.Options{
		Required: false,
		Help: "A URL to a proxy to forward traffic to",
	Default: ""})

	var threads *int = parser.Int("t", "threads", &argparse.Options{
		Required: false,
		Help: "Number of threads to use. Default: 1",
	Default: 1})
// qpr: number of queries to send each request
	var qpr *int = parser.Int("n", "nestedQueries", &argparse.Options{
		Required: false,
		Help: "Number of queries to send each request. Default: 5",
	Default: 5})

	var wordlist *string = parser.String("w", "wordlist", &argparse.Options{
		Required: true,
		Help: "A path to a wordlist to use for fuzzing",
	Default: ""})

	var url *string = parser.String("u", "url", &argparse.Options{
		Required: false,
		Help: "A URL to the GraphQl API",
	Default: nil})

	var headersList *[]string = parser.StringList("H", "header", &argparse.Options{
		Required: false,
		Help: "HTTP header to user 'Same sytax as CURL'",
	Default: nil})
	var output *string = parser.String("o", "output", &argparse.Options{
		Required: false,
		Help: "A path to a file to save the output to",
		Default: nil,})
	var reqData *string = parser.String("B", "body", &argparse.Options{
		Required: false,
		Help: "JSON object contains query object and FUZZ in the place you would like to fuzz",
		Default: "",
	})
	var delay *string = parser.String("d", "delay", &argparse.Options{
		Required: false,
		Help: "Time to delay between requests in milliseconds",
	Default: "0"})
	
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	if (*req == "" && *url == ""){
		fmt.Println("You must provide an input either a request or a URL")
		os.Exit(1)
	}
	if len(*reqData) != 0{
		BodyData = *reqData
	}
	if *url == "" {
		method, baseUrl, headers = ParseReq(*req)
	}else{
		baseUrl = *url
	}
	_, err1 := os.Stat(*wordlist)
	if err1 != nil {
		fmt.Println("Could not open the wordlist file")
		os.Exit(1)
	}

	if len(*headersList) != 0{
		for _, v := range *headersList {
			parts := strings.SplitN(v, ":", 2)
			headers[parts[0]] = strings.TrimSpace(parts[1])

		}
	}
	if (*reqData == "" && *req == ""){
		fmt.Println("You must provide a request body")
		os.Exit(1)
	}

	d, err := strconv.Atoi(*delay)
	if err != nil {
		log.Fatal(err)
	}

	return baseUrl, method, headers, *threads, *output, *wordlist, *qpr, *proxy, d


}


func SubtractSlice(sliceA, sliceB []string) []string {
    var diff []string

    // Create a map to store the elements of sliceB for quick lookup.
    lookup := make(map[string]struct{})
    for _, item := range sliceB {
        lookup[item] = struct{}{}
    }

    // Add only the elements from sliceA that are not found in the lookup map.
    for _, item := range sliceA {
        if _, found := lookup[item]; !found {
            diff = append(diff, item)
        }
    }

    return diff
}