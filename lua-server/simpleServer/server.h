#include<WINSOCK2.H>
#include<STDIO.H>
#include<iostream>
#include<cstring>
#include<string>
#include "luna.h"
using namespace std;
#pragma comment(lib, "ws2_32.lib")
class SimpleServer {
    string ip;
    int port;
    SOCKET sclient;
public:
    DECLARE_LUA_CLASS(SimpleServer);
    bool Init();
    void Uninit();
    bool Connect();
    int Send(const char* sendData);
    char* Recv();
};