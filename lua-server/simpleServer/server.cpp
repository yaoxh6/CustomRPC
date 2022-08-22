#include "server.h"

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

int SimpleServer::Send(lua_State* L)
{
	const char* sendData = lua_tostring(L, 1);
	send(sclient, sendData, strlen(sendData), 0);
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
