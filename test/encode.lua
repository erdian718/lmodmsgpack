local test = {}
local string = require 'string'
local msgpack = require 'msgpack'

function test.base()
	assert(msgpack.encode(nil) == string.char(192))
	assert(msgpack.encode(false) == string.char(194))
	assert(msgpack.encode(true) == string.char(195))
end

function test.float()
	-- TODO float32
	assert(msgpack.encode(1e100) == string.char(203, 84, 178, 73, 173, 37, 148, 195, 125))
end

function test.integer()
	assert(msgpack.encode(0) == string.char(0))
	assert(msgpack.encode(127) == string.char(127))
	assert(msgpack.encode(-1) == string.char(255))
	assert(msgpack.encode(-32) == string.char(224))

	assert(msgpack.encode(-33) == string.char(208, 223))
	assert(msgpack.encode(-128) == string.char(208, 128))
	-- TODO uint
end

function test.string()
	assert(msgpack.encode('') == string.char(160))

	local s = string.rep('*', 31)
	local xs = msgpack.encode(s)
	assert(string.byte(xs) == 191)
	assert(string.sub(xs, 2) == s)

	local s = string.rep('*', 32)
	local xs = msgpack.encode(s)
	assert(string.byte(xs, 1) == 217)
	assert(string.byte(xs, 2) == 32)
	assert(string.sub(xs, 3) == s)

	local s = string.rep('*', 255)
	local xs = msgpack.encode(s)
	assert(string.byte(xs, 1) == 217)
	assert(string.byte(xs, 2) == 255)
	assert(string.sub(xs, 3) == s)

	local s = string.rep('*', 256)
	local xs = msgpack.encode(s)
	assert(string.byte(xs, 1) == 218)
	assert(string.byte(xs, 2) == 1)
	assert(string.byte(xs, 3) == 0)
	assert(string.sub(xs, 4) == s)

	local s = string.rep('*', 65535)
	local xs = msgpack.encode(s)
	assert(string.byte(xs, 1) == 218)
	assert(string.byte(xs, 2) == 255)
	assert(string.byte(xs, 3) == 255)
	assert(string.sub(xs, 4) == s)

	local s = string.rep('*', 65536)
	local xs = msgpack.encode(s)
	assert(string.byte(xs, 1) == 219)
	assert(string.byte(xs, 2) == 0)
	assert(string.byte(xs, 3) == 1)
	assert(string.byte(xs, 4) == 0)
	assert(string.byte(xs, 5) == 0)
	assert(string.sub(xs, 6) == s)
end

return test
