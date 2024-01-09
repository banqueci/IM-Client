package protocol

const MsgHeartbeat = 1
const MsgAuthToken = 2

const MsgAuthStatus = 3

//persistent
const MsgIm = 4
const MsgAck = 5

//deprecated
const MSGRst = 6

//persistent
const MsgGroupNotification = 7
const MsgGroupIm = 8

//deprecated
const MsgPeerAck = 9
const MsgInPutting = 10

//消息回执
const MsgReceipt = 11

//消息撤回
const ImMsgRetract = 12
const GroupMsgRetract = 13

const MsgOrderReceipt = 14

const MsgEnterRoom = 18
const MsgLeaveRoom = 19
const MsgRoomIm = 20

//转发消息
const MsgForward = 21

//客户端->服务端
const MsgSync = 26 //同步消息
//服务端->客服端
const MsgSyncBegin = 27
const MsgSyncEnd = 28

//通知客户端有新消息
const MsgSyncNotify = 29

//客户端->服务端
const MsgSyncGroup = 30 //同步超级群消息
//服务端->客服端
const MsgSyncGroupBegin = 31
const MsgSyncGroupEnd = 32

//通知客户端有新消息
const MsgSyncGroupNotify = 33

//客服端->服务端,更新服务器的syncKey
const MsgSyncKey = 34
const MsgGroupSyncKey = 35

const MsgSystem = 36

const MsgTransmission = 37

const MsgResult = 38

const MsgRoomGetState = 39

const MsgRoomState = 40

const MsgRoomSetState = 41

const MsgRoomNotify = 42

// 另一客户端登录
const SysAnother = 51

// 有人进入房间
const SysRoomEnter = 52

// 有人退出房间
const SysRoomLeave = 53

// 上层自定义的系统消息
const SysCustom = 54

// 单聊消息全部已读
const SysSingleMsgRead = 55

const MsgSyncRoomNotify = 56

const PlatformMobile = 101
const PlatformPc = 301

const RpcPublish = 103
const RpcPublishGroup = 104
const RpcLogin = 105
const RpcLogout = 106
const RpcPublishRoom = 107
const RpcSubscribeRoom = 108
const RpcUnSubscribeRoom = 109
const RpcRegister = 110
const RpcUNREGISTER = 111
const RpcUpdateRoomState = 112

const RpcSubscribeRoomMember = 113
const RpcUnsubscribeRoomMember = 114

const ImrRpcRegister = 115
const ImrRpcRequest = 116
const ImrRpcLogin = 117
const ImrRpcPublish = 118
const ImrRpcPublishRoom = 119
const ImrRpcPublishGroup = 120

const RpcPublishNewRoom = 121
const ImrRpcPublishNewRoom = 122
const TcpPing = 123
const TcpPong = 124

//客户端打开插件
const OpenPlugin = 201

//客户端关闭插件
const ClosePlugin = 202

//客户端插件执行结果
const ResultPlugin = 203

//客户端指令消息
const MsgOrder = 204

//客户端指令消息确认
const MsgOrderConfirm = 205

//客户端自定义指令消息
const MsgOrderCustomize = 206

//客户端自定义指令消息确认
const MsgOrderConfirmCustomize = 207

const OtherInstanceUser = 300
