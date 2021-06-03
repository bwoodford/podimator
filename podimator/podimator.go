package podimator

import (
	"fmt"
	"os"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/mmcdole/gofeed"

	"github.com/IveGotNorto/podimator/config"
)

const workers int = 5

type Podimator struct {
	// Debug printing statements
	Debug bool
	// Verbose output to terminal
	Verbose bool
	// Path to configuration file
	ConfigPath string
	// Parsed configuration file
	Config *config.Config
	// HTTP Grab client
	Client *grab.Client
}

func New() Podimator {
	return Podimator{
		false,
		false,
		"/etc/podimator/podcasts.toml",
		nil,
		grab.NewClient(),
	}
}

type Command interface {
	Run(podi *Podimator) error
}

func (podi *Podimator) Start(com Command) error {
	var errors = []error{}
	podi.Config, errors = config.Parse(podi.ConfigPath)
	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Configuration file errors encountered:\n")
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		return nil
	}
	podi.Config.Setup()
	return com.Run(podi)
}

// filter will take the user podcast list and "filter" it down to a single specified podcast.
// This method is used heavily when a single podcast is specified through the user interface.
func (podi *Podimator) filter(name string) error {
	i, err := findIndex(podi.Config.Podcasts, name)
	if err != nil {
		return fmt.Errorf("%v: %v", "unable to filter podcasts", err)
	}
	podi.Config.Podcasts = []*config.Podcast{podi.Config.Podcasts[i]}
	return nil
}

//findIndex searches for a specified podcast name from the user config file and returns the index location, if it exists.
func findIndex(podcasts []*config.Podcast, name string) (int, error) {
	for i, p := range podcasts {
		if p.Name == name {
			return i, nil
		}
	}
	return -1, fmt.Errorf("%v", "podcast not found in user config")
}

// buildRequests creates web requests for downloading specific (items) that are given.
// downloadpath is used to tie the request to a particular local file location
// for use later.
func buildRequests(items []*gofeed.Item, downloadPath string) []*grab.Request {
	var reqs []*grab.Request
	for _, item := range items {
		enc, err := findEnclosure(item.Enclosures)
		if err != nil {
			//fmt.Fprintf(os.Stderr, "failed to find audio/mpeg in rss enclosure")
			continue
		}
		req, err := grab.NewRequest(downloadPath, enc.URL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		reqs = append(reqs, req)
	}
	return reqs
}

// findEnclosure searches for audio/mpeg files that should be the episode.
// The first encountered audio/mpeg type Enclosure is taken as the episode.
func findEnclosure(enclosure []*gofeed.Enclosure) (*gofeed.Enclosure, error) {
	var sel *gofeed.Enclosure
	var err error
	for _, e := range enclosure {
		if e.Type == "audio/mpeg" {
			sel = e
			break
		}
	}
	if sel == nil {
		err = fmt.Errorf("%v", "unable to find rss enclosure of type audio/mpeg")
	}
	return sel, err
}

// download fulfills downloading those requests to an underlying file location concurrently.
// Number of workers for this method is specified globally.
func (podi *Podimator) download(reqs []*grab.Request) {
	respch := podi.Client.DoBatch(workers, reqs...)
	t := time.NewTicker(80 * time.Millisecond)

	completed := 0
	inProgress := 0
	responses := make([]*grab.Response, 0)

	for completed < len(reqs) {
		select {
		case resp := <-respch:
			if resp != nil {
				responses = append(responses, resp)
			}
		case <-t.C:
			if inProgress > 0 {
				fmt.Printf("\033[%dA\033[K", inProgress)
			}

			// update completed downloads
			for i, resp := range responses {
				if resp != nil && resp.IsComplete() {
					// print final result
					if err := resp.Err(); err != nil {
						fmt.Fprintf(os.Stderr, "[\033[0;31merror\033[0m] failed to download from %s: %v\n", resp.Request.URL(), err)
					} else {
						fmt.Printf("[\033[0;32mcomplete\033[0m] downloaded %s\n", resp.Filename)
					}
					// mark completed
					responses[i] = nil
					completed++
				}
			}
			// update downloads in progress
			inProgress = 0
			for _, resp := range responses {
				if resp != nil {
					inProgress++
					if podi.Verbose {
						fmt.Printf("[\033[0;33mdownloading\033[0m] %s (%d%%)\033[K\n", resp.Filename, int(100*resp.Progress()))
					} else {
						fmt.Printf("[\033[0;33mdownloading\033[0m] %s\033[K\n", resp.Filename)
					}
				}
			}
		}
	}
	t.Stop()
}
