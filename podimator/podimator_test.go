package podimator

import (
	"errors"
	"os"
	"testing"

	"github.com/mmcdole/gofeed"

	"github.com/IveGotNorto/podimator/config"
	"github.com/IveGotNorto/podimator/test"
)

var podi Podimator
var items []*gofeed.Item

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func setup() {
	podi = New()
	podi.Config = &config.Config{
		Location: "",
		Podcasts: test.TestPodcasts,
	}
	items = test.TestItems
}

func TestPodIndex(t *testing.T) {
	var tests = []struct {
		input    string
		expected int
		err      error
	}{
		{"Automated Humans, Automating Humans", 0, nil},
		{"Bill Bryson Sports Podcast", 1, nil},
		{"Gene Simmons Hardcore History", 2, nil},
		{"I Don't Exist!", -1, PodcastNotFound},
	}
	for _, test := range tests {
		got, err := findIndex(podi.Config.Podcasts, test.input)
		if got != test.expected || !errors.Is(err, test.err) {
			t.Errorf("Podimator.podIndex(%q) = %v", test.input, test.expected)
		}
	}
}

func TestBuildRequests(t *testing.T) {
	var tests = []struct {
		input		*gofeed.Item
		expected	int
	}{
		{items[0], 1},
		{items[1], 1},
		{items[2], 0},
		{items[3], 0},
	}
	for i, test := range tests {
		got := buildRequests([]*gofeed.Item{test.input}, "")
		if len(got) != test.expected {
			t.Errorf("Test buildRequest(#%q) = len(%v)", i+1, test.expected)
		}
	}
}

func TestFilter(t *testing.T) {
	// Might need to make a copy of podi and use it here*
	var tests = []struct {
		input string
		expected int
	}{
		{"", len(podi.Config.Podcasts)},
		{"I Don't Exist!", len(podi.Config.Podcasts)},
		{"Automated Humans, Automating Humans", 1},
	}
	for _, test := range tests {
		podi.filter(test.input)
		if len(podi.Config.Podcasts) != test.expected {
			t.Errorf("Test Podimator.filter(%q) = len(%v)", test.input, test.expected)
		}
	}
}

// Need some form of indirection built in this method for testing
/*
func TestDownloadPass(t *testing.T) {
}
*/
