package.path = package.path .. ";../../app/?.lua;../../app/demo/?.lua"

_G.log_tree = require "tree"
_G.s2s = {}
local SimpleServer = require "SimpleServer"
local lua_server = SimpleServer.create()

function s2s.hello_reply_reply(...)
    log_tree("[hello_reply_reply] param:", {...})
end

function s2s.hello_reply(...)
    log_tree("[hello_reply] param:", {...})
    lua_server.Send("SayHello2", "hello_reply_reply", "param2", 12345)
end


lua_server.on_call_with_handle = function(msg, ...)
    log_tree("msg", {msg, ...})
    local proc = s2s[msg]
    if not proc then
        print("function ", msg, " is not exist")
        return
    end
    local ok, err = xpcall(proc, debug.traceback, ...)
    if not ok then
        print("function ", msg, " execution failed err", err)
        return
    end
end
-- 调用go服务的SayHello函数, 参数是lua_server_test
-- 回调函数是hello_reply, 即go服务回包是s2s.hello_reply的参数
lua_server.Send("SayHello", "hello_reply", "param1")
