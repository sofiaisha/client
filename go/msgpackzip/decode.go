package msgpackzip

import (
	"errors"
	"io"
	"math"
)

var ErrExtUnsupported = errors.New("ext types not supported")
var ErrMaxDepth = errors.New("input exceeded maximum allowed depth")
var ErrContainerTooBig = errors.New("container allocation is too big")
var ErrStringTooBig = errors.New("string allocation is too big")
var ErrLenTooBig = errors.New("Lenghts bigger than 0x8000000 are too big")
var ErrIntTooBig = errors.New("Cannot handle ints largers than int64 max")

type count uint64

type intType int

const (
	intTypeFixedUint intType = 1  // <= 0x7f
	intTypeFixedInt  intType = 2  // >= 0xe0
	intTypeUint8     intType = 3  // 0xcc
	intTypeInt8      intType = 4  // 0xd0
	intTypeUint16    intType = 5  // 0xcd
	intTypeInt16     intType = 6  // 0xd1
	intTypeUint32    intType = 7  // 0xce
	intTypeInt32     intType = 8  // 0xd2
	intTypeUint64    intType = 9  // 0xcf
	intTypeInt64     intType = 10 // 0xd3
)

type msgpackInt struct {
	typ  intType
	val  int64
	uval uint64
}

const bigLen = 0x8000000
const bigString = bigLen

func (i msgpackInt) toLen() (int, error) {
	if i.typ == intTypeUint64 {
		if i.uval >= uint64(bigLen) {
			return 0, ErrLenTooBig
		}
		return int(i.uval), nil
	}
	if i.val >= int64(bigLen) {
		return 0, ErrLenTooBig
	}
	return int(i.val), nil
}

func (i msgpackInt) toInt64() (int64, error) {
	if i.typ == intTypeUint64 {
		if i.uval >= uint64(math.MaxInt64) {
			return 0, ErrIntTooBig
		}
		return int64(i.uval), nil
	}
	return i.val, nil
}

type msgpackDecoderHooks struct {
	mapKeyHook      func(d decodeStack) (decodeStack, error)
	mapValueHook    func(d decodeStack) (decodeStack, error)
	mapStartHook    func(d decodeStack, i msgpackInt) (decodeStack, error)
	arrayStartHook  func(d decodeStack, i msgpackInt) (decodeStack, error)
	arrayValueHook  func(d decodeStack, i interface{}) (decodeStack, error)
	stringHook      func(l msgpackInt, s string) error
	rawHook         func(l msgpackInt, b []byte) error
	nilHook         func() error
	intHook         func(i msgpackInt) error
	float32Hook     func(f float32) error
	float64Hook     func(f float64) error
	boolHook        func(b bool) error
	fallthroughHook func(i interface{}, s string) error
}

func readByte(r io.Reader) (byte, error) {
	var buf [1]byte
	_, err := r.Read(buf[:])
	if err != nil {
		return byte(0), err
	}
	return buf[0], nil
}

type decodeStack struct {
	depth int
	hooks msgpackDecoderHooks
}

func (d decodeStack) descend() decodeStack {
	d.depth++
	return d
}

type msgpackDecoder struct {
	r io.Reader
}

func newMsgpackDecoder(r io.Reader) *msgpackDecoder {
	return &msgpackDecoder{r: r}
}

func (m *msgpackDecoder) run(h msgpackDecoderHooks) error {
	d := decodeStack{hooks: h}
	return m.decode(d)
}

func (m *msgpackDecoder) produceInt(s decodeStack, i msgpackInt) (err error) {
	if s.hooks.intHook != nil {
		return s.hooks.intHook(i)
	}
	if s.hooks.fallthroughHook != nil {
		return s.hooks.fallthroughHook(i, "int")
	}
	return nil
}

func (m *msgpackDecoder) decodeString(s decodeStack, i msgpackInt) (err error) {
	l, err := i.toLen()
	if err != nil {
		return err
	}
	if l > bigString {
		return ErrStringTooBig
	}
	buf := make([]byte, l)
	_, err = io.ReadFull(m.r, buf)
	if err != nil {
		return err
	}
	return m.produceString(s, i, string(buf))
}

func (m *msgpackDecoder) produceString(s decodeStack, i msgpackInt, str string) (err error) {
	if s.hooks.stringHook != nil {
		return s.hooks.stringHook(i, str)
	}
	if s.hooks.fallthroughHook != nil {
		return s.hooks.fallthroughHook(str, "string")
	}
	return nil
}

func (m *msgpackDecoder) decode(s decodeStack) (err error) {
	if s.depth > 100 {
		return ErrMaxDepth
	}

	b, err := readByte(m.r)
	if err != nil {
		return err
	}

	switch {

	// positive or negative fix bytes
	case b <= 0x7f:
		return m.produceInt(s, msgpackInt{typ: intTypeFixedUint, val: int64(b)})
	case b >= 0xe0:
		return m.produceInt(s, msgpackInt{typ: intTypeFixedInt, val: int64(b)})

	// fix length string
	case b >= 0xa0 && b <= 0xbf:
		return m.decodeString(s, msgpackInt{typ: intTypeUint8, val: int64(b & byte(0x1f))})
	}

	return nil
}
