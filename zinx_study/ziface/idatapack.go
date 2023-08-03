package ziface

type IDataPackage interface {
	// 获取包头的方法
	GetHeadLen() uint32
	// 封包的方法
	Pack(msg IMessage) ([]byte, error)
	// 拆包的方法
	Unpack([]byte) (IMessage, error)
}
