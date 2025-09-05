-- recharge.lua
-- KEYS[1] = 用户id
-- ARGV[1] = 更新金额
-- 返回值: {new_balance, old_balance}, 修改失败返回 {-1, -1}

local amount = tonumber(ARGV[1])
if not amount or amount <= 0 then
    return {-1, -1} -- 参数错误
end

local user_key = KEYS[1]

-- 充值前账户金额
local old_balance = redis.call("HGET", user_key, "balance")
if not old_balance then
    old_balance = 0
else
    old_balance = tonumber(old_balance)
end

-- 增加余额
local new_balance = redis.call("HINCRBY", user_key, "balance", amount)

-- 返回 {充值后余额, 充值前余额}
return {new_balance, old_balance}