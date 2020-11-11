// +build gofuzz

package vpnrouting

import (
	"github.com/coredns/coredns/plugin/pkg/fuzz"
)

// Fuzz fuzzes cache.
func Fuzz(data []byte) int {
	w := Vpnrouting{}
	return fuzz.Do(w, data)
}
