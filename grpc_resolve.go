package vpnrouting

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:8888"
)

func (r Vpnrouting) grpcResolve(url string, geo *GeoIP) (*Resolver, error) {
	conn, err := grpc.Dial(r.ResolveGrpcAddress(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewServiceDiscoveryServiceClient(conn)

	lat := fmt.Sprint(geo.Lat)
	lng := fmt.Sprint(geo.Lon)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result := new(Resolver)
	response, err := c.GetClosestNode(ctx, &ServiceRequest{Url: url, Lat: lat, Lng: lng})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", response.GetIp())
	result.IP = response.GetIp()
	return result, nil
}
