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

// DecodeBytes decodes the Lua value from the bytes.
func DecodeBytes(l *lua.State, xs []byte) (int64, error) {
	return DecodeReader(l, bytes.NewReader(xs))
}

// DecodeReader decodes the Lua value from the reader.
func DecodeReader(l *lua.State, r io.Reader) (m int64, err error) {
	var xs []byte
	xs, err = readn(r, 1)
	if err != nil {
		m = int64(len(xs))
		return
	}
	x := xs[0]
	switch {
	case x <= 0x7f:
		l.Push(int64(x))
	case x <= 0x8f:
		m, err = decodeMap(l, r, x)
	case x <= 0x9f:
		m, err = decodeArray(l, r, x)
	case x <= 0xbf:
		m, err = decodeString(l, r, x)
	case x <= 0xc1:
		l.Push(nil)
	case x <= 0xc2:
		l.Push(false)
	case x <= 0xc3:
		l.Push(true)
	case x <= 0xc6:
		m, err = decodeString(l, r, x)
	case x <= 0xc9:
		m, err = decodeExt(l, r, x)
	case x <= 0xcb:
		m, err = decodeFloat(l, r, x)
	case x <= 0xd3:
		m, err = decodeInteger(l, r, x)
	case x <= 0xd8:
		m, err = decodeExt(l, r, x)
	case x <= 0xdb:
		m, err = decodeString(l, r, x)
	case x <= 0xdd:
		m, err = decodeArray(l, r, x)
	case x <= 0xdf:
		m, err = decodeMap(l, r, x)
	case x <= 0xff:
		l.Push(int64(int8(x)))
	}
	m += 1
	return
}

func decodeString(l *lua.State, r io.Reader, x byte) (m int64, err error) {
	n, k := 0, 0
	if x <= 0xbf {
		n = int(x - 0xa0)
	} else {
		switch x {
		case 0xd9, 0xc4:
			n, k, err = declen(r, 1)
		case 0xda, 0xc5:
			n, k, err = declen(r, 2)
		case 0xdb, 0xc6:
			n, k, err = declen(r, 4)
		}
	}
	m = int64(k)
	if err != nil {
		return
	}

	var xs []byte
	xs, err = readn(r, n)
	m += int64(len(xs))
	if err != nil {
		return
	}
	l.Push(string(xs))
	return
}

func decodeFloat(l *lua.State, r io.Reader, x byte) (m int64, err error) {
	if x == 0xca {
		xs, e := readn(r, 4)
		m, err = int64(len(xs)), e
		if err == nil {
			l.Push(math.Float32frombits(endian.Uint32(xs)))
		}
	} else {
		xs, e := readn(r, 8)
		m, err = int64(len(xs)), e
		if err == nil {
			l.Push(math.Float64frombits(endian.Uint64(xs)))
		}
	}
	return
}

func decodeInteger(l *lua.State, r io.Reader, x byte) (m int64, err error) {
	switch x {
	case 0xcc, 0xd0:
		xs, e := readn(r, 1)
		m, err = int64(len(xs)), e
		if err == nil {
			if x == 0xcc {
				l.Push(int64(xs[0]))
			} else {
				l.Push(int64(int8(xs[0])))
			}
		}
	case 0xcd, 0xd1:
		xs, e := readn(r, 2)
		m, err = int64(len(xs)), e
		if err == nil {
			if x == 0xcd {
				l.Push(int64(endian.Uint16(xs)))
			} else {
				l.Push(int64(int16(endian.Uint16(xs))))
			}
		}
	case 0xce, 0xd2:
		xs, e := readn(r, 4)
		m, err = int64(len(xs)), e
		if err == nil {
			if x == 0xce {
				l.Push(int64(endian.Uint32(xs)))
			} else {
				l.Push(int64(int32(endian.Uint32(xs))))
			}
		}
	case 0xcf, 0xd3:
		xs, e := readn(r, 8)
		m, err = int64(len(xs)), e
		if err == nil {
			if x == 0xcf {
				v0 := endian.Uint64(xs)
				v1 := int64(v0)
				if v1 >= 0 {
					l.Push(v1)
				} else {
					l.Push(float64(v0))
				}
			} else {
				l.Push(int64(endian.Uint64(xs)))
			}
		}
	}
	return
}

func decodeArray(l *lua.State, r io.Reader, x byte) (m int64, err error) {
	n, k := 0, 0
	switch x {
	case 0xdc:
		n, k, err = declen(r, 16)
	case 0xdd:
		n, k, err = declen(r, 32)
	default:
		n = int(x - 0x90)
	}
	m = int64(k)
	if err != nil {
		return
	}
	l.NewTable(n, 0)
	for i := 0; i < n; i++ {
		l.Push(i + 1)
		k, e := DecodeReader(l, r)
		m, err = m+k, e
		if err == nil {
			l.SetTableRaw(-3)
		} else {
			l.Pop(2)
			break
		}
	}
	return
}

func decodeMap(l *lua.State, r io.Reader, x byte) (m int64, err error) {
	n, k := 0, 0
	switch x {
	case 0xde:
		n, k, err = declen(r, 16)
	case 0xdf:
		n, k, err = declen(r, 32)
	default:
		n = int(x - 0x80)
	}
	m = int64(k)
	if err != nil {
		return
	}
	l.NewTable(0, n)
	for i := 0; i < n; i++ {
		k, e := DecodeReader(l, r)
		m, err = m+k, e
		if err != nil {
			l.Pop(1)
			break
		}
		k, e = DecodeReader(l, r)
		m, err = m+k, e
		if err != nil {
			l.Pop(2)
			break
		}
		l.SetTableRaw(-3)
	}
	return
}

func decodeExt(l *lua.State, r io.Reader, x byte) (m int64, err error) {
	n, k := 0, 0
	if 0xd4 <= x && x <= 0xd8 {
		n = 1 << (x - 0xd4)
	} else {
		switch x {
		case 0xc7:
			n, k, err = declen(r, 8)
		case 0xc8:
			n, k, err = declen(r, 16)
		case 0xc9:
			n, k, err = declen(r, 32)
		}
	}
	m = int64(k)
	if err != nil {
		return
	}
	var xs []byte
	xs, err = readn(r, n+1)
	m += int64(len(xs))
	if err == nil {
		l.Push(nil)
	}
	return
}
