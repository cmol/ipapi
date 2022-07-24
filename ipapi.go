// Package ipapi allows for easy fetching of IP data, while still retaining the
// rate limiting specified
package ipapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Fields configures what to query from the API. This can be either as a comma
// separated string or the string of the numeric value eg. "61439". For full
// documentation see: https://ip-api.com/docs/api:json
// Note that
var Fields = "?fields=status,message,country,countryCode,region,regionName,city,zip,lat,lon,timezone,isp,org,as,query"

// APIKey holds the key for a paying customer. Blank by default
var APIKey = ""

// Endpoint is the query location for all queries
var Endpoint = "http://ip-api.com/json/"

// MaxQueueLength limits the total number of elements in the queue. This is
// only relevant for rate limited usage.
var MaxQueueLength = 50

// Response holds data for each of the possible data points
type Response struct {
	Query         string   `json:"query,omitempty"`
	Status        string   `json:"status,omitempty"`
	Message       string   `json:"message,omitempty"`
	Continent     string   `json:"continent,omitempty"`
	ContinentCode string   `json:"continentCode,omitempty"`
	Country       string   `json:"country,omitempty"`
	CountryCode   string   `json:"countryCode,omitempty"`
	Region        string   `json:"region,omitempty"`
	RegionName    string   `json:"regionName,omitempty"`
	City          string   `json:"city,omitempty"`
	District      string   `json:"district,omitempty"`
	ZIP           string   `json:"zip,omitempty"`
	Latitude      *float64 `json:"lat,omitempty"`
	Longtitude    *float64 `json:"lon,omitempty"`
	Timezone      string   `json:"timezone,omitempty"`
	Offset        *int64   `json:"offset,omitempty"`
	Currency      string   `json:"currency,omitempty"`
	ISP           string   `json:"isp,omitempty"`
	Organization  string   `json:"org,omitempty"`
	AS            string   `json:"as,omitempty"`
	ASName        string   `json:"asname,omitempty"`
	Reverse       string   `json:"reverse,omitempty"`
	Mobile        *bool    `json:"mobile,omitempty"`
	Proxy         *bool    `json:"proxy,omitempty"`
	Hosting       *bool    `json:"hosting,omitempty"`
}

var queue = make(chan queueElement, MaxQueueLength)

type queueElement struct {
	address  string
	response chan Response
}

// Lookup adds query to queue and returns the result channel
func Lookup(address string) (chan Response, error) {
	if len(queue) >= MaxQueueLength {
		return nil, errors.New("too many requests in queue")
	}
	c := make(chan Response)
	queue <- queueElement{address, c}
	return c, nil
}

func checkTTLAndSleep(r *http.Response) error {
	ttlString := r.Header.Get("X-Ttl")
	rlString := r.Header.Get("X-Rl")
	rl, err := strconv.Atoi(rlString)
	if err != nil {
		return errors.New("unable to get X-Rl parameter")
	}
	if rl > 0 {
		return nil
	}
	ttl, err := strconv.Atoi(ttlString)
	if err != nil {
		return errors.New("unable to get X-Ttl parameter")
	}
	time.Sleep(time.Duration(ttl) * time.Second)
	return nil
}

// Run will start processing elements in the queue
func Run() {
	go run()
}

func run() {
	for q := range queue {
		for i := 2; i > 0; i-- { // Try multiple times for rate limit
			resp, err := http.Get(Endpoint + q.address + Fields)
			if err != nil {
				log.Fatalln(err)
				resp.Body.Close()
				break
			}

			if resp.StatusCode != http.StatusOK &&
				resp.StatusCode != http.StatusTooManyRequests {
				log.Fatalln(resp.Status)
				resp.Body.Close()
				break
			}

			if resp.StatusCode == http.StatusTooManyRequests {
				if err := checkTTLAndSleep(resp); err != nil {
					log.Println(err.Error())
					time.Sleep(10 * time.Second)
				}
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
				resp.Body.Close()
				break
			}
			var result Response
			if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
				fmt.Println("Can not unmarshal JSON")
			}
			q.response <- result
			if err := checkTTLAndSleep(resp); err != nil {
				log.Println(err.Error())
				time.Sleep(10 * time.Second)
			}
		}
		close(q.response)
	}
}
