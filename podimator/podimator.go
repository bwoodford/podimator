package podimator

import(
    "fmt"
    "errors"
    "os"
    "time"

    "github.com/cavaliercoder/grab"
    "github.com/mmcdole/gofeed"

    "github.com/IveGotNorto/podimator/config"
    "github.com/IveGotNorto/podimator/podcast"
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
    // GoFeed Parser
    FeedParser *gofeed.Parser
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

func (podi *Podimator) Run(com interface{}) {
    podi.Config = config.Parse(podi.ConfigPath)
    podi.Config.Setup()

    switch com.(type) {
    case Update:
        podi.Update(com.(Update))
    default:
        // woo default
        fmt.Println("")
    }
}

func (podi *Podimator) filter(name string) error {
    i, err := podi.podIndex(name)
    if err != nil {
        return fmt.Errorf("unable to find index for podcast: %v", err)
    }
    podi.Config.Podcasts = []podcast.Podcast{podi.Config.Podcasts[i]}
    return nil
}

func (podi *Podimator) podIndex(podcast string) (int, error) {
    for i, p := range podi.Config.Podcasts {
        if p.Name == podcast {
            return i, nil
        }
    }
    return -1, errors.New("invalid podcast name given")
}

func buildRequests(items []*gofeed.Item, downloadPath string) ([]*grab.Request) {
    var reqs []*grab.Request
    for _, i := range items {
        var enclosure *gofeed.Enclosure
        for _, j := range i.Enclosures {
            if j.Type == "audio/mpeg" {
                enclosure = j;
                break
            }
        }
        if enclosure == nil {
            fmt.Fprintf(os.Stderr, "failed to find audio/mpeg in rss enclosure")
            continue
        }
        req, err := grab.NewRequest(downloadPath, enclosure.URL)
        if err != nil {
            fmt.Fprintf(os.Stderr, "%v\n", err)
            continue
        }
        reqs = append(reqs, req)
    }
    return reqs
}

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
					fmt.Printf("[\033[0;33mdownloading\033[0m] %s (%d%%)\033[K\n", resp.Filename, int(100*resp.Progress()))
				}
			}
        }
    }
    t.Stop()
}
