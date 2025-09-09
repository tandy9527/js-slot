-- bet.lua
-- KEYS[1] = 用户id
-- ARGV[1] = 下注金额
-- 返回值: -1 下注金额小于0, -2 余额不足,  {下注后余额, 下注前余额}

local bet = tonumber(ARGV[1])
if not bet or bet <= 0 then
    return -1
end

local user_key = KEYS[1]

-- 获取余额
local balance = tonumber(redis.call("HGET", user_key, "balance")) or 0
if balance < bet then
    return -2
end

-- 扣掉余额
local new_balance = redis.call("HINCRBY", user_key, "balance", -bet)

-- 返回 {下注后余额, 下注前余额}
return new_balance
