package vpnrouting

import (
	"encoding/json"
	"fmt"

	"github.com/maxmind/mmdbinspect/pkg/mmdbinspect"
)

const (
	CityDBPath = "/root/coredns/GeoLite2-City.mmdb"
)

// GeoIP struct
type GeoIP struct {
	// The right side is the name of the JSON variable
	IP string `json:"ip"`
	// CountryCode string  `json:"country_code"`
	// CountryName string  `json:"country_name"`
	// RegionCode  string  `json:"region_code"`
	// RegionName  string  `json:"region_name"`
	// City        string  `json:"city"`
	// Zipcode     string  `json:"zipcode"`
	Lat float32 `json:"latitude"`
	Lon float32 `json:"longitude"`
	// MetroCode   int     `json:"metro_code"`
	// AreaCode    int     `json:"area_code"`
}

// GeoLookup fn
func geoLookup(address string) (*GeoIP, error) {
	// Use freegeoip.net to get a JSON response
	// There is also /xml/ and /csv/ formats available
	// http://ip-api.com/json
	fmt.Println("Geo lookup for:" + address)

	reader, err := mmdbinspect.OpenDB(CityDBPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	records, err := mmdbinspect.RecordsForNetwork(*reader, address)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	stringJSON, err := mmdbinspect.RecordToString(records)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var rec []map[string]map[string]map[string]float32
	json.Unmarshal([]byte(stringJSON), &rec)
	location := rec[0]["Record"]["location"]

	fmt.Println(location["latitude"], location["longitude"])

	// response, err := http.Get("http://api.ipstack.com/" + address + "?access_key=bab3a03dd8ff3f11b2212a2f57c91089")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer response.Body.Close()

	// response.Body() is a reader type. We have
	// to use ioutil.ReadAll() to read the data
	// in to a byte slice(string)
	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Geo lookup result:" + fmt.Sprint(response.Body))
	geo := &GeoIP{
		Lat: location["latitude"],
		Lon: location["longitude"],
		IP:  address,
	}
	// Unmarshal the JSON byte slice to a GeoIP struct
	// err = json.Unmarshal(body, &geo)
	if err != nil {
		fmt.Println(err)
	}

	return geo, nil
}
