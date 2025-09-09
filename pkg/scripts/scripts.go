package scripts

//----------------------- 用redis执行lua脚本保证原子性 -----------------------

import _ "embed"

//go:embed recharge.lua
var RechargeLua string

//go:embed bet.lua
var BetLua string
