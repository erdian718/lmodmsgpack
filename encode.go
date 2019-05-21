/*
Copyright 2019 by ofunc

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

package lmodmsgpack

import (
	"bytes"
	"io"
	"math"

	"ofunc/lua"
)

// EncodeBytes encodes the Lua value to bytes.
func EncodeBytes(l *lua.State, idx int) []byte {
	idx = l.AbsIndex(idx)
	switch l.TypeOf(idx) {
	case lua.TypeString:
		return encodeString(l.ToString(idx))
	case lua.TypeNumber:
		if v, e := l.TryInteger(idx); e == nil {
			return encodeInteger(v)
		} else {
			return encodeFloat(l.ToFloat(idx))
		}
	case lua.TypeBoolean:
		if l.ToBoolean(idx) {
			return []byte{0xc3}
		} else {
			return []byte{0xc2}
		}
	case lua.TypeTable:
		buf := new(bytes.Buffer)
		encodeTable(l, idx, buf)
		return buf.Bytes()
	}
	if l.GetMetaField(idx, "__len") != lua.TypeNil {
		buf := new(bytes.Buffer)
		encodeArray(l, idx, buf)
		return buf.Bytes()
	}
	if l.GetMetaField(idx, "__pairs") != lua.TypeNil {
		buf := new(bytes.Buffer)
		encodePairs(l, idx, buf)
		return buf.Bytes()
	}
	return []byte{0xc0}
}

// EncodeWriter encodes the Lua value to the writer.
func EncodeWriter(l *lua.State, idx int, w io.Writer) (int64, error) {
	idx = l.AbsIndex(idx)
	switch l.TypeOf(idx) {
	case lua.TypeString:
		k, e := w.Write(encodeString(l.ToString(idx)))
		return int64(k), e
	case lua.TypeNumber:
		if v, e := l.TryInteger(idx); e == nil {
			k, e := w.Write(encodeInteger(v))
			return int64(k), e
		} else {
			k, e := w.Write(encodeFloat(l.ToFloat(idx)))
			return int64(k), e
		}
	case lua.TypeBoolean:
		if l.ToBoolean(idx) {
			k, e := w.Write([]byte{0xc3})
			return int64(k), e
		} else {
			k, e := w.Write([]byte{0xc2})
			return int64(k), e
		}
	case lua.TypeTable:
		return encodeTable(l, idx, w)
	}
	if l.GetMetaField(idx, "__len") != lua.TypeNil {
		return encodeArray(l, idx, w)
	}
	if l.GetMetaField(idx, "__pairs") != lua.TypeNil {
		return encodePairs(l, idx, w)
	}
	k, e := w.Write([]byte{0xc0})
	return int64(k), e
}

func encodeString(v string) []byte {
	var xs []byte
	n := int64(len(v))
	switch {
	case n < 1<<5:
		xs = make([]byte, n+1)
		xs[0] = 0xa0 | byte(n)
		copy(xs[1:], v)
	case n < 1<<8:
		xs = make([]byte, n+2)
		xs[0] = 0xd9
		xs[1] = byte(n)
		copy(xs[2:], v)
	case n < 1<<16:
		xs = make([]byte, n+3)
		xs[0] = 0xda
		endian.PutUint16(xs[1:], uint16(n))
		copy(xs[3:], v)
	case n < 1<<32:
		xs = make([]byte, n+5)
		xs[0] = 0xdb
		endian.PutUint32(xs[1:], uint32(n))
		copy(xs[5:], v)
	default:
		panic(ErrLengthOverflow)
	}
	return xs
}

func encodeFloat(v float64) []byte {
	var xs []byte
	xs = make([]byte, 9)
	xs[0] = 0xcb
	endian.PutUint64(xs[1:], math.Float64bits(v))
	return xs
}

func encodeInteger(v int64) []byte {
	var xs []byte
	if v >= 0 {
		switch {
		case v < 128:
			xs = make([]byte, 1)
			xs[0] = byte(v)
		case int64(uint8(v)) == v:
			xs = make([]byte, 2)
			xs[0] = 0xcc
			xs[1] = byte(v)
		case int64(uint16(v)) == v:
			xs = make([]byte, 3)
			xs[0] = 0xcd
			endian.PutUint16(xs[1:], uint16(v))
		case int64(uint32(v)) == v:
			xs = make([]byte, 5)
			xs[0] = 0xce
			endian.PutUint32(xs[1:], uint32(v))
		default:
			xs = make([]byte, 9)
			xs[0] = 0xcf
			endian.PutUint64(xs[1:], uint64(v))
		}
	} else {
		switch {
		case v >= -32:
			xs = make([]byte, 1)
			xs[0] = byte(v)
		case int64(int8(v)) == v:
			xs = make([]byte, 2)
			xs[0] = 0xd0
			xs[1] = byte(v)
		case int64(int16(v)) == v:
			xs = make([]byte, 3)
			xs[0] = 0xd1
			endian.PutUint16(xs[1:], uint16(v))
		case int64(int32(v)) == v:
			xs = make([]byte, 5)
			xs[0] = 0xd2
			endian.PutUint32(xs[1:], uint32(v))
		default:
			xs = make([]byte, 9)
			xs[0] = 0xd3
			endian.PutUint64(xs[1:], uint64(v))
		}
	}
	return xs
}

func encodeTable(l *lua.State, idx int, w io.Writer) (m int64, err error) {
	var xs []byte
	n := l.Count(idx)
	switch {
	case n < 1<<4:
		xs = []byte{0x80 | byte(n)}
	case n < 1<<16:
		xs = make([]byte, 3)
		xs[0] = 0xde
		endian.PutUint16(xs[1:], uint16(n))
	case n < 1<<32:
		xs = make([]byte, 5)
		xs[0] = 0xdf
		endian.PutUint32(xs[1:], uint32(n))
	default:
		panic(ErrLengthOverflow)
	}
	k, e := w.Write(xs)
	m, err = int64(k), e
	if err != nil {
		return
	}
	l.ForEachRaw(idx, func() bool {
		k, e := EncodeWriter(l, -2, w)
		m, err = m+k, e
		if err == nil {
			k, e := EncodeWriter(l, -1, w)
			m, err = m+k, e
		}
		return err == nil
	})
	return
}

func encodeArray(l *lua.State, idx int, w io.Writer) (m int64, err error) {
	var xs []byte
	l.PushIndex(idx)
	l.Call(1, 1)
	n := l.ToInteger(-1)
	l.Pop(1)
	switch {
	case n < 1<<4:
		xs = []byte{0x90 | byte(n)}
	case n < 1<<16:
		xs = make([]byte, 3)
		xs[0] = 0xdc
		endian.PutUint16(xs[1:], uint16(n))
	case n < 1<<32:
		xs = make([]byte, 5)
		xs[0] = 0xdd
		endian.PutUint32(xs[1:], uint32(n))
	default:
		panic(ErrLengthOverflow)
	}
	k, e := w.Write(xs)
	m, err = int64(k), e
	if err != nil {
		return
	}
	for i := int64(1); i <= n; i++ {
		l.Push(i)
		l.GetTable(idx)
		k, e := EncodeWriter(l, -1, w)
		m, err = m+k, e
		l.Pop(1)
		if err != nil {
			break
		}
	}
	return
}

func encodePairs(l *lua.State, idx int, w io.Writer) (m int64, err error) {
	var n int64
	var xs, buf []byte
	l.Pop(1)
	l.ForEach(idx, func() bool {
		buf = append(buf, EncodeBytes(l, -2)...)
		buf = append(buf, EncodeBytes(l, -1)...)
		n++
		return true
	})
	switch {
	case n < 1<<4:
		xs = []byte{0x80 | byte(n)}
	case n < 1<<16:
		xs = make([]byte, 3)
		xs[0] = 0xde
		endian.PutUint16(xs[1:], uint16(n))
	case n < 1<<32:
		xs = make([]byte, 5)
		xs[0] = 0xdf
		endian.PutUint32(xs[1:], uint32(n))
	default:
		panic(ErrLengthOverflow)
	}
	k, e := w.Write(xs)
	m, err = int64(k), e
	if err != nil {
		return
	}
	k, e = w.Write(buf)
	m, err = m+int64(k), e
	return
}
