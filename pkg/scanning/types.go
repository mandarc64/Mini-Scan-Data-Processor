package scanning

import (
	"encoding/json"
)

const (
	Version = iota
	V1
	V2
)

type Scan struct {
	Ip          string          `json:"ip"`
	Port        uint32          `json:"port"`
	Service     string          `json:"service"`
	Timestamp   int64           `json:"timestamp"`
	DataVersion int             `json:"data_version"`
	Data        json.RawMessage `json:"data"` // Use RawMessage to delay unmarshalling
}

type V1Data struct {
	ResponseBytesUtf8 string `json:"response_bytes_utf8"` // Change from []byte to string for easy JSON handling
}

type V2Data struct {
	ResponseStr string `json:"response_str"`
}
