// package heron provides a streaming JSON parser that emits JSON
// objects and arrays as they are parsed. Non-JSON lines and errors
// are emitted separately.
package heron

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/dcilke/goj"
)

const (
	// The default buffer size used to accumulate non json bytes.
	DefaultBufSize = 512

	newLine = 10
)

type Heron struct {
	buf     *bytes.Buffer
	bufSize int
	json    func(any)
	bytes   func([]byte)
	err     func(error)
}

type Option func(*Heron)

// WithBufSize sets the buffer size used to accumulate non json bytes.
func WithBufSize(size int) Option {
	return func(p *Heron) {
		p.bufSize = size
	}
}

// WithJSON sets the function to call when a JSON object or array is parsed.
func WithJSON(f func(any)) Option {
	return func(p *Heron) {
		p.json = f
	}
}

// WithBytes sets the function to call when a non JSON line is parsed. This
// function is called when either a newline is encountered or the buffer is
// full. Setting to 0 disables internal buffering and emitting of non JSON
// bytes.
func WithBytes(f func([]byte)) Option {
	return func(p *Heron) {
		p.bytes = f
	}
}

// WithError sets the function to call when an error is encountered when
// processing the stream.
func WithError(f func(error)) Option {
	return func(p *Heron) {
		p.err = f
	}
}

// New creates a new Heron parser with the given options.
func New(opts ...Option) *Heron {
	p := &Heron{
		bufSize: DefaultBufSize,
		json:    func(a any) {},
		bytes:   func(s []byte) {},
		err:     func(e error) {},
	}

	for _, opt := range opts {
		opt(p)
	}

	if p.bufSize > 0 {
		p.buf = bytes.NewBuffer(make([]byte, 0, p.bufSize))
	}

	return p
}

// Process reads the given stream and emits values.
func (p *Heron) Process(stream io.Reader) {
	decoder := goj.NewDecoder(stream)
	decoder.UseNumber()

	for {
		// we will only emit objects or arrays, so if the current character is not { or [, process as byte stream
		if c, _ := decoder.Peek(); goj.IsBegin(c) {
			var msg any
			err := decoder.Decode(&msg)
			if err == nil {
				switch k := reflect.TypeOf(msg).Kind(); k {
				case reflect.Map, reflect.Array, reflect.Slice:
					p.Flush()
					p.json(msg)
				default:
					p.err(fmt.Errorf("unexpected decoded line type %q", k))
				}
				continue
			}
		}

		if err := p.processByte(decoder); err == io.EOF {
			return
		}
	}
}

// BufSize returns the buffer size used to accumulate non json bytes.
func (p *Heron) BufSize() int {
	return p.bufSize
}

// Flush flushes the buffer.
func (p *Heron) Flush() {
	if p.buf == nil {
		return
	}
	if p.buf.Len() > 0 {
		p.bytes(p.buf.Bytes())
		p.buf.Reset()
	}
}

func (p *Heron) processByte(decoder *goj.Decoder) error {
	b, err := decoder.ReadByte()
	if err == io.EOF {
		p.Flush()
		return err
	}

	p.push(b)
	if b == newLine {
		p.Flush()
	} else {
		p.flushFull()
	}
	return nil
}

func (p *Heron) push(b byte) {
	if p.buf == nil {
		return
	}
	p.buf.WriteByte(b)
}

func (p *Heron) flushFull() {
	if p.buf == nil {
		return
	}
	if p.buf.Len() >= p.bufSize {
		p.Flush()
	}
}
