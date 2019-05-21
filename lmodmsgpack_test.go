package lmodmsgpack_test

import (
	"fmt"
	"testing"

	"ofunc/lmodmsgpack"
	"ofunc/lua"
	"ofunc/lua/util"
)

type test struct{}

func TestMain(m *testing.M) {
	l := util.NewState()
	l.Preload("msgpack", lmodmsgpack.Open)

	l.Push(test{})
	l.NewTable(0, 2)
	l.Push("__index")
	l.Push(func(l *lua.State) int {
		l.Push(2 * l.ToInteger(2))
		return 1
	})
	l.SetTableRaw(-3)
	l.Push("__len")
	l.Push(func(l *lua.State) int {
		l.Push(10000)
		return 1
	})
	l.SetTableRaw(-3)
	l.SetMetaTable(-2)
	l.SetGlobal("testlen")

	l.Push(test{})
	l.NewTable(0, 2)
	l.Push("__pairs")
	l.Push(func(l *lua.State) int {
		i := 0
		l.Push(func(l *lua.State) int {
			if i >= 10000 {
				return 0
			}
			l.Push(i)
			l.Push(2 * i)
			i++
			return 2
		})
		l.PushIndex(1)
		l.Push(nil)
		return 3
	})
	l.SetTableRaw(-3)
	l.SetMetaTable(-2)
	l.SetGlobal("testpairs")

	if err := util.Test(l, "test"); err != nil {
		fmt.Println("error:", err)
	}
}
