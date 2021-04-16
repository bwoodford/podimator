package test

import(
	"github.com/mmcdole/gofeed"

	"github.com/IveGotNorto/podimator/podcast"
)

var TestPodcasts = []*podcast.Podcast{
	{
		URL:     "www.google.com",
		Name:    "Automated Humans, Automating Humans",
		Updated: "01/10/1001",
	},
	{
		URL:     "www.yahoo.com",
		Name:    "Bill Bryson Sports Podcast",
		Updated: "01/10/1001",
	},
	{
		URL:     "www.aol.com",
		Name:    "Gene Simmons Hardcore History",
		Updated: "01/10/1001",
	},
}

var TestItems = []*gofeed.Item{
	&gofeed.Item{
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
	&gofeed.Item{
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
	&gofeed.Item{
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
	&gofeed.Item{
		Enclosures: []*gofeed.Enclosure{
		},
	},
}
