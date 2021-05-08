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

func (r Vpnrouting) resolve(address string, c *GeoIP) (*Resolver, error) {
	postBody, _ := json.Marshal(map[string]string{
		"url": address,
		"Lat": fmt.Sprint(c.Lat),
		"Lon": fmt.Sprint(c.Lon),
	})
	responseBody := bytes.NewBuffer(postBody)
	response, err := http.Post(r.ResolverHostname()+`/discovery/closest-node`, "application/json", responseBody)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	result := new(Resolver)
	// Unmarshal the JSON byte slice to a GeoIP struct
	err = json.Unmarshal(body, result)
	if err != nil {
		fmt.Println(err)
	}

	return result, nil
}
