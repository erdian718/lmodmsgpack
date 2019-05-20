local test = {}
local string = require 'string'
local msgpack = require 'msgpack'

function test.base()
	assert(msgpack.encode(nil) == string.char(192))
	assert(msgpack.encode(false) == string.char(194))
	assert(msgpack.encode(true) == string.char(195))
end

function test.float()
	assert(msgpack.encode(1e100) == string.char(203, 84, 178, 73, 173, 37, 148, 195, 125))
end

function test.integer()
	assert(msgpack.encode(0) == string.char(0))
	assert(msgpack.encode(127) == string.char(127))
	assert(msgpack.encode(128) == string.char(204, 128))
	assert(msgpack.encode(255) == string.char(204, 255))
	assert(msgpack.encode(256) == string.char(205, 1, 0))
	assert(msgpack.encode(65535) == string.char(205, 255, 255))
	assert(msgpack.encode(65536) == string.char(206, 0, 1, 0, 0))
	assert(msgpack.encode(4294967295) == string.char(206, 255, 255, 255, 255))
	assert(msgpack.encode(4294967296) == string.char(207, 0, 0, 0, 1, 0, 0, 0, 0))

	assert(msgpack.encode(-1) == string.char(255))
	assert(msgpack.encode(-32) == string.char(224))
	assert(msgpack.encode(-33) == string.char(208, 223))
	assert(msgpack.encode(-128) == string.char(208, 128))
	assert(msgpack.encode(-129) == string.char(209, 255, 127))
	assert(msgpack.encode(-32768) == string.char(209, 128, 0))
	assert(msgpack.encode(-32769) == string.char(210, 255, 255, 127, 255))
	assert(msgpack.encode(-2147483648) == string.char(210, 128, 0, 0, 0))
	assert(msgpack.encode(-2147483649) == string.char(211, 255, 255, 255, 255, 127, 255, 255, 255))
end

function test.string()
	assert(msgpack.encode('') == string.char(160))

	local s = string.rep('*', 31)
	assert(msgpack.encode(s) == string.char(191)..s)

	local s = string.rep('*', 32)
	assert(msgpack.encode(s) == string.char(217, 32)..s)

	local s = string.rep('*', 255)
	assert(msgpack.encode(s) == string.char(217, 255)..s)

	local s = string.rep('*', 256)
	assert(msgpack.encode(s) == string.char(218, 1, 0)..s)

	local s = string.rep('*', 65535)
	assert(msgpack.encode(s) == string.char(218, 255, 255)..s)

	local s = string.rep('*', 65536)
	assert(msgpack.encode(s) == string.char(219, 0, 1, 0, 0)..s)
end

return test
