package main

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	url       string
	urlList   string
	threads   int
	output    string
	results   []map[string]interface{}
	mu        sync.Mutex
	wg        sync.WaitGroup
	regex     = `(?m)((\"|\'|\x60)([\\n\s])*FUZZ(\s|\n)+[=\[\]\@\#\d:\s\w\$\("\)!,\_\-\;\.\{\}\\/]+)`
	components = []string{"query", "mutation", "fragment", "subscription",
		"scalar", "enum", "interface", "union", "directive"}
	hashes []string
)

type headerSlice []string

func (h *headerSlice) String() string {
	return fmt.Sprint(*h)
}

func (h *headerSlice) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func init() {
	flag.StringVar(&url, "u", "", "A single URL to extract GraphQl queries from")
	flag.StringVar(&urlList, "l", "", "A list of URL to extract graphql queries")
	flag.IntVar(&threads, "t", 3, "Number of threads to use \"Works only with multiple URL\"")
	flag.StringVar(&output, "o", "", "A JSON file to save the output to")
}

func main() {
	var headers headerSlice
	flag.Var(&headers, "H", "A header to include in the request")
	flag.Parse()

	if url != "" {
		ext(url, headers)
	} else if urlList != "" {
		urls, err := readLines(urlList)
		if err != nil {
			fmt.Println("Error reading URL list:", err)
			return
		}

		ch := make(chan string, threads)
		for i := 0; i < threads; i++ {
			wg.Add(1)
			go worker(ch, headers)
		}

		for _, u := range urls {
			ch <- u
		}

		close(ch)
		wg.Wait()
	}

	if output != "" {
		saveResults(output)
	}
}

func ext(url string, headers []string) {
	mu.Lock()
	defer mu.Unlock()

	allQueries := make(map[string]interface{})
	allQueries["url"] = url
	numQueries := 0

	req, err := http.NewRequest("GET", strings.Trim(url, " \n\r"), nil)
	if err != nil {
		fmt.Printf("Error creating request for %s: %v\n", url, err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		req.Header.Set(parts[0], parts[1])
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url , err)
		return
	}
	defer resp.Body.Close()


	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error while trying to fetch %s, status code: %d\n", url, resp.StatusCode)
		return
	}


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body for %s: %v\n", url, err)
		return
	}


	for _, component := range components {

		re := regexp.MustCompile(strings.ReplaceAll(regex, "FUZZ", component))
		findings := re.FindAllStringSubmatch(string(body), -1)
		
		queries := make([]string, 0)
		for _, match := range findings {
			
			hash := hashString(cleanQuery(match[1]))
			if ifExists(hashes, hash) {
				continue
			}
			hashes = append(hashes, hash)
			
			queries = append(queries, cleanQuery(match[1]))
			
		}
		numQueries += len(queries)
		allQueries[component] = queries
	}


	allQueries["num_of_queries"] = numQueries
	results = append(results, allQueries)
}



func worker(ch chan string, headers []string) {
	defer wg.Done()
	for url := range ch {
		ext(url, headers)
	}
}



func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}



func saveResults(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		fmt.Println("Error encoding results:", err)
	}
}

func cleanQuery(query string) string {
	query = strings.Trim(query, "\"'` \n")
	query = strings.ReplaceAll(query, "\\n", "\n")
	return query
}



func ifExists(list []string, str string) bool {
    for _, s := range list {
        if s == str {
            return true
        }
    }
    return false
}

func hashString(str string) string {
    hash := md5.New()
    hash.Write([]byte(str))
    hashBytes := hash.Sum(nil)
    return fmt.Sprintf("%x", hashBytes)
}