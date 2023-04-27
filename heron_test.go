package heron_test

import (
	"bytes"
	"testing"

	"github.com/dcilke/heron"
	"github.com/stretchr/testify/require"
)

func TestEmitter(_ *testing.T) {
	file := bytes.NewBufferString("noop test")
	e := heron.New()
	e.Process(file)
}

func TestEmitter_Object(t *testing.T) {
	file := bytes.NewBufferString(`{"foo": "bar"}`)
	out := make([]any, 0, 1)
	e := heron.New(
		heron.WithJSON(func(o any) { out = append(out, o) }),
		heron.WithBytes(func(o []byte) { out = append(out, string(o)) }),
		heron.WithError(func(o error) { out = append(out, o) }),
	)
	e.Process(file)
	require.Equal(t, []any{
		map[string]any{"foo": "bar"},
	}, out)
}

func TestEmitter_Array(t *testing.T) {
	file := bytes.NewBufferString(`["a","b","c"]`)
	out := make([]any, 0, 1)
	e := heron.New(
		heron.WithJSON(func(o any) { out = append(out, o) }),
		heron.WithBytes(func(o []byte) { out = append(out, string(o)) }),
		heron.WithError(func(o error) { out = append(out, o) }),
	)
	e.Process(file)
	require.Equal(t, []any{
		[]any{"a", "b", "c"},
	}, out)
}

func TestEmitter_NDJSON(t *testing.T) {
	file := bytes.NewBuffer([]byte("{\"foo\": \"bar\"}\n{\"bin\": \"baz\"}"))
	out := make([]any, 0, 2)
	e := heron.New(
		heron.WithJSON(func(o any) { out = append(out, o) }),
		heron.WithBytes(func(o []byte) { out = append(out, string(o)) }),
		heron.WithError(func(o error) { out = append(out, o) }),
	)
	e.Process(file)
	require.Equal(t, []any{
		map[string]any{"foo": "bar"},
		map[string]any{"bin": "baz"},
	}, out)
}

func TestEmitter_JSONs(t *testing.T) {
	file := bytes.NewBuffer([]byte("{\"foo\": \"bar\"}{\"foo\": \"bar\"}"))
	out := make([]any, 0, 2)
	e := heron.New(
		heron.WithJSON(func(o any) { out = append(out, o) }),
		heron.WithBytes(func(o []byte) { out = append(out, string(o)) }),
		heron.WithError(func(o error) { out = append(out, o) }),
	)
	e.Process(file)
	require.Equal(t, []any{
		map[string]any{"foo": "bar"},
		map[string]any{"foo": "bar"},
	}, out)
}

func TestEmitter_Bytes(t *testing.T) {
	file := bytes.NewBufferString("not json")
	out := make([]any, 0, 2)
	e := heron.New(
		heron.WithJSON(func(o any) { out = append(out, o) }),
		heron.WithBytes(func(o []byte) { out = append(out, string(o)) }),
		heron.WithError(func(o error) { out = append(out, o) }),
	)
	e.Process(file)
	require.Equal(t, []any{"not json"}, out)
}

func TestEmitter_Order(t *testing.T) {
	file := bytes.NewBufferString(`string {"foo": "bar"}`)
	out := make([]any, 0, 2)
	e := heron.New(
		heron.WithJSON(func(o any) { out = append(out, o) }),
		heron.WithBytes(func(o []byte) { out = append(out, string(o)) }),
		heron.WithError(func(o error) { out = append(out, o) }),
	)
	e.Process(file)
	require.Equal(t, []any{
		"string ",
		map[string]any{"foo": "bar"},
	}, out)
}

func TestEmitter_BufLimit(t *testing.T) {
	file := bytes.NewBufferString("{\"foo\": \"bar\"}a\nbcd")
	out := make([]any, 0, 3)
	e := heron.New(
		heron.WithBufSize(2),
		heron.WithJSON(func(o any) { out = append(out, o) }),
		heron.WithBytes(func(o []byte) { out = append(out, string(o)) }),
		heron.WithError(func(o error) { out = append(out, o) }),
	)
	e.Process(file)
	require.Equal(t, []any{
		map[string]any{"foo": "bar"},
		"a",
		"bc",
		"d",
	}, out)
}

func TestEmitter_NoBuf(t *testing.T) {
	file := bytes.NewBufferString("{\"foo\": \"bar\"}a\nbcd")
	out := make([]any, 0, 3)
	e := heron.New(
		heron.WithBufSize(0),
		heron.WithJSON(func(o any) { out = append(out, o) }),
		heron.WithBytes(func(o []byte) { out = append(out, string(o)) }),
		heron.WithError(func(o error) { out = append(out, o) }),
	)
	e.Process(file)
	require.Equal(t, []any{
		map[string]any{"foo": "bar"},
	}, out)
}
