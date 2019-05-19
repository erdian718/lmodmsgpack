package lmodmsgpack_test

import (
	"fmt"
	"testing"

	"ofunc/lmodmsgpack"
	"ofunc/lua/util"
)

func TestMain(m *testing.M) {
	l := util.NewState()
	l.Preload("msgpack", lmodmsgpack.Open)
	if err := util.Test(l, "test"); err != nil {
		fmt.Println("error:", err)
	}
}
