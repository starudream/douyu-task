package ws

import (
	"bytes"
	"encoding/binary"
	"strings"
)

const (
	client2Server uint16 = 689
	server2Client uint16 = 690
)

func Encode(kv ...string) []byte {
	sb := bytes.Buffer{}
	for i := 0; i < len(kv); i += 2 {
		sb.WriteString(escape(kv[i]))
		sb.WriteString("@=")
		sb.WriteString(escape(kv[i+1]))
		sb.WriteString("/")
	}
	return EncodeRaw(sb.Bytes())
}

func EncodeRaw(data []byte) []byte {
	length := uint32(13 + len(data))
	bs := make([]byte, length)
	binary.LittleEndian.PutUint32(bs[0:4], length-4)
	binary.LittleEndian.PutUint32(bs[4:8], length-4)
	binary.LittleEndian.PutUint16(bs[8:10], client2Server)
	bs[10] = 0
	bs[11] = 0
	copy(bs[12:length-1], data)
	bs[length-1] = 0
	return bs
}

func Decode(data []byte) map[string]string {
	bs := DecodeRaw(data)
	if bs == nil {
		return map[string]string{}
	}
	raw := string(bs)
	m := make(map[string]string)
	for _, v := range strings.Split(raw, "/") {
		if len(v) == 0 {
			continue
		}
		a := strings.Split(v, "@=")
		if len(a) != 2 {
			continue
		}
		b := unescape(a[0])
		c := unescape(a[1])
		m[b] = c
	}
	return m
}

func DecodeRaw(data []byte) []byte {
	if len(data) < 13 {
		return nil
	}
	l1 := binary.LittleEndian.Uint32(data[0:4])
	l2 := binary.LittleEndian.Uint32(data[4:8])
	t := binary.LittleEndian.Uint16(data[8:10])
	if l1 != l2 || t != server2Client || data[10] != 0 || data[11] != 0 || data[len(data)-1] != 0 {
		return nil
	}
	return data[12 : l1+3]
}

func escape(s string) string {
	s = strings.ReplaceAll(s, "@", "@A")
	s = strings.ReplaceAll(s, "/", "@S")
	return s
}

func unescape(s string) string {
	s = strings.ReplaceAll(s, "@A", "@")
	s = strings.ReplaceAll(s, "@S", "/")
	return s
}
