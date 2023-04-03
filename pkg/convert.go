package pkg

import (
	"strconv"
	"strings"
)

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

func (s StrTo) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

func (s StrTo) UInt32() (uint32, error) {
	v, err := strconv.Atoi(s.String())
	return uint32(v), err
}

func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}

type Strs []string

func (m Strs) Scan(val interface{}) ([]string, error) {
	s := val.(string)
	ss := strings.Split(string(s), "|")
	return ss, nil
}

func (m Strs) Value() (string, error) {
	str := strings.Join(m, "|")
	return str, nil
}
