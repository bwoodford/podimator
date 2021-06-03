package config

import (
	"testing"
)

type ValidationTests struct {
	in   *Config
	want int
}

var ValidationConfigs = []*ValidationTests{
	{
		in: &Config{
			Location: "",
			Podcasts: nil,
		},
		want: 1,
	},
	{
		in: &Config{
			Location: "",
			Podcasts: []*Podcast{
				{
					URL:  "www.google.com",
					Name: "",
				},
				{
					URL:  "",
					Name: "Bill Bryson Sports Podcast",
				},
				{
					URL:  "",
					Name: "",
				},
			},
		},
		want: 5,
	},
	{
		in: &Config{
			Location: "A real location",
			Podcasts: []*Podcast{
				{
					URL:  "www.valid.com",
					Name: "Valid",
				},
				{
					URL:  "www.veryvalid.com",
					Name: "Very Valid",
				},
				{
					URL:  "www.veryveryvalid.com",
					Name: "Very Very Valid",
				},
			},
		},
		want: 0,
	},
}

func TestValidate(t *testing.T) {
	for _, conf := range ValidationConfigs {
		err := validate(conf.in)
		if len(err) != conf.want {
			t.Errorf("Test validate(config *Config), expected: %d errors, got: %d errors", conf.want, len(err))
		}
	}
}
