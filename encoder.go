package buffer

import (
	"encoding/binary"
	"bytes"
	"math"
	"reflect"
	"fmt"
)

type CustomEncoder interface {
	SaveEncoder(*Encoder) error
}

type Encoder struct {
	buf 	*bytes.Buffer
}

func NewEncoder(b *bytes.Buffer) *Encoder {
	return &Encoder{buf: b}
}

func (e *Encoder) Bytes() []byte {
	return e.buf.Bytes()
}

func (e *Encoder) Encode(v interface{}) error {
	switch v := v.(type) {
	case bool:
		return e.encodeBool(v)
	case int8:
		return e.encodeUint8(uint8(v))
	case uint8:
		return e.encodeUint8(v)
	case int16:
		return e.encodeUint16(uint16(v))
	case uint16:
		return e.encodeUint16(v)
	case int32:
		return e.encodeUint32(uint32(v))
	case uint32:
		return e.encodeUint32(v)
	case int64:
		return e.encodeUint64(uint64(v))
	case uint64:
		return e.encodeUint64(v)
	case float32:
		return e.encodeFloat32(v)
	case float64:
		return e.encodeFloat64(v)
	case int:
		return e.encodeInt(v)
	case uint:
		return e.encodeUint(v)
	case string:
		return e.encodeString(v)
	}

	return encodeValue(e, reflect.ValueOf(v))
}

func (e *Encoder) encodeBool(v bool) error {
	return binary.Write(e.buf, binary.BigEndian, v)
}

func (e *Encoder) encodeUint8(v uint8) error {
	return binary.Write(e.buf, binary.BigEndian, v)
}

func (e *Encoder) encodeUint16(v uint16) error {
	return binary.Write(e.buf, binary.BigEndian, v)
}

func (e *Encoder) encodeUint32(v uint32) error {
	return binary.Write(e.buf, binary.BigEndian, v)
}

func (e *Encoder) encodeUint64(v uint64) error {
	return binary.Write(e.buf, binary.BigEndian, v)
}

func (e *Encoder) encodeInt(v int) error {
	return binary.Write(e.buf, binary.BigEndian, int64(v))
}

func (e *Encoder) encodeUint(v uint) error {
	return binary.Write(e.buf, binary.BigEndian, uint64(v))
}

func (e *Encoder) encodeFloat32(v float32) error {
	return binary.Write(e.buf, binary.BigEndian, math.Float32bits(v))
}

func (e *Encoder) encodeFloat64(v float64) error {
	return binary.Write(e.buf, binary.BigEndian, math.Float64bits(v))
}

func (e *Encoder) encodeNil() error {
	return e.encodeUint8(uint8(Nil))
}

func (e *Encoder) encodeValid() error {
	return e.encodeUint8(uint8(Valid))
}

func (e *Encoder) encodeString(v string) error {
	err := binary.Write(e.buf, binary.BigEndian, uint16(len(v)))
	if err != nil {
		return err
	}

	e.buf.Write([]byte(v))

	return nil
}

func encodeBool(e *Encoder, v reflect.Value) error {
	return e.encodeBool(v.Bool())
}

func encodeInt8(e *Encoder, v reflect.Value) error {
	return e.encodeUint8(uint8(v.Int()))
}

func encodeUint8(e *Encoder, v reflect.Value) error {
	return e.encodeUint8(uint8(v.Uint()))
}

func encodeInt16(e *Encoder, v reflect.Value) error {
	return e.encodeUint16(uint16(v.Int()))
}

func encodeUint16(e *Encoder, v reflect.Value) error {
	return e.encodeUint16(uint16(v.Uint()))
}

func encodeInt32(e *Encoder, v reflect.Value) error {
	return e.encodeUint32(uint32(v.Int()))
}

func encodeUint32(e *Encoder, v reflect.Value) error {
	return e.encodeUint32(uint32(v.Uint()))
}

func encodeInt64(e *Encoder, v reflect.Value) error {
	return e.encodeUint64(uint64(v.Int()))
}

func encodeUint64(e *Encoder, v reflect.Value) error {
	return e.encodeUint64(uint64(v.Uint()))
}

func encodeFloat32(e *Encoder, v reflect.Value) error {
	return e.encodeFloat32(float32(v.Float()))
}

func encodeFloat64(e *Encoder, v reflect.Value) error {
	return e.encodeFloat64(float64(v.Float()))
}

func encodeInt(e *Encoder, v reflect.Value) error {
	return e.encodeUint64(uint64(v.Int()))
}

func encodeUint(e *Encoder, v reflect.Value) error {
	return e.encodeUint64(v.Uint())
}

func encodeString(e *Encoder, v reflect.Value) error {
	return e.encodeString(v.String())
}

func encodeArray(e *Encoder, v reflect.Value) error {
	l := v.Len()
	if err := e.encodeUint16(uint16(l)); err != nil {
		return err
	}

	for i := 0; i < l; i++ {
		if err := encodeValue(e, v.Index(i)); err != nil {
			return err
		}
	}

	return nil
}

func encodeSlice(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		return e.encodeNil()
	} else {
		e.encodeValid()
	}

	return encodeArray(e, v)
}

func encodeMap(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		return e.encodeNil()
	} else {
		e.encodeValid()
	}

	if err := e.encodeUint16(uint16(v.Len())); err != nil {
		return err
	}

	for _, key := range v.MapKeys() {
		if err := encodeValue(e, key); err != nil {
			return err
		}

		if err := encodeValue(e, v.MapIndex(key)); err != nil {
			return err
		}
	}

	return nil
}

func encodeStruct(e *Encoder, v reflect.Value) error {
	encoder := v.Interface().(CustomEncoder)

	return encoder.SaveEncoder(e)
}

func encodePtr(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		return e.encodeNil()
	} else {
		e.encodeValid()
	}

	if v.MethodByName("SaveEncoder").IsValid() {
		encoder := v.Interface().(CustomEncoder)
		return encoder.SaveEncoder(e)
	} else {
		return fmt.Errorf("buffer: Encode(unsupported %s)", v.Type())
	}
}

func encodeValue(e *Encoder, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Bool:
		return encodeBool(e, v)
	case reflect.String:
		return encodeString(e, v)
	case reflect.Int8:
		return encodeInt8(e, v)
	case reflect.Uint8:
		return encodeUint8(e, v)
	case reflect.Int16:
		return encodeInt16(e, v)
	case reflect.Uint16:
		return encodeUint16(e, v)
	case reflect.Int32:
		return encodeInt32(e, v)
	case reflect.Uint32:
		return encodeUint32(e, v)
	case reflect.Int64:
		return encodeInt64(e, v)
	case reflect.Uint64:
		return encodeUint64(e, v)
	case reflect.Float32:
		return encodeFloat32(e, v)
	case reflect.Float64:
		return encodeFloat64(e, v)
	case reflect.Int:
		return encodeInt(e, v)
	case reflect.Uint:
		return encodeUint(e, v)
	case reflect.Array:
		return encodeArray(e, v)
	case reflect.Slice:
		return encodeSlice(e, v)
	case reflect.Map:
		return encodeMap(e, v)
	case reflect.Struct:
		return encodeStruct(e, v)
	case reflect.Ptr:
		return encodePtr(e, v)
	}

	return nil
}