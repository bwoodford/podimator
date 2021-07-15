package podimator

import (
	"os"
	"testing"

	"github.com/mmcdole/gofeed"

	"github.com/IveGotNorto/podimator/config"
	test "github.com/IveGotNorto/podimator/testdata"
)

var podi *Podimator
var items []*gofeed.Item

type BasicTest struct {
	input    string
	expected int
}

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
	var tests = []BasicTest{
		{"Automated Humans, Automating Humans", 0},
		{"Bill Bryson Sports Podcast", 1},
		{"Gene Simmons Hardcore History", 2},
		{"I Don't Exist!", -1},
		{"", -1},
	}
	for _, test := range tests {
		if got, _ := findIndex(podi.Config.Podcasts, test.input); got != test.expected {
			t.Errorf("Podimator.podIndex(%q) = %v", test.input, test.expected)
		}
	}
}

func TestBuildRequests(t *testing.T) {
	var tests = []struct {
		input    *gofeed.Item
		expected int
	}{
		{items[0], 1},
		{items[1], 1},
		{items[2], 0},
		{items[3], 0},
		{items[4], 0},
	}
	for i, test := range tests {
		if got := buildRequests([]*gofeed.Item{test.input}, ""); len(got) != test.expected {
			t.Errorf("Test buildRequest(#%q) = len(%v)", i+1, test.expected)
		}
	}
}

func TestFilter(t *testing.T) {
	// Copy the global struct before making changes to it
	var podiCop = podi
	var tests = []struct {
		input    string
		errors   bool
		expected int
	}{
		{"I Don't Exist!", true, len(podiCop.Config.Podcasts)},
		{"", true, len(podiCop.Config.Podcasts)},
		{"Automated Humans, Automating Humans", false, 1},
	}
	for _, test := range tests {
		err := podiCop.filter(test.input)
		length := len(podiCop.Config.Podcasts)
		if length != test.expected && (err != nil) == test.errors {
			t.Errorf("Test Podimator.filter(%q) = len(%v), expected len(%v)", test.input, length, test.expected)
		}
	}
}
