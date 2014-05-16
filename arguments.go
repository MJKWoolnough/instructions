package instructions

import (
	"reflect"
	"strconv"
)

type argument interface {
	Set(string) error
	Get() reflect.Value
}

type stringT string

func (s *stringT) Set(a string) error {
	*s = stringT(a)
	return nil
}

func (s *stringT) Get() reflect.Value {
	return reflect.ValueOf(string(*s))
}

type intT struct {
	bitSize int
	number  int64
}

func (i *intT) Set(a string) (err error) {
	i.number, err = strconv.ParseInt(a, 0, i.bitSize)
	return err
}

func (i *intT) Get() reflect.Value {
	switch i.bitSize {
	case 0:
		return reflect.ValueOf(int(i.number))
	case 8:
		return reflect.ValueOf(int8(i.number))
	case 16:
		return reflect.ValueOf(int16(i.number))
	case 32:
		return reflect.ValueOf(int32(i.number))
	}
	return reflect.ValueOf(i.number)
}

type uintT struct {
	bitSize int
	number  uint64
}

func (u *uintT) Set(a string) (err error) {
	u.number, err = strconv.ParseUint(a, 0, u.bitSize)
	return err
}

func (u *uintT) Get() reflect.Value {
	switch u.bitSize {
	case 0:
		return reflect.ValueOf(uint(u.number))
	case 8:
		return reflect.ValueOf(uint8(u.number))
	case 16:
		return reflect.ValueOf(uint16(u.number))
	case 32:
		return reflect.ValueOf(uint32(u.number))
	}
	return reflect.ValueOf(u.number)
}

type floatT struct {
	bitSize int
	number  float64
}

func (f *floatT) Set(a string) (err error) {
	f.number, err = strconv.ParseFloat(a, f.bitSize)
	return err
}

func (f *floatT) Get() reflect.Value {
	if f.bitSize == 32 {
		return reflect.ValueOf(float32(f.number))
	}
	return reflect.ValueOf(f.number)
}
