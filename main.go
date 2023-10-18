package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

var version = "dev"

const (
	Green        = "\033[32m"
	Reset        = "\033[0m"
	Ghostty      = "ghostty"
	XtemKitty    = "xtem-kitty"
	ITermApp     = "iTerm.app"
	OpenGraphPre = "og:"
)

func main() {
	pFlag := flag.Bool("p", false, "Show og:image")
	jsonFlag := flag.Bool("json", false, "Output as JSON")

	// Parse the flags
	flag.Parse()

	// Get the non-flag arguments
	args := flag.Args()

	// Check if we have a URL
	if len(args) != 1 {
		fmt.Println("Usage: ogpk [options] <url>")
		flag.PrintDefaults()
		fmt.Printf("\nVersion %s\n", version)
		return
	}

	url := parseURL(args[0])
	ogData, err := getOpenGraphData(url)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if *jsonFlag {
		displayAsJSON(ogData)
	} else {
		displayInTerminal(ogData)
	}

	if imageURL, ok := ogData["og:image"]; ok && *pFlag {
		displayImage(imageURL)
	}
}

// parseURL parses a URL and adds the protocol prefix if it's missing.
func parseURL(input string) string {
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		return "http://" + input
	}
	return input
}

// getOpenGraphData fetches OpenGraph data from a given URL.
func getOpenGraphData(url string) (map[string]string, error) {
	doc, err := fetchHTML(url)
	if err != nil {
		return nil, fmt.Errorf("fetching URL: %w", err)
	}
	return extractOpenGraphData(doc), nil
}

// displayAsJSON displays OpenGraph data as JSON.
func displayAsJSON(ogData map[string]string) {
	jsonData, err := json.MarshalIndent(ogData, "", "  ")
	if err != nil {
		log.Fatalf("Error converting data to JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}

// displayInTerminal displays OpenGraph data in the terminal.
func displayInTerminal(ogData map[string]string) {
	var keys []string
	for k := range ogData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s%s%s: %s\n", Green, k, Reset, ogData[k])
	}
}

// displayImage downloads and displays an image from a given URL.
func displayImage(imageURL string) {
	imgData, err := downloadImage(imageURL)
	if err != nil {
		log.Fatalf("Error downloading image: %v", err)
	}

	filePath, err := saveImage(imgData)
	if err != nil {
		log.Fatalf("Error saving image: %v", err)
	}

	_, err = exec.LookPath("timg")
	if err == nil {
		if err := displayImageWithTimg(filePath); err != nil {
			log.Fatalf("Error displaying image with timg: %v", err)
		}
	} else {
		fmt.Println("timg not found, image saved to:", filePath)
	}
}

// fetchHTML fetches and parses HTML from a given URL.
func fetchHTML(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return html.Parse(resp.Body)
}

// extractOpenGraphData extracts OpenGraph data from parsed HTML and returns it as a map.
func extractOpenGraphData(doc *html.Node) map[string]string {
	data := make(map[string]string)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			var property, content string
			for _, a := range n.Attr {
				if a.Key == "property" && a.Val[:3] == "og:" {
					property = a.Val
				}
				if a.Key == "content" {
					content = a.Val
				}
			}
			if property != "" && content != "" {
				data[property] = content
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return data
}

func displayImageWithTimg(path string) error {
	fmt.Println()
	cmdArgs := []string{path}

	switch terminalName() {
	case Ghostty:
		cmdArgs = append(cmdArgs, "-pk")
	case XtemKitty:
		cmdArgs = append(cmdArgs, "-pk")
	case ITermApp:
		cmdArgs = append(cmdArgs, "-pi")
	}

	// Set the height to 12 grid units
	cmdArgs = append(cmdArgs, "-gx12")

	cmd := exec.Command("timg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// terminalName returns the name of the terminal running the application.
func terminalName() string {
	term := os.Getenv("TERM_PROGRAM")
	if term == "" {
		// Fetch the name of the terminal from the $TERM environment variable
		term = os.Getenv("TERM")

		if term == "" {
			return "unknown"
		}
	}
	return term
}

// downloadImage downloads an image from a given URL.
func downloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// saveImage saves image data to a temporary file and returns the file path.
func saveImage(data []byte) (string, error) {
	tmpFile, err := ioutil.TempFile("", "ogpk-*.jpg")
	if err != nil {
		return "", err
	}

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return "", err
	}
	tmpFile.Close()

	return tmpFile.Name(), nil
}
