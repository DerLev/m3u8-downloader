package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
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

	if isValidUrl(*urlPtr) {
		url = *urlPtr
	} else {
		inputUrl, err := getUrlFromUser()
		if err != nil {
			fmt.Printf("An error occurred: %v\n", err)
			return
		}
		url = inputUrl
	}

	fmt.Println(url)
}

func getUrlFromUser() (string, error) {
	var err error = nil

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a URL to be downloaded: ")

	input, readErr := reader.ReadString('\n')
	if readErr != nil {
		return "", readErr
	}

	input = strings.TrimSpace(input)

	if isValidUrl(input) {
		return input, err
	} else {
		err = errors.New("the URL is not valid")
		return "", err
	}
}

func isValidUrl(rawUrl string) bool {
	match, _ := regexp.MatchString(`^http(s)?:\/\/.*\..*\/$`, rawUrl)
	return match
}
