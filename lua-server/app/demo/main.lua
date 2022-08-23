package.path = package.path .. ";../../app/?.lua;../../app/demo/?.lua"

_G.log_tree = require "tree"
_G.s2s = {}
local SimpleServer = require "SimpleServer"

function s2s.hello_reply(param)
    print("hello call back! param:", param)
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
-- 调用go服务的SayHello函数, 参数是lua_server_test
-- 回调函数是hello_reply, 即go服务回包是s2s.hello_reply的参数
-- 默认第二个参数是回调函数, 是否有回调函数用go侧的pb文件区分
lua_server.Send("SayHello", "hello_reply", "lua_server_test")
