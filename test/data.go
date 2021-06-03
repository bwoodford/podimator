package test

import (
	"github.com/mmcdole/gofeed"

	"github.com/IveGotNorto/podimator/config"
)

var TestPodcasts = []*config.Podcast{
	{
		URL:  "www.google.com",
		Name: "Automated Humans, Automating Humans",
	},
	{
		URL:  "www.yahoo.com",
		Name: "Bill Bryson Sports Podcast",
	},
	{
		URL:  "www.aol.com",
		Name: "Gene Simmons Hardcore History",
	},
	{
		URL:  "",
		Name: "Bleh",
	},
	{
		URL:  "www.askjeeves.com",
		Name: "",
	},
	{},
}

var TestItems = []*gofeed.Item{
	{
		Enclosures: []*gofeed.Enclosure{
			{
				URL:  "https://woo.com",
				Type: "audio/mpeg",
			},
			{
				URL:  "https://thislinkshouldbeskipped.com",
				Type: "text",
			},
		},
	},
	{
		Enclosures: []*gofeed.Enclosure{
			{
				URL:  "https://thislinkshouldbeskipped.com",
				Type: "text",
			},
			{
				URL:  "https://woo.com",
				Type: "audio/mpeg",
			},
		},
	},
	{
		Enclosures: []*gofeed.Enclosure{
			{
				URL:  "https://thislinkshouldbeskipped.com",
				Type: "text",
			},
			{
				URL:  "https://thislinkshouldbeskipped.com",
				Type: "text",
			},
		},
	},
	{
		Enclosures: []*gofeed.Enclosure{},
	},
}
