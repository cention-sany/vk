package vk

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// EpochTime is time in seconds. Use it to successfully parse with json
type EpochTime time.Time

// MarshalJSON unix timestamp strings
func (t EpochTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

// UnmarshalJSON unix timestamp strings
func (t *EpochTime) UnmarshalJSON(s []byte) error {
	var err error
	var q int64

	if q, err = strconv.ParseInt(string(s), 10, 64); err != nil {
		return err
	}

	*(*time.Time)(t) = time.Unix(q, 0)
	return err
}

// Bool is a bool variable parsed from int
type Bool bool

func (b Bool) MarshalJSON() ([]byte, error) {
	if bool(b) {
		return []byte("1"), nil
	} else {
		return []byte("0"), nil
	}
}

func (b *Bool) UnmarshalJSON(s []byte) error {
	var (
		err error
		q   int64
	)
	if q, err = strconv.ParseInt(string(s), 10, 32); err != nil {
		return err
	}

	*b = Bool(q == 1)
	return err
}

type Int int

func (i Int) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(i))), nil
}

func (i *Int) UnmarshalJSON(s []byte) error {
	q, err := strconv.ParseInt(string(s), 10, 32)
	if err != nil {
		return err
	}
	*i = Int(q)
	return nil
}

// IdList is a list of comma-separated ids
type IdList []int

func (l IdList) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}
func (l IdList) String() string {
	var str = make([]string, len(l))
	for i, v := range []int(l) {
		str[i] = fmt.Sprint(v)
	}
	return strings.Join(str, ",")
}

func (l *IdList) UnmarshalJSON(s []byte) error {
	if string(s) == "" {
		*l = nil
		return nil
	}
	arr := strings.Split(string(s), ",")
	*l = make(IdList, len(arr))
	for i := range arr {
		q, err := strconv.ParseInt(arr[i], 10, 64)
		if err != nil {
			return err
		}
		(*l)[i] = int(q)
	}
	return nil
}
