package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/AGou-ops/zinx/utils"
	"github.com/AGou-ops/zinx/ziface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadlen() uint32 {
	// DataLen uint32 4bytes + DataId uint32 4 bytes
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuf := bytes.NewBuffer([]byte{})

	// 写datalen
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	// 写dataId
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 写data
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	dataBuf := bytes.NewReader(binaryData)

	msg := &Message{}
	// 读dataLen
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读dataID
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	if utils.GlobalObject.MaxPacketSize > 0 &&
		msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large data packet")
	}

	// 读data
	return msg, nil
}
