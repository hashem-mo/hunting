package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/akamensky/argparse"
)

var LOCK sync.Mutex
var FOUNDED []string
var BLOCKED string

// Defining the base of the request body
var BodyData string
var outputFile *os.File
var saveOutput = false


var notfound = regexp.MustCompile(`^(?i)Cannot query field "\w+" on type "\w+"\.$`)
var othercase = regexp.MustCompile(`(?i)Cannot query field "\w+" on type "\w+"\. Did you mean ("\w+")(, "\w+")*(, or "\w+")*`)
var gettext = regexp.MustCompile(`\"(\w+)\"`)



func PrintBanner(){
	fmt.Printf("   _____                 _      ____  _ ______                      \n")
	fmt.Printf("  / ____|               | |    / __ \\| |  ____|                     \n")
	fmt.Printf(" | |  __ _ __ __ _ _ __ | |__ | |  | | | |__ _   _ ___________ _ __ \n")
	fmt.Printf(" | | |_ | '__/ _` | '_ \\| '_ \\| |  | | |  __| | | |_  /_  / _ \\ '__|\n")
	fmt.Printf(" | |__| | | | (_| | |_) | | | | |__| | | |  | |_| |/ / / /  __/ |   \n")
	fmt.Printf("  \\_____|_|  \\__,_| .__/|_| |_|\\___\\_\\_|_|   \\__,_/___/___\\___|_|   \n")
	fmt.Printf("                  | |                                                \n")
	fmt.Printf("                  |_|\n")

	fmt.Println("\n By: hashem_mo")
}


// A function that sends a GET request to a URL and returns the status code and the response body
func sendRequest(URL string, method string, headers map[string]string, reqBody string, word string, proxy string, words []string) (int, string) {
	// Making a proxy string, this is the way of definging a proxy in GoLang
	var client http.Client
	if proxy != ""{
	// Parsing the proxy url
	proxyUrl, _ := url.Parse(proxy)
	tr := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	client = http.Client{Transport: tr}

	}

// Creating the request NOTE: URL is upper case to avoid duplication with url.Parse()
	req, err := http.NewRequest(method, URL, strings.NewReader(reqBody))

	if err != nil {
		return 0, err.Error()
	}
// Setting the request headers
	for key, value := range headers {
		// Set each header to the request
		req.Header.Set(key, value)
	  }
// Closing the request body
	defer req.Body.Close()
// Sending the request
	resp, err := client.Do(req)
	if err != nil{
		fmt.Println(err)
		return 1, "error occured"
	}
	// Checking if we are blocked and deciding what to do on user input
	if resp.StatusCode == 403 {
		fmt.Printf("Response Status Code: %v \n", resp.StatusCode)
		var input string
		if BLOCKED == "c" {
			return resp.StatusCode, "Forbidden"
		}else{
		// Locking thread until the user decides what to do
		LOCK.Lock()
		fmt.Print("Seems like we have been blocked: c to continue, q to quit or u to update Cookies: ")
		fmt.Scanln(&input)
		switch {
		case input == "q": 
			fmt.Println("Exiting....")
			os.Exit(1)
		case input == "c":
			BLOCKED = "c"
		case input == "u":
			fmt.Print("Paste your cookies here: ")
			var cookie string
			fmt.Scanln(&cookie)
			headers["Cookie"] = cookie
		}
		LOCK.Unlock()
		}
	}
// Closing the response body
	defer resp.Body.Close()
	contentTypeSlice := resp.Header["Content-Type"]
	if len(contentTypeSlice) == 0 {
		return resp.StatusCode, "NO CONTENT-TYPE HEADER"
	}

	if !strings.Contains(contentTypeSlice[0], "application/json"){
		return resp.StatusCode, "Content type is not JSON"
	}
// Reading the response body
	var body []byte
   // Check if the response is encoded in gzip format
   if resp.Header.Get("Content-Encoding") == "gzip" {
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	// Read the decompressed response body
	body, err = io.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
	}
} else {
	// The response is not gzip encoded, so read it directly
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

}


	if err != nil {
		return resp.StatusCode, err.Error()
	}
	parseResp(string(body), words)

	return resp.StatusCode, string(body)
}


func DoesSliceContains(mySlice []string, myStr string) bool { // This func is specially created to check if the slice has an item! Cuz Golang doesn't have a function to do this job!
	for _, Value := range mySlice {
		if Value == myStr {
			return true
		}
	}
	return false
}



func handleOutput(word string)(){

	if saveOutput{
		outputFile.WriteString(fmt.Sprintf("%s\n", word))	
	}
	fmt.Println(word)

}


func applyRegex(errors []interface{}, word string)(){

	if len(errors) > 1{
		fmt.Println(word)
	}
	message := errors[0].(map[string]interface{})["message"].(string)

	switch {
			case notfound.MatchString(message):
				break
			case othercase.MatchString(message):
				for _, w := range othercase.FindAllStringSubmatch(message, -1)[0][1:]{
					if strings.TrimSpace(w) == "" {
						continue
					}
					LOCK.Lock()
					str := gettext.FindAllStringSubmatch(w, -1)[0][1]
			 		if !DoesSliceContains(FOUNDED, str){
			 		FOUNDED = append(FOUNDED, str)
			 		handleOutput(str)
			 	}
			 		LOCK.Unlock()
					}

			default:
				fmt.Printf("Unhandled Error1: %s, Query used %s\n", message, word)
			}
}


			// 	for _, m := range othercase.FindAllStringSubmatch(message, -1)[0][2:] {
			// 		fmt.Println(othercase.FindAllStringSubmatch(message, -1)[0][2:])
			// 		if len(gettext.FindAllStringSubmatch(m, -1)) <= 1{
			// 			continue
			// 		}
			// 		if len(gettext.FindAllStringSubmatch(m, -1)[0]) > 2 {
			// 			continue
			// 		}
			// 		LOCK.Lock()
			// 		if !DoesSliceContains(FOUNDED, gettext.FindAllStringSubmatch(m, -1)[0][1]){
			// 		FOUNDED = append(FOUNDED, gettext.FindAllStringSubmatch(m, -1)[0][1])
			// 		handleOutput(gettext.FindAllStringSubmatch(m, -1)[0][1])
			// 	}
			// 		LOCK.Unlock()
			// }


func parseResp(body string, words []string){
// 	Defining a data interface to store the parsed response in
	var  data interface{}
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println(err)
	}

	switch data := data.(type) {
		// Parsing the response if it's an array
	case []interface{}:
		// looping over each item of the array
		for i, v := range data{
			// checking if the 'errors' object exsits
			if arr, ok := v.(map[string]interface{})["errors"].([]interface{}); ok {
				applyRegex(arr, words[i]) 
			}
			if _, ok := v.(map[string]interface{})["data"].(interface{}); ok {
					handleOutput(words[i])
			}
		}
	case interface{}:
		if _, ok := data.(map[string]interface{})["data"].(interface{}); ok{
			handleOutput(words[0])
		}else{
		var temp []interface{}
		applyRegex(append(temp, data), words[0])
	}
	default:
		fmt.Printf("Unhandled Error, Query Used: %s, Data Returned: %s", words, data)
	}

}



func process(url string, method string, headers map[string]string,  wordlist []string, workers int, qpr int, proxy string){

	wordlistChan := make(chan string, len(wordlist)) // create a channel of strings

	// Create a sync.WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Start a goroutine to send to the channel

		for _, word := range wordlist {
			wordlistChan <- word // send one element to the channel
		}
		close(wordlistChan) // close the channel after sending all elements
	

	// Start the workers
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			var quries  []string
			var words [] string
			defer wg.Done()
			for word := range wordlistChan { // receive one element from the channel
				var body = strings.Replace(BodyData, "FUZZ", word, 1)
				if len(quries) < qpr{
					quries = append(quries, body)
					words = append(words, word)
					continue
				}
				body = fmt.Sprintf("[%s]",strings.Join(quries, ","))
				
				sendRequest(url, method, headers, body, word, proxy, words) // pass the element to the function
				quries = nil
				words = nil
			}
			
		}()
	} 
	wg.Wait() // wait for all workers to finish
}



func parseReq(fileName string)(string, string, map[string]string){

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	defer file.Close()
	headers := map[string]string{}
	var method string;
	var path string;

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	firstLine := scanner.Text()

	method = strings.Split(firstLine, " ")[0]
	path = strings.Split(firstLine, " ")[1]

	var bodyline = false 
	var reqbody strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if bodyline{
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




func parseArgs()(string, string, map[string]string, int, string, string, int, string){
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
	var reqData *string = parser.String("d", "data", &argparse.Options{
		Required: false,
		Help: "JSON object contains query object and FUZZ in the place you would like to fuzz",
		Default: "",
	})
	
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
		method, baseUrl, headers = parseReq(*req)
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


	return baseUrl, method, headers, *threads, *output, *wordlist, *qpr, *proxy


}



func main() {
	PrintBanner()
	// Define the base URL and the wordlist
	baseUrl, method, headers, threads, output, wordlistFile, qpr, proxy := parseArgs()
	
	if output != ""{
		saveOutput = true
		_, err := os.Stat(output)
		if err == nil {
		outputFile, err = os.OpenFile(wordlistFile, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
		  fmt.Print(err)
		  os.Exit(0)
		}
	}else{
		outputFile, err = os.Create(output)
		if err != nil{
			fmt.Println(err)
			os.Exit(1)
		}
	} 
		defer outputFile.Close()
	}

	file, err := os.Open(wordlistFile)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer file.Close()

	// create a scanner
	scanner := bufio.NewScanner(file)

	// create a slice of strings
	wordlist := make([]string, 0)

	// loop over the file and scan each line
	for scanner.Scan() {
		// get the current line as a string
		line := scanner.Text()
		// append the line to the slice
		wordlist = append(wordlist, line)
	}

	// check for any error during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	process(baseUrl, method, headers, wordlist, threads, qpr, proxy)

}
