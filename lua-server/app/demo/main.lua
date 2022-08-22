package.path = package.path .. ";../../app/?.lua;../../app/demo/?.lua"

_G.log_tree = require "tree"
local SimpleServer = require "SimpleServer"

local lua_server = SimpleServer.create()
lua_server.Send("hello\n")