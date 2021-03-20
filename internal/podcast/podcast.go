package podcast

type Podcast struct {
    Url string `json:"url"`
    Name string `json:"name"`
    Updated string `json:"updated"`
    Process bool
}

// TODO: Create method for updating update date in json file

