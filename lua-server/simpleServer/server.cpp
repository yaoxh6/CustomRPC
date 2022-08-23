#include "server.h"
#include "json.h"
using json_t = nlohmann::json;
LUA_EXPORT_CLASS_BEGIN(SimpleServer)
LUA_EXPORT_METHOD(Send)
LUA_EXPORT_CLASS_END()

bool SimpleServer::Init()
{
	WORD sockVersion = MAKEWORD(2, 2);
	WSADATA data;
	if (WSAStartup(sockVersion, &data) != 0)
	{
		return false;
	}
	sclient = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
	if (sclient == INVALID_SOCKET)
	{
		return false;
	}
	ip = "127.0.0.1";
	port = 8888;
	return true;
}

void SimpleServer::Uninit()
{
	closesocket(sclient);
	WSACleanup();
}

bool SimpleServer::Connect()
{
	sockaddr_in serAddr;
	serAddr.sin_family = AF_INET;
	serAddr.sin_port = htons(port);
	serAddr.sin_addr.S_un.S_addr = inet_addr(ip.c_str());
	if (connect(sclient, (sockaddr*)&serAddr, sizeof(serAddr)) == SOCKET_ERROR)
	{
		printf("connect error !");
		closesocket(sclient);
		return false;
	}
	return true;
}

bool EncodeData(lua_State* L, json_t& j, int index) {
	int type = lua_type(L, index);
	switch (type) {
	case LUA_TNIL:
		j[index-1] = "";
		return true;
	case LUA_TNUMBER:
		lua_isinteger(L, index) ? j[index-1] = (lua_tointeger(L, index)) : j[index-1] = (lua_tonumber(L, index));
		return true;
	case LUA_TBOOLEAN:
		j[index-1] = (!!lua_toboolean(L, index));
		return true;
	case LUA_TSTRING:
		j[index-1] = lua_tostring(L, index);
		return true;
	case LUA_TTABLE:
		//暂时不支持table
		return false;
	default:
		break;
	}
	return false;
}

int SimpleServer::Send(lua_State* L)
{
	json_t j;
	int top = lua_gettop(L);
	for (int i = 1; i<=top; i++)
	{
		if (!EncodeData(L, j, i)) {
			return 0;
		}
	}
	send(sclient, j.dump().c_str(), strlen(j.dump().c_str()), 0);
	char recData[255];
	int ret = recv(sclient, recData, 255, 0);
	if (ret > 0) {
		recData[ret] = 0x00;
		printf(recData);
	}
	Recv(L, recData, strlen(recData));
}

void SimpleServer::Recv(lua_State* L, const char* data, size_t data_len)
{
	if (!lua_get_object_function(L, this, "on_call_with_handle")) {
		printf("SimpleServer::Recv on_call_with_hanldle failed");
		return;
	}
	lua_pushstring(L, data);
	lua_call_function(L, nullptr, 1, 0);
}
