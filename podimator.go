package podimator

import(
    "fmt"
    "math"
    "os"
    "time"

    "github.com/cavaliercoder/grab"
    "github.com/mmcdole/gofeed"
)

type Podimator struct {
    Config *Config
    Client *grab.Client
}

type Podcast struct {
    URL string `json:"url"`
    Name string `json:"name"`
    Updated string `json:"updated"`
    Process bool
}

var podi Podimator
const workers int = 5

func init() {
    podi = Podimator {
        Config: ConfigParse("podcasts.json"),
        Client: grab.NewClient(),
    }

    podi.Config.Setup()
    podi.Client.UserAgent = "Podimator"
}

// Process existing podcast(s) for all episodes
func All(podcast string) {
    if len(podcast) > 0 && !contains(podcast) {
         return 
    } else if len(podcast) == 0 {
        markProcess(podcast)
    }
    podi.Process(math.MaxInt32)
}

// Process existing podcast(s) for most recent episodes
func Update(podcast string, episodes int) {
    if len(podcast) > 0 && !contains(podcast) {
        // Podcast does not exist
        return
    } else if len(podcast) == 0 {
        markProcess(podcast)
    }
    podi.Process(episodes)
}

func (podi *Podimator) Process(episodes int) {
    fp := gofeed.NewParser()

    for _, podcast := range podi.Config.Podcasts {
        feed, err := fp.ParseURL(podcast.URL)
        if err != nil {
            fmt.Fprintf(os.Stderr, "unable to parse %s: %v\n", podcast.Name, err)
            continue
        }
        items := feed.Items[0:episodes]
        fmt.Printf("Processing \"%v\"\n", podcast.Name)
        requests := getRequests(items, podi.Config.Location + "/" + podcast.Name)
        fmt.Printf("Downloading %d files...\n", len(requests))
        download(requests, podi.Client)
    }
}

func markProcess(name string) {
    for _, podcast := range podi.Config.Podcasts {
        podcast.Process = true
    }
}

func contains(name string) bool {
    for _, podcast := range podi.Config.Podcasts {
        if podcast.Name == name {
            podcast.Process = true
            return true
        }
    }
    return false
}

func getRequests(items []*gofeed.Item, downloadPath string) ([]*grab.Request) {
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

func download(reqs []*grab.Request, client *grab.Client) {

    respch := client.DoBatch(workers, reqs...)

    t := time.NewTicker(200 * time.Millisecond)

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
            // clear lines
            if inProgress > 0 {
                fmt.Printf("\033[%dA\033[K", inProgress)
            }

            // update completed downloads
            for i, resp := range responses {
                if resp != nil && resp.IsComplete() {
                    // print final result
                    if err := resp.Err(); err != nil {
                        fmt.Fprintf(os.Stderr, "downloading failed %s: %v\n", resp.Request.URL(), err)
                    } else {
                        fmt.Printf("Finished %s %d / %d bytes (%d%%)\n", 
                                    resp.Filename, resp.BytesComplete(), 
                                    resp.Size, int(100*resp.Progress()))
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
                    fmt.Printf("Downloading %s %d / %d bytes (%d%%)\033[K\n", 
                                resp.Filename,
                                resp.BytesComplete(), 
                                resp.Size, 
                                int(100*resp.Progress()))
                }
            }
        }
    }

    t.Stop()
    fmt.Printf("%d files successfully downloaded.\n", len(reqs))
}


