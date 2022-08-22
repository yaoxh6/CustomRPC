package.path = package.path .. ";../../app/?.lua;../../app/demo/?.lua"

_G.log_tree = require "tree"
_G.s2s = {}
local SimpleServer = require "SimpleServer"

function s2s.hello()
    print("hello call back!")
end

local lua_server = SimpleServer.create()
lua_server.on_call_with_handle = function(msg, ...)
    log_tree("msg", {msg, ...})
    local proc = s2s[msg]
    if not proc then
        print("function ", msg, " is not exist")
        return
    end
    local ok, err = xpcall(proc, debug.traceback, ...)
    if not ok then
        print("function ", msg, " execution failed")
        return
    end
end
lua_server.Send("hello")
