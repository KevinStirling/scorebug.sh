package mlbstats

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func New() *Client {
	return &Client{
		BaseURL: "https://statsapi.mlb.com",
		HTTP: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Schedule(date *time.Time) (Schedule, error) {
	var d time.Time
	if date == nil {
		d = time.Now()
	} else {
		d = *date
	}
	url := fmt.Sprintf("%s/api/v1/schedule?sportId=1&date=%s&hydrate=linescore,team", c.BaseURL, d.Format("01/02/2006"))

	resp, err := c.HTTP.Get(url)
	if err != nil {
		return Schedule{}, err
	}
	defer resp.Body.Close()

	var out Schedule
	json.NewDecoder(resp.Body).Decode(&out)

	return out, nil
}
