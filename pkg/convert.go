package pkg

import (
	"strconv"
	"strings"
)

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

func (s StrTo) Int() (int, *Error) {
	v, err := strconv.Atoi(s.String())
	if err != nil {
		return 0, NewError(ServerError, err.Error())
	}
	return v, nil
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

func (s StrTo) UInt32() (uint32, *Error) {
	v, err := strconv.Atoi(s.String())
	if err != nil {
		return 0, NewError(ServerError, err.Error())
	}
	return uint32(v), nil
}

func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}

type Strs []string

func (m Strs) Scan(val interface{}) ([]string, *Error) {
	s, ok := val.(string)
	if !ok {
		return nil, NewError(ServerError, "fail to scan")
	}
	ss := strings.Split(string(s), "|")
	return ss, nil
}

func (m Strs) Value() (string, *Error) {
	str := strings.Join(m, "|")
	return str, nil
}
