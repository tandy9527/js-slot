--- zincrby.lua
--- KEYS[1] = redis key
--- ARGV[1] = member
--- ARGV[2] = score
--- 返回值: score
--- zincrBy 增加分数,score 最低为0

local key = KEYS[1]
local member = KEYS[2]
local addscore = tonumber(ARGV[1])

local score = redis.call('ZSCORE', key, member)
if not score then
  score = 0
else
  score = tonumber(score)
end

score = score + addscore
if score < 0 then
  score = 0
end

redis.call('ZINCRBY', key, score, member)
return score