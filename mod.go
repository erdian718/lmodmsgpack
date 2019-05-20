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

// MessagePack for Lua.
package lmodmsgpack

import (
	"bytes"
	"io"
	"strings"

	"ofunc/lua"
)

// Open opens the module.
func Open(l *lua.State) int {
	l.NewTable(0, 4)

	l.Push("version")
	l.Push("0.0.1")
	l.SetTableRaw(-3)

	l.Push("encode")
	l.Push(lEncode)
	l.SetTableRaw(-3)

	l.Push("decode")
	l.Push(lDecode)
	l.SetTableRaw(-3)

	return 1
}

func lEncode(l *lua.State) int {
	s, buf := 2, new(bytes.Buffer)
	w, ok := l.GetRaw(1).(io.Writer)
	if !ok {
		s, w = 1, buf
	}

	var m int64
	var err error
	n := l.AbsIndex(-1)
	for i := s; i <= n && err == nil; i++ {
		k, e := EncodeWriter(l, i, w)
		m, err = m+k, e
		if err != nil {
			break
		}
	}
	if ok {
		l.Push(m)
		if err == nil {
			return 1
		} else {
			l.Push(err.Error())
			return 2
		}
	} else {
		l.Push(buf.String())
		return 1
	}
}

func lDecode(l *lua.State) int {
	var e error
	var r io.Reader
	var i int64
	if l.TypeOf(1) == lua.TypeUserData {
		r = toReader(l, 1)
	} else {
		s := l.ToString(1)
		i = l.OptInteger(2, 1)
		if i < 0 {
			i += int64(len(s)) + 1
		}
		r = strings.NewReader(s[i-1:])
	}
	k, e := DecodeReader(l, r)
	if e == nil {
		l.Push(nil)
		l.Push(k + i)
	} else {
		l.Push(nil)
		l.Push(e.Error())
		l.Push(k + i)
	}
	return 3
}
