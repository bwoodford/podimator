package podimator

import(
    "fmt"
    "os"
    "net/http"
//    "io"
//    "path"
//    "time"

    "github.com/mmcdole/gofeed"

    . "github.com/IveGotNorto/podimator/internal/config"
    . "github.com/IveGotNorto/podimator/internal/podcast"
)

func Podimator() {

    feed := make(chan *gofeed.Feed)

    config := ConfigParse("podcasts.json")
    config.Setup()

    // only iter on selected podcasts
    // TODO: look into pflag

    go func() {
        defer close(feed)
        for _, podcast := range config.Podcasts {
            parsed, err := feedParse(podcast)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", podcast.Name, err)
                continue
            }
            feed <- parsed
        }
    }()

    for f := range feed {
        download(f.Items[0].Enclosures[0].URL, config.Location)
    }

}

func feedParse(podcast Podcast) (*gofeed.Feed, error) {
    fp := gofeed.NewParser()
    return fp.ParseURL(podcast.Url)
}

func download(url string, filePath string) (err error) {

    req := http.NewRequest(, url,)

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    //filePath = fmt.Sprintf("%v/%v", filePath, path.Base(url))

    //fmt.Println(filePath)

    /*
    // Create the file
    out, err := os.Create(path)
    if err != nil  {
        return err
    }
    defer out.Close()

    
    // Writer the body to file
    _, err = io.Copy(out, resp.Body)
    if err != nil  {
        return err
    }
    */

    return nil  
}


