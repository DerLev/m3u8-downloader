package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

func main() {
	/* print some ascii art */
	fmt.Println(`
       _       _                                          
 ._ _  _)     (_)    _|  _       ._  |  _   _.  _|  _  ._ 
 | | | _) |_| (_)   (_| (_) \/\/ | | | (_) (_| (_| (/_ |  
                                                          
  `)

	urlPtr := flag.String("url", "", "url to be downloaded from")
	// uaPtr := flag.String("user-agent", "", "user agent to be used for downloading")
	// referrerPtr := flag.String("referrer", "", "referrer to be used for downloading")
	// dataPtr := flag.String("data", "", "base64 config for use with other applications")
	flag.Parse()

	url := ""

	/* check if the url passed by flags is valid, otherwise get the input from the user */
	if isValidUrl(*urlPtr) {
		url = *urlPtr
	} else {
		/* get the url from the user */
		inputUrl, err := getUrlFromUser()
		if err != nil {
			fmt.Printf("An error occurred: %v\n", err)
			return
		}
		url = inputUrl
	}

	fmt.Println(url)
	/* download the file */
	now := time.Now().Format("2006-01-02 15-04-05")
	downloadFile(url, now)
}

/* prompts the user for the url input */
func getUrlFromUser() (string, error) {
	var err error = nil

	/* create the prompt */
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a URL to be downloaded: ")

	/* read the input from the cli */
	input, readErr := reader.ReadString('\n')
	if readErr != nil {
		return "", readErr
	}

	input = strings.TrimSpace(input)

	/* check whether the input is a valid url */
	if isValidUrl(input) {
		return input, err
	} else {
		err = errors.New("the URL is not valid")
		return "", err
	}
}

/* checks whether a string is a valid url */
func isValidUrl(rawUrl string) bool {
	match, _ := regexp.MatchString(`^http(s)?:\/\/.*\..*\/.*$`, rawUrl)
	return match
}

/* downloads a file from an http endpoint */
func downloadFile(urlInput string, subpath string) {
	/* directory to download into */
	dir := "m3u8s/" + subpath

	/* create the directory */
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Println("Error when creating directory:", err)
		return
	}

	/* parse the url */
	parsedUrl, err := url.Parse(urlInput)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	/* extract the filename from the url */
	filename := path.Base(parsedUrl.Path)

	/* fetch the file */
	res, err := http.Get(urlInput)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer res.Body.Close()

	/* return error if the response code is not in the 200s */
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		fmt.Println("Bad response status code:", res.Status)
		return
	}

	/* create the output file */
	out, err := os.Create(dir + "/" + filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer out.Close()

	/* copy the body's data into the output file */
	bytes, err := io.Copy(out, res.Body)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}

	fmt.Printf("Downloaded %d bytes to %s\n", bytes, filename)
}
