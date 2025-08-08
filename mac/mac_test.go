package mac

import (
	"testing"

	ut "github.com/cuityhj/cement/unittest"
)

func TestParseMAC(t *testing.T) {
	macs := []string{
		"00:00:5e:00:53:01",
		"02:00:5e:10:00:00:00:01",
		"00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01",
		"00-00-5e-00-53-01",
		"02-00-5e-10-00-00-00-01",
		"00-00-00-00-fe-80-00-00-00-00-00-00-02-00-5e-10-00-00-00-01",
		"0000.5e00.5301",
		"0200.5e10.0000.0001",
		"0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001",
		"0000:5e00:5301",
		"0000-5e00-5301",
		"00005e005301",
	}

	for _, mac := range macs {
		_, err := ParseMAC(mac)
		ut.Assert(t, err == nil, mac+" is invalid mac")
	}
}
