local test = {}
local math = require 'math'
local msgpack = require 'msgpack'

local function equal(x, y)
	if x == y then
		return true
	end
	if type(x) ~= 'table' or type(y) ~= 'table' then
		return false
	end
	for k, v in pairs(x) do
		if not equal(v, y[k]) then
			return false
		end
	end
	for k, v in pairs(y) do
		if not equal(v, x[k]) then
			return false
		end
	end
	return true
end

function test.base()
	local x = {}
	assert(equal(x, msgpack.decode(msgpack.encode(x))))

	local x = {'ABC', 123, true, false, 12.34}
	assert(equal(x, msgpack.decode(msgpack.encode(x))))

	local x = {'ABC', 123, true, nil, false, 12.34}
	assert(equal(x, msgpack.decode(msgpack.encode(x))))

	local x = {
		k1 = 'ABC';
		k2 = 123;
		k3 = true;
		k4 = false;
		k5 = 12.34;
		k6 = {
			kk1 = 'EFG';
			kk2 = 456;
			kk3 = true;
			kk4 = false;
			kk5 = 56.78;
			kk6 = {'A', 'B', 'C', 'E', 'F', 'G', {1, 2, 3, 4, 5, 6, 7, 8, 9}};
		};
	}
	assert(equal(x, msgpack.decode(msgpack.encode(x))))

	local x = {}
	for i = 1, 10000 do
		x[i] = math.random()
	end
	assert(equal(x, msgpack.decode(msgpack.encode(x))))

	local x = {}
	for i = 1, 10000 do
		x[math.random(10000)] = i
	end
	assert(equal(x, msgpack.decode(msgpack.encode(x))))
end

return test
