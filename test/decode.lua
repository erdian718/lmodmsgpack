local test = {}
local string = require 'string'
local msgpack = require 'msgpack'

function test.base()
	assert(msgpack.decode(string.char(192)) == nil)
	assert(msgpack.decode(string.char(194)) == false)
	assert(msgpack.decode(string.char(195)) == true)
end

function test.float()
	assert(msgpack.decode(string.char(203, 84, 178, 73, 173, 37, 148, 195, 125)) == 1e100)
end

function test.integer()
	assert(msgpack.decode(string.char(0)) == 0)
	assert(msgpack.decode(string.char(127)) == 127)
	assert(msgpack.decode(string.char(204, 128)) == 128)
	assert(msgpack.decode(string.char(204, 255)) == 255)
	assert(msgpack.decode(string.char(205, 1, 0)) == 256)
	assert(msgpack.decode(string.char(205, 255, 255)) == 65535)
	assert(msgpack.decode(string.char(206, 0, 1, 0, 0)) == 65536)
	assert(msgpack.decode(string.char(206, 255, 255, 255, 255)) == 4294967295)
	assert(msgpack.decode(string.char(207, 0, 0, 0, 1, 0, 0, 0, 0)) == 4294967296)

	assert(msgpack.decode(string.char(255)) == -1)
	assert(msgpack.decode(string.char(224)) == -32)
	assert(msgpack.decode(string.char(208, 223)) == -33)
	assert(msgpack.decode(string.char(208, 128)) == -128)
	assert(msgpack.decode(string.char(209, 255, 127)) == -129)
	assert(msgpack.decode(string.char(209, 128, 0)) == -32768)
	assert(msgpack.decode(string.char(210, 255, 255, 127, 255)) == -32769)
	assert(msgpack.decode(string.char(210, 128, 0, 0, 0)) == -2147483648)
	assert(msgpack.decode(string.char(211, 255, 255, 255, 255, 127, 255, 255, 255)) == -2147483649)
end

function test.string()
	assert(msgpack.decode(string.char(160)) == '')

	local s = string.rep('*', 31)
	assert(msgpack.decode(string.char(191)..s) == s)

	local s = string.rep('*', 32)
	assert(msgpack.decode(string.char(217, 32)..s) == s)

	local s = string.rep('*', 255)
	assert(msgpack.decode(string.char(217, 255)..s) == s)

	local s = string.rep('*', 256)
	assert(msgpack.decode(string.char(218, 1, 0)..s) == s)

	local s = string.rep('*', 65535)
	assert(msgpack.decode(string.char(218, 255, 255)..s) == s)

	local s = string.rep('*', 65536)
	assert(msgpack.decode(string.char(219, 0, 1, 0, 0)..s) == s)
end

return test
