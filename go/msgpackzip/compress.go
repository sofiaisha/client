package msgpackzip

import (
	"bytes"
	"fmt"
	"sort"
)

type compressor struct {
	input *bytes.Buffer
}

func newCompressor(b []byte) *compressor {
	return &compressor{input: bytes.NewBuffer(b)}
}

func Compress(input []byte) (output []byte, err error) {
	return newCompressor(input).run()
}

func (c *compressor) run() (output []byte, err error) {

	freqs, err := c.collectFrequencies()
	if err != nil {
		return nil, err
	}
	keys, err := c.sortIntoKeys(freqs)
	if err != nil {
		return nil, err
	}
	output, err = c.output(keys)
	return output, err
}

func (c *compressor) collectFrequencies() (ret map[interface{}]int, err error) {

	ret = make(map[interface{}]int)
	f := func(d decodeStack) (decodeStack, error) {
		d.hooks.stringHook = func(l msgpackInt, s string) error {
			ret[s]++
			return nil
		}
		d.hooks.intHook = func(l msgpackInt) error {
			i, err := l.toInt64()
			if err != nil {
				return err
			}
			ret[i]++
			return nil
		}
		d.hooks.fallthroughHook = func(i interface{}, s string) error {
			return fmt.Errorf("bad map key (type %T)", i)
		}
		return d, nil
	}

	err = newMsgpackDecoder(c.input).run(msgpackDecoderHooks{mapKeyHook: f})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *compressor) sortIntoKeys(freqs map[interface{}]int) (keys map[interface{}]int, err error) {
	type tuple struct {
		key  interface{}
		freq int
	}
	tuples := make([]tuple, len(freqs))
	var i int
	for k, v := range freqs {
		tuples[i] = tuple{k, v}
		i++
	}
	sort.SliceStable(tuples, func(i, j int) bool { return tuples[i].freq > tuples[j].freq })

	ret := make(map[interface{}]int, len(freqs))
	for i, tup := range tuples {
		ret[tup.key] = i
	}
	return ret, nil
}

func (c *compressor) output(keys map[interface{}]int) (output []byte, err error) {
	return nil, nil
}
