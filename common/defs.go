package common

/*
CS: client -> server
SC: server-> client 
*/
var (
	CS_REGISTER byte = 1    //注册
	SC_SYNC_POS byte = 2	//同步位置信息
	SC_PUSH_TOKEN byte = 3   //得到token
)