package podimator

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/cheggaaa/pb/v3"
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

func New() *Podimator {
	return &Podimator{
		false,
		false,
		"/etc/podimator/podcasts.toml",
		nil,
		grab.NewClient(),
	}
}

type PodCommand interface {
	Run(podi *Podimator) error
}

func (podi *Podimator) Start(com PodCommand) {
	var err error
	podi.Config, err = config.Parse(podi.ConfigPath)
	podi.Config.Setup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v; exiting...", err)
	} else {
		if err = com.Run(podi); err != nil {
			fmt.Fprintf(os.Stderr, "%v; exiting...", err)
		}
	}
}

// Takes the user podcast list and filters it down to a single specified podcast.
// Filter will remove podcasts from the config struct
func (podi *Podimator) filter(name string) error {
	i, err := findIndex(podi.Config.Podcasts, name)
	if err != nil {
		return fmt.Errorf("%w; unable to filter podcasts", err)
	}
	podi.Config.Podcasts = []*config.Podcast{podi.Config.Podcasts[i]}
	return nil
}

// Searches for a specified podcast name from the user config file and returns the index location, if it exists.
func findIndex(podcasts []*config.Podcast, name string) (int, error) {
	for i, p := range podcasts {
		if p.Name == name {
			return i, nil
		}
	}
	return -1, errors.New("podcast not found in user config")
}

// Creates web requests for downloading specific (items) that are given.
// downloadpath is used to tie the request to a particular local file location
// for use later.
func buildRequests(items []*gofeed.Item, downloadPath string) []*grab.Request {
	var reqs []*grab.Request
	for _, item := range items {
		enc, err := findEnclosure(item.Enclosures)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
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

// Searches for audio/mpeg files that should be the episode.
// The first encountered audio/mpeg type Enclosure is taken as the episode.
func findEnclosure(enclosure []*gofeed.Enclosure) (*gofeed.Enclosure, error) {
	var sel *gofeed.Enclosure
	var err error
	for _, e := range enclosure {
		if e.Type == "audio/mpeg" {
			// Take the first "audio/mpeg" selection
			sel = e
			break
		}
	}
	if sel == nil {
		err = errors.New("unable to find rss enclosure of type audio/mpeg")
	}
	return sel, err
}

// Fulfills downloading requests to an underlying file location concurrently.
func (podi *Podimator) download(reqs []*grab.Request) error {
	respch := podi.Client.DoBatch(workers, reqs...)
	t := time.NewTicker(200 * time.Millisecond)

	errCollect := errors.New("")
	completed := 0
	requests := len(reqs)
	responses := make([]*grab.Response, 0)
	bar := pb.StartNew(requests)

	for completed < requests {
		select {
		case resp := <-respch:
			if resp != nil {
				responses = append(responses, resp)
			}
		case <-t.C:
			for i, resp := range responses {
				if resp != nil && resp.IsComplete() {
					if err := resp.Err(); err != nil {
						errCollect = fmt.Errorf("%w; ", err)
					}
					// mark completed
					responses[i] = nil
					completed++
					// Increment bar by 1
					bar.Add(1)
				}
			}
		}
	}
	bar.Finish()
	t.Stop()
	return errCollect
}
