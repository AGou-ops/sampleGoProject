package ziface

// Irequest接口
type IRequest interface {
	// 得到当前链接
	GetConnection() IConnection

	// 得到请求的消息数据
	GetData() []byte

	// 得到当前请求的消息ID
	GetMsgID() uint32
}
