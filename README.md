# Podimator

Podimator is a linux terminal program for downloading podcast episodes.

Multiple go packages make Podimator possible:
- [grab](github.com/cavaliercoder/grab)
- [pb](github.com/cheggaaa/pb)
- [gofeed](github.com/mmcdole/gofeed)
- [go-toml](github.com/pelletier/go-toml)
- [urfave/cli](github.com/urfave/cli)

Requires Go v1.7 and above.

## Installation 

1. `git clone` this repository to /usr/local/src, or a prefered location.
  - `git clone https://www.github.com/IveGotNorto/podimator`

2. `cd` to the cloned directory
  - `cd podimator`

3. `cd` to into the cmd/podimator directory.
  - `cd cmd/podimator`

4. Run `go install` 

5. The podimator program has now been added to GOBIN.

## Getting Started

1. Create a new configuration directory for podimator under `~/.config/podimator` 
  - `mkdir ~/.config/podimator`

2. Copy the podcasts.toml (located in the git repository) file to `~/.config/podimator`
  - `cp /usr/local/src/podimator/podcasts.toml ~/.config/podimator`

3. Edit the podcasts.toml file to your preferences.

### Configuration File
An example configuration file (podcasts.toml) is included with the repository. An explanation of the fields in the file follow:
  - location= the download location for all of the podcasts.
  - [[podcasts]] is a required array heading for all podcasts in the file.
    - name= the podcast name that you would like to refer to the podcast by.
      - Every podcast will create a directory under *location* using this specified name.
    - url= a valid Url pointing to the podcasts rss feed.
      - These Urls can usually be found on a podcasts main website.

## Future Improvements
- Add XDG support for configuration file.
- Automatic downloading of new episodes.
- Updating based on last updated date.
