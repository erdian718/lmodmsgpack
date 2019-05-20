local test = {}
local math = require 'math'
local msgpack = require 'msgpack'

local function sametable(x, y)
	for k, v in pairs(x) do
		if y[k] ~= v then
			return false
		end
	end
	for k, v in pairs(y) do
		if x[k] ~= v then
			return false
		end
	end
	return true
end

function test.table()
	local x = {}
	assert(sametable(x, msgpack.decode(msgpack.encode(x))))

	local x = {'ABC', 123, true, false, 12.34}
	assert(sametable(x, msgpack.decode(msgpack.encode(x))))

	local x = {'ABC', 123, true, nil, false, 12.34}
	assert(sametable(x, msgpack.decode(msgpack.encode(x))))

	local x = {
		k1 = 'ABC';
		k2 = 123;
		k3 = true;
		k4 = false;
		k5 = 12.34;
	}
	assert(sametable(x, msgpack.decode(msgpack.encode(x))))

	local x = {}
	for i = 1, 10000 do
		x[i] = math.random()
	end
	assert(sametable(x, msgpack.decode(msgpack.encode(x))))

	local x = {}
	for i = 1, 10000 do
		x[math.random(10000)] = i
	end
	assert(sametable(x, msgpack.decode(msgpack.encode(x))))
end

return test
