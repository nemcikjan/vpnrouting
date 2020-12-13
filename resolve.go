package vpnrouting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Resolver -
type Resolver struct {
	IP string `json:"ip"`
}

func resolve(address string, c GeoIP) Resolver {
	postBody, _ := json.Marshal(map[string]string{
		"url": "www.example.com",
		"Lat": fmt.Sprint(c.Lat),
		"Lon": fmt.Sprint(c.Lon),
	})
	responseBody := bytes.NewBuffer(postBody)
	response, err := http.Post("http://localhost:3000/discovery/closest-node", "application/json", responseBody)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	var result Resolver
	// Unmarshal the JSON byte slice to a GeoIP struct
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
	}

	return result
}
