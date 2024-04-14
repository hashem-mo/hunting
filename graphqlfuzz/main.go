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
	"time"
)
type AllValid struct {
    Data struct{} `json:"data"`
}
var LOCK sync.Mutex
var FOUNDED []string
var BLOCKED string

// Defining the base of the request body
var BodyData string
var outputFile *os.File
var saveOutput = false


var notfound = regexp.MustCompile(`^(?i)Cannot query field ['"][a-zA-Z0-9!\[\]_]+['"] on type ['"][a-zA-Z0-9!\[\]_]+['"]\.$`)
var othercase = regexp.MustCompile(`(?i)Cannot query field ['"][a-zA-Z0-9!\[\]_]+['"] on type ['"][a-zA-Z0-9!\[\]_]+['"]\. Did you mean (['"][a-zA-Z0-9!\[\]_]+['"])(, ['"][a-zA-Z0-9!\[\]_]+['"])*(, or ['"][a-zA-Z0-9!\[\]_]+['"])*`)
var gettext = regexp.MustCompile(`['"]([a-zA-Z0-9!\[\]_]+)['"]`)
var missingfields = regexp.MustCompile(`Field ['"][a-zA-Z0-9!\[\]_]+['"] of type ['"][a-zA-Z!\[\]]+['"] must have a selection of subfields. Did you mean ['"][a-zA-Z0-9!\[\]_]+ \{ \.\.\. \}['"]?`)
var missingargs = regexp.MustCompile(`Field ['"][a-zA-Z0-9!\[\]_]+['"] argument ['"][a-zA-Z0-9!\[\]_]+['"] of type ['"][a-zA-Z!\[\]]+['"] is required,? but (it was )?not provided.?`)



func PrintBanner(){
	fmt.Printf("   _____                 _      ____  _ ______                      \n")
	fmt.Printf("  / ____|               | |    / __ \\| |  ____|                     \n")
	fmt.Printf(" | |  __ _ __ __ _ _ __ | |__ | |  | | | |__ _   _ ___________ _ __ \n")
	fmt.Printf(" | | |_ | '__/ _` | '_ \\| '_ \\| |  | | |  __| | | |_  /_  / _ \\ '__|\n")
	fmt.Printf(" | |__| | | | (_| | |_) | | | | |__| | | |  | |_| |/ / / /  __/ |   \n")
	fmt.Printf("  \\_____|_|  \\__,_| .__/|_| |_|\\___\\_\\_|_|   \\__,_/___/___\\___|_|   \n")
	fmt.Printf("                  | |                                                \n")
	fmt.Printf("                  |_|\n")

	fmt.Println("\t\t\tTwitter: hashem_mo0")
}


// A function that sends a GET request to a URL and returns the status code and the response body
func sendRequest(URL string, method string, headers map[string]string, reqBody string, word string, proxy string, words []string) {
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
		log.Println(0, err.Error())
		return 
	}

// Setting the request headers
	for key, value := range headers {
		// Set each header to the request
		req.Header.Set(key, value)
	  }
	req.Header.Set("Accept-Encoding", "gzip")
// Closing the request body
	defer req.Body.Close()
// Sending the request
	resp, err := client.Do(req)
	if err != nil{
		log.Println(err)
		
	}
	// Checking if we are blocked and deciding what to do on user input

	if resp.StatusCode == 403 || resp.StatusCode == 429{
		
		fmt.Printf("Response Status Code: %v \n", resp.StatusCode)
		var input string
		if BLOCKED == "c" {
			return
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
		return
		}
	}
// Closing the response body
	defer resp.Body.Close()
	contentTypeSlice := resp.Header["Content-Type"]
	if len(contentTypeSlice) == 0 {
		log.Fatalln(resp.StatusCode, "NO CONTENT-TYPE HEADER")
	}

	if !strings.Contains(contentTypeSlice[0], "application/json"){
		log.Fatalf("Content type is not JSON, Status Code: %d, Body:\n %s \n\n", resp.StatusCode, reqBody )
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
		log.Println( resp.StatusCode, err.Error())
	}
	parseResp(string(body), words)

	
}






func handleOutput(word string)(){

	if saveOutput{
		outputFile.WriteString(fmt.Sprintf("%s\n", word))	
	}
	fmt.Println(word)

}


func applyRegex(message string)([]string){

	switch {
			case notfound.MatchString(message):
				return nil
			case othercase.MatchString(message):
				var words []string
				for _, w := range othercase.FindAllStringSubmatch(message, -1)[0][1:]{
					if strings.TrimSpace(w) == "" {
						continue
					}
					words = append(words, gettext.FindAllStringSubmatch(w, -1)[0][1])
					}
				return words
			case missingfields.MatchString(message):
				var words []string
				word := missingfields.FindAllStringSubmatch(message, -1)[0][0]
				words = append(words, gettext.FindAllStringSubmatch(word, -1)[0][1])
				return	words

			case missingargs.MatchString(message):
				var words []string
				word := missingargs.FindAllStringSubmatch(message, -1)[0][0]
				
				words = append(words, gettext.FindStringSubmatch(word)[1])
				return words
			default:
				log.Printf("Unhandled Error: %s", message)
				return nil
			}
}



func parseResp(body string, words []string){
// 	Defining a data interface to store the parsed response in
	var  data interface{}
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println(err)
	}

	switch data:= data.(type) {
	case AllValid:
		for _, word := range(words){
			LOCK.Lock()
			 if !DoesSliceContains(FOUNDED, word){
			 FOUNDED = append(FOUNDED, word)
			 handleOutput(word)
		 }
			 LOCK.Unlock()
		}

	case map[string]interface{}:

		var valids []string
		
	 	for  _, v := range data["errors"].([]interface{}){
 			errMsg := v.(map[string]interface{})["message"].(string)
			queries := applyRegex(errMsg)
			valids = append(valids, queries...)
		}
		LOCK.Lock()
		for _, q := range valids {
		if !DoesSliceContains(FOUNDED, q){
		FOUNDED = append(FOUNDED, q)
		handleOutput(q)
	}}
		LOCK.Unlock()


	default:
	 	log.Printf("Unhandled Error, Query Used: %s, Data Returned: %s", words, data)
	}

}



func process(url string, method string, headers map[string]string,  wordlist []string, workers int, qpr int, proxy string, delay int){

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

			var words [] string
			defer wg.Done()
			for word := range wordlistChan { // receive one element from the channel
				
				if len(words) < qpr{
					words = append(words, word)
					continue
				}
				body := strings.Replace(BodyData, "FUZZ", strings.Join(words, " "), 1)	
				sendRequest(url, method, headers, body, word, proxy, words) // pass the element to the function
				time.Sleep(time.Duration(delay) * time.Millisecond)
				words = nil
			}
			
		}()
	} 
	wg.Wait() // wait for all workers to finish
}



func main() {
	PrintBanner()
	// Define the base URL and the wordlist
	baseUrl, method, headers, threads, output, wordlistFile, qpr, proxy, delay := ParseArgs()
	
	if output != ""{
		saveOutput = true
		_, err := os.Stat(output)
		if err == nil {
		outputFile, err = os.OpenFile(output, os.O_WRONLY|os.O_APPEND, 0666)
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

	process(baseUrl, method, headers, wordlist, threads, qpr, proxy, delay)

}
