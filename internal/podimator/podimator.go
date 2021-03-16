package podimator

import(
    "fmt"
    "os"
    "time"

    "github.com/cavaliercoder/grab"
    "github.com/mmcdole/gofeed"

    . "github.com/IveGotNorto/podimator/internal/config"
    . "github.com/IveGotNorto/podimator/internal/podcast"
)

type Podimator struct{
    Config *Config
    Client *grab.Client
}

const workers int = 5

func (podi *Podimator) Process() {
    for _, podcast := range podi.Config.Podcasts {
        feed, err := feedParse(podcast)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", podcast.Name, err)
            continue
        }
        fmt.Printf("Processing \"%v\"\n", podcast.Name)
        requests := episodes(feed.Items, podi.Config.Location + "/" + podcast.Name)
        fmt.Printf("Downloading %d files...\n", len(requests))
        download(requests, podi.Client)
    }
}


func feedParse(podcast Podcast) (*gofeed.Feed, error) {
    fp := gofeed.NewParser()
    return fp.ParseURL(podcast.Url)
}

func episodes(items []*gofeed.Item, downloadPath string) ([]*grab.Request){
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
            fmt.Fprintf(os.Stderr, "Failed to recieve Enclosure from Feed Item\n")
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
                        fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", resp.Request.URL(), err)
                    } else {
                        fmt.Printf("Finished %s %d / %d bytes (%d%%)\n", resp.Filename, resp.BytesComplete(), resp.Size, int(100*resp.Progress()))
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
                    fmt.Printf("Downloading %s %d / %d bytes (%d%%)\033[K\n", resp.Filename, resp.BytesComplete(), resp.Size, int(100*resp.Progress()))
                }
            }
        }
    }

    t.Stop()
    fmt.Printf("%d files successfully downloaded.\n", len(reqs))
}



