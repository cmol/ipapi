// Package ipapi allows for easy fetching of IP data, while still retaining the
// rate limiting specified
package ipapi

import (
	"reflect"
	"testing"
	"time"
)

func TestLookup(t *testing.T) {
	tests := []struct {
		name    string
		address string
		fields  string
		want    Response
		wantErr bool
	}{
		{
			name:    "Get google public DNS INET4",
			address: "8.8.8.8",
			fields:  "?fields=status,message,country,countryCode,region,regionName,city,zip,timezone,isp,org,as,query",
			want: Response{
				Query:        "8.8.8.8",
				Status:       "success",
				Country:      "United States",
				CountryCode:  "US",
				Region:       "VA",
				RegionName:   "Virginia",
				City:         "Ashburn",
				ZIP:          "20149",
				Latitude:     nil,
				Longtitude:   nil,
				Timezone:     "America/New_York",
				Offset:       nil,
				ISP:          "Google LLC",
				Organization: "Google Public DNS",
				AS:           "AS15169 Google LLC",
				Mobile:       nil,
				Proxy:        nil,
				Hosting:      nil,
			},
			wantErr: false,
		},
		{
			name:    "Get cloudflare DNS INET6",
			address: "2606:4700:4700::1111",
			fields:  "?fields=status,message,country,countryCode,region,regionName,city,zip,timezone,isp,org,as,query",
			want: Response{
				Query:        "2606:4700:4700::1111",
				Status:       "success",
				Country:      "Canada",
				CountryCode:  "CA",
				Region:       "QC",
				RegionName:   "Quebec",
				City:         "Montreal",
				ZIP:          "H4X",
				Latitude:     nil,
				Longtitude:   nil,
				Timezone:     "America/Toronto",
				Offset:       nil,
				ISP:          "Cloudflare, Inc.",
				Organization: "Cloudflare, Inc.",
				AS:           "AS13335 Cloudflare, Inc.",
				Mobile:       nil,
				Proxy:        nil,
				Hosting:      nil,
			},
			wantErr: false,
		},
		{
			name:    "Get RFC1918 address",
			address: "192.168.0.1",
			fields:  "?fields=status,message,country,countryCode,region,regionName,city,zip,timezone,isp,org,as,query",
			want: Response{
				Query:      "192.168.0.1",
				Status:     "fail",
				Message:    "private range",
				Latitude:   nil,
				Longtitude: nil,
				Offset:     nil,
				Mobile:     nil,
				Proxy:      nil,
				Hosting:    nil,
			},
			wantErr: false,
		},
		{
			name:    "Get non IP address",
			address: "1 2 3 4",
			fields:  "?fields=status,message,country,countryCode,region,regionName,city,zip,timezone,isp,org,as,query",
			want: Response{
				Query:      "1 2 3 4",
				Status:     "fail",
				Message:    "invalid query",
				Latitude:   nil,
				Longtitude: nil,
				Offset:     nil,
				Mobile:     nil,
				Proxy:      nil,
				Hosting:    nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Fields = tt.fields
			c, err := Lookup(tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			select {
			case got := <-c:
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Lookup() =\n%+v, want\n%+v", got, tt.want)
				}
			case <-time.After(10 * time.Second):
				t.Errorf("Timed out waiting for response")
			}
		})
	}
}
