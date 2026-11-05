//go:build !android && (linux || cgo)

package kiss

// bleScanDedup decides which scan sightings are worth reporting. A device
// is reported on first sight and again when its name first arrives, because
// the name often lives in the scan response and the first sighting lacks it.
type bleScanDedup struct {
	named map[string]bool // addr -> a name has been reported
}

func newBLEScanDedup() *bleScanDedup {
	return &bleScanDedup{named: make(map[string]bool)}
}

// shouldReport records the sighting and reports whether it should be emitted.
func (d *bleScanDedup) shouldReport(addr, name string) bool {
	hadName, seen := d.named[addr]
	if !seen {
		d.named[addr] = name != ""
		return true
	}
	if !hadName && name != "" {
		d.named[addr] = true
		return true
	}
	return false
}

// unnamed returns the addresses reported without a name, in no fixed order.
func (d *bleScanDedup) unnamed() []string {
	var out []string
	for addr, named := range d.named {
		if !named {
			out = append(out, addr)
		}
	}
	return out
}
