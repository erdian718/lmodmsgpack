# lmodmsgpack

MessagePack for [Lua](https://github.com/ofunc/lua).

## Usage

```go
package main

import (
	"ofunc/lmodmsgpack"
	"ofunc/lua/util"
)

func main() {
	l := util.NewState()
	l.Preload("msgpack", lmodmsgpack.Open)
	util.Run(l, "main.lua")
}
```

```lua
local msgpack = require 'msgpack'

local x = msgpack.encode(v)
local y = msgpack.decode(x)
```

## Dependencies

* [ofunc/lua](https://github.com/ofunc/lua)

## Documentation

### msgpack.encode([w, ]v1[, v2, ...])

Encodes the values.
If writer `w` is provided, the encoded data will be writed to `w`.
Otherwise, the encoded data will be returned.

### msgpack.decode(x[, i])

Decodes the fisrt value.
`x` can be a reader or a string.
`i` is the start position of the string, default value is `1`.
