package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegration_OpenGraphExtraction(t *testing.T) {
	// Mock an HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<meta property="og:title" content="Test Title" />
				<meta property="og:description" content="Test Description" />
			</head>
			<body>
				Hello, World!
			</body>
			</html>
		`))
	}))
	defer ts.Close()

	ogData, err := getOpenGraphData(ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTitle := "Test Title"
	if ogData["og:title"] != expectedTitle {
		t.Errorf("Expected og:title to be %s, but got %s", expectedTitle, ogData["og:title"])
	}

	expectedDescription := "Test Description"
	if ogData["og:description"] != expectedDescription {
		t.Errorf("Expected og:description to be %s, but got %s", expectedDescription, ogData["og:description"])
	}
}

func TestIntegration_OpenGraphExtraction_RealURL(t *testing.T) {
	url := "https://replay.software"

	ogData, err := getOpenGraphData(url)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedData := map[string]string{
		"og:description":  "A tiny studio making delightful apps for your Mac.",
		"og:image":        "https://replay.software/replay/opengraph-23.png",
		"og:image:height": "315",
		"og:image:width":  "600",
		"og:locale":       "en_GB",
		"og:site_name":    "Replay Software",
		"og:title":        "Replay Software",
		"og:type":         "website",
		"og:url":          "https://replay.software/",
	}

	// Check if the OpenGraph data matches the expecte values
	for key, expectedValue := range expectedData {
		if ogData[key] != expectedValue {
			t.Errorf("For key %s, expected value to be %s, but got %s", key, expectedValue, ogData[key])
		}
	}
}
