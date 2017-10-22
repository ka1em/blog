package model

import (
	"errors"
	"net"

	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func stringToIP(s string) (net.IP, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, errors.New("invalid ip address")
	}
	return ip.To4(), nil
}

// machineID retrieves the private IP address of the Amazon EC2 instance
// and returns its lower 16 bits.
func machineID() (uint16, error) {
	ip, err := stringToIP("47.93.11.105")
	//todo ip ?
	if err != nil {
		return 0, err
	}

	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

func init() {
	var st sonyflake.Settings
	st.MachineID = machineID
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}
