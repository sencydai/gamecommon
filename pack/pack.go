package pack

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

const (
	DEFAULT_TAG     = 0xccee
	DEFAULT_CRC_KEY = 0x765d

	HEAD_SIZE = 12
)

func NewWriter(datas ...interface{}) *bytes.Buffer {
	writer := bytes.NewBuffer([]byte{})
	Write(writer, datas...)
	return writer
}

func GetBytes(datas ...interface{}) []byte {
	writer := NewWriter(datas...)
	return writer.Bytes()
}

func Read(reader *bytes.Reader, datas ...interface{}) {
	for _, data := range datas {
		switch v := data.(type) {
		case *bool, *int8, *uint8, *int16, *uint16, *int32, *uint32, *int64, *uint64, *float32, *float64:
			err := binary.Read(reader, binary.LittleEndian, v)
			if err != nil {
				panic(err.Error())
			}
		case *int:
			var vv int32
			err := binary.Read(reader, binary.LittleEndian, &vv)
			if err != nil {
				panic(err)
			}
			*v = int(vv)
		case *string:
			var l uint16
			err := binary.Read(reader, binary.LittleEndian, &l)
			if err != nil {
				panic(err.Error())
			}
			s := make([]byte, l)
			for i := uint16(0); i < l; i++ {
				s[i], err = reader.ReadByte()
				if err != nil {
					panic(err.Error())
				}
			}
			*v = string(s)
			_, err = reader.ReadByte()
			if err != nil {
				panic(err.Error())
			}
		default:
			panic("pack.Read invalid type " + reflect.TypeOf(data).String())
		}
	}
}

func Write(writer *bytes.Buffer, datas ...interface{}) {
	for _, data := range datas {
		switch v := data.(type) {
		case bool, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
			binary.Write(writer, binary.LittleEndian, v)
		case int:
			binary.Write(writer, binary.LittleEndian, int32(v))
		case []byte:
			writer.Write(v)
		case string:
			binary.Write(writer, binary.LittleEndian, uint16(len(v)))
			writer.Write([]byte(v))
			binary.Write(writer, binary.LittleEndian, byte(0))
		default:
			panic("pack.Write invalid type " + reflect.TypeOf(data).String())
		}
	}
}

func AllocPack(sysId, cmdId byte, data ...interface{}) *bytes.Buffer {
	writer := NewWriter(DEFAULT_TAG, 0, int16(0), DEFAULT_CRC_KEY, sysId, cmdId)
	Write(writer, data...)
	return writer
}

func EncodeWriter(writer *bytes.Buffer) []byte {
	data := writer.Bytes()
	encode(data)
	return data
}

func EncodeData(sysId, cmdId byte, data ...interface{}) []byte {
	writer := AllocPack(sysId, cmdId, data...)
	return EncodeWriter(writer)
}

func encode(data []byte) {
	Len := GetBytes(len(data) - HEAD_SIZE)
	for i := 0; i < len(Len); i++ {
		data[i+4] = Len[i]
	}
}
