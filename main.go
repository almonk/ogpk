package main

import (
	"encoding/json"
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
	Green = "\033[32m"
	Reset = "\033[0m"
)

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
	case "ghostty":
		cmdArgs = append(cmdArgs, "-pk")
	case "xtem-kitty":
		cmdArgs = append(cmdArgs, "-pk")
	case "iTerm.app":
		cmdArgs = append(cmdArgs, "-pi")
	}

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

func main() {
	var pFlag bool
	var jsonFlag bool

	args := os.Args[1:]
	for i, arg := range args {
		if arg == "--p" {
			pFlag = true
			args = append(args[:i], args[i+1:]...)
			break
		}

		if arg == "--json" {
			jsonFlag = true
			args = append(args[:i], args[i+1:]...)
			break
		}
	}

	if len(args) != 1 {
		fmt.Println("Usage: ogpk <url> [options]")
		// Print out options
		fmt.Println("\nOptions:")
		fmt.Println("  --p\t\tPreview image")
		fmt.Println("  --json\tOutput as JSON")

		fmt.Printf("\nVersion %s\n", version)
		return
	}
	url := args[0]

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	doc, err := fetchHTML(url)
	if err != nil {
		log.Fatalf("Error fetching URL: %v", err)
	}

	ogData := extractOpenGraphData(doc)

	// Collect and sort the keys
	var keys []string
	for k := range ogData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Display the OpenGraph data in order by key name
	if jsonFlag {
		// Convert the map to JSON
		jsonData, err := json.MarshalIndent(ogData, "", "  ")
		if err != nil {
			log.Fatalf("Error converting data to JSON: %v", err)
		}
		fmt.Println(string(jsonData))
	} else {
		// Display the data in the terminal
		for _, k := range keys {
			fmt.Printf("%s%s%s: %s\n", Green, k, Reset, ogData[k])
		}
	}

	if imageURL, ok := ogData["og:image"]; ok && pFlag {
		imgData, err := downloadImage(imageURL)
		if err != nil {
			log.Fatalf("Error downloading image: %v", err)
		}

		filePath, err := saveImage(imgData)
		if err != nil {
			log.Fatalf("Error saving image: %v", err)
		}

		// Check if timg is available
		_, err = exec.LookPath("timg")
		if err == nil {
			// If timg is available, display the image using timg
			if err := displayImageWithTimg(filePath); err != nil {
				log.Fatalf("Error displaying image with timg: %v", err)
			}
		} else {
			fmt.Println("timg not found, image saved to:", filePath)
		}
	}

}
