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
	-- TODO
end

function test.table()
	assert(msgpack.encode({}) == string.char(128))
	assert(msgpack.encode({
		a = 1;
		b = 2;
		c = 3;
	}) == string.char(131, 161, 97, 1, 161, 98, 2, 161, 99, 3))
end

function test.userdata()
	-- TODO
end

return test
