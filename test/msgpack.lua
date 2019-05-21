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

function test.table()
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

function test.len()
	local xs = msgpack.decode(msgpack.encode(testlen))
	assert(#xs == 10000)
	for i = 1, 10000 do
		assert(xs[i] == 2*i)
	end
end

function test.pairs()
	local xs = msgpack.decode(msgpack.encode(testpairs))
	local s = 0
	for k, v in pairs(xs) do
		assert(v == 2 * k)
		s = s + v
	end
	assert(s == 99990000)
end

function test.multi()
	local x1, x2, x3, x4, x5, x6, x7 = 'ABC', 123, {'E', 'F', 'G'}, 12.34, nil, true, false
	local s = msgpack.encode(x1, x2, x3, x4, x5, x6, x7)

	local y1, err, k = msgpack.decode(s)
	assert(err == nil)
	assert(equal(x1, y1))

	local y2, err, k = msgpack.decode(s, k)
	assert(err == nil)
	assert(equal(x2, y2))

	local y3, err, k = msgpack.decode(s, k)
	assert(err == nil)
	assert(equal(x3, y3))

	local y4, err, k = msgpack.decode(s, k)
	assert(err == nil)
	assert(equal(x4, y4))

	local y5, err, k = msgpack.decode(s, k)
	assert(err == nil)
	assert(equal(x5, y5))

	local y6, err, k = msgpack.decode(s, k)
	assert(err == nil)
	assert(equal(x6, y6))

	local y7, err, k = msgpack.decode(s, k)
	assert(err == nil)
	assert(equal(x7, y7))
end

return test
