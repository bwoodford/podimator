package podimator

import (
	"fmt"
	"os"

	"github.com/mmcdole/gofeed"
)

type All struct {
	// Podcast name to update
	PodcastName string
}

func (all All) Run(podi *Podimator) error {
	if len(all.PodcastName) > 0 {
		if err := podi.filter(all.PodcastName); err != nil {
			return fmt.Errorf("unable to filter podcast %s: %v", all.PodcastName, err)
		}
	}

	parser := gofeed.NewParser()

	for _, p := range podi.Config.Podcasts {
		feed, err := parser.ParseURL(p.URL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to parse %s: %s\n", p, err)
			continue
		}
		requests := buildRequests(feed.Items, podi.Config.Location+"/"+p.Name)
		fmt.Printf("[\033[0;35mupdating\033[0m] %s\n", p.Name)
		podi.download(requests)
	}
	return nil
}
