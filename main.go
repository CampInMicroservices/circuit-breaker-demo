package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	gclient "github.com/machinebox/graphql"
)

type LocationResponse struct {
	Cities struct {
		Data []struct {
			City         string  `json:"city"`
			Latitude     float64 `json:"latitude"`
			Longitude    float64 `json:"longitude"`
			Population   int     `json:"population"`
			TempC        int     `json:"tempC"`
			Weather      string  `json:"weather"`
			WeatherShort string  `json:"weatherShort"`
		} `json:"data"`
		Error string `json:"error"`
	} `json:"cities"`
}

func DoGqlReq(u string) (interface{}, error) {

	url := fmt.Sprintf("http://%s/v1/locations", u)

	graphqlClient := gclient.NewClient(url)

	query := `{
				cities {
					data {
						city
						latitude
						longitude
						population
					}
					error
				}
			}`

	graphqlRequest := gclient.NewRequest(query)

	var graphqlResponse interface{}
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		log.Fatalf("Error querying GeoDB, error: %v", err)
	}

	// Convert map to json string
	rJSON, err := json.Marshal(graphqlResponse)
	if err != nil {
		log.Panic("Cannot marshal graphqlResponse to JSON")
	}

	// Convert struct
	var locationsResponse LocationResponse
	err = json.Unmarshal(rJSON, &locationsResponse)
	if err != nil {
		log.Panic("Cannot unmarshal LocationResponse")
	}

	return locationsResponse, nil
}

func main() {

	for i := 0; i < 50; i++ {
		st, err := DoGqlReq("0.0.0.0:8081")

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(st)
		time.Sleep(100 * time.Millisecond)

	}
}
