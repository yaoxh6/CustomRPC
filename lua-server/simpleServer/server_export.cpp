#include "luna.h"
#include "server_export.h"
#include "server.h"
#include <cassert>

static int lua_server(lua_State* L)
{
	auto m_server = new SimpleServer();
	if (!m_server->Init() || !m_server->Connect()) {
		return 0;
	}
	lua_push_object(L, m_server);
	return 1;
}

int luaopen_server(lua_State* L)
{
	luaL_checkversion(L);
	luaL_Reg l[] = {
		{ "create", lua_server},
		{ NULL, NULL },
	};
	luaL_newlib(L, l);
	return 1;
}
