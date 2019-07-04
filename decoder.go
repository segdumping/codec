package buffer

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
	"fmt"
)

type CustomDecoder interface {
	LoadDecoder(*Decoder) error
}

type Decoder struct {
	buf 	*bytes.Buffer
}

func NewDecoder(b *bytes.Buffer) *Decoder {
	return &Decoder{buf: b}
}

func (d *Decoder) Decode(v interface{}) error {
	var err error
	switch v := v.(type) {
	case *string:
		if v != nil {
			*v, err = d.decodeString()
			return err
		}
	case *int:
		if v != nil {
			*v, err = d.decodeInt()
			return err
		}
	case *int8:
		if v != nil {
			*v, err = d.decodeInt8()
			return err
		}
	case *int16:
		if v != nil {
			*v, err = d.decodeInt16()
			return err
		}
	case *int32:
		if v != nil {
			*v, err = d.decodeInt32()
			return err
		}
	case *int64:
		if v != nil {
			*v, err = d.decodeInt64()
			return err
		}
	case *uint:
		if v != nil {
			*v, err = d.decodeUint()
			return err
		}
	case *uint8:
		if v != nil {
			*v, err = d.decodeUint8()
			return err
		}
	case *uint16:
		if v != nil {
			*v, err = d.decodeUint16()
			return err
		}
	case *uint32:
		if v != nil {
			*v, err = d.decodeUint32()
			return err
		}
	case *uint64:
		if v != nil {
			*v, err = d.decodeUint64()
			return err
		}
	case *bool:
		if v != nil {
			*v, err = d.decodeBool()
			return err
		}
	case *float32:
		if v != nil {
			*v, err = d.decodeFloat32()
			return err
		}
	case *float64:
		if v != nil {
			*v, err = d.decodeFloat64()
			return err
		}
	}

	return decodeValue(d, reflect.ValueOf(v).Elem())
}

func (d *Decoder) decodeLength() int {
	bs := d.buf.Next(2)
	len := binary.BigEndian.Uint16(bs)

	return int(len)
}

func (d *Decoder) decodeString() (string, error) {
	var l uint16
	if err := binary.Read(d.buf, binary.BigEndian, &l); err != nil {
		return "", err
	}

	return string(d.buf.Next(int(l))), nil
}

func (d *Decoder) decodeBool() (bool, error) {
	bs := d.buf.Next(1)

	return bs[0] != 0, nil
}

func (d *Decoder) decodeInt8() (int8, error) {
	bs := d.buf.Next(1)

	return int8(bs[0]), nil
}

func (d *Decoder) decodeUint8() (uint8, error) {
	bs := d.buf.Next(1)

	return bs[0], nil
}

func (d *Decoder) decodeInt16() (int16, error) {
	bs := d.buf.Next(2)

	return int16(binary.BigEndian.Uint16(bs)), nil
}

func (d *Decoder) decodeUint16() (uint16, error) {
	bs := d.buf.Next(2)

	return binary.BigEndian.Uint16(bs), nil
}

func (d *Decoder) decodeInt32() (int32, error) {
	bs := d.buf.Next(4)

	return int32(binary.BigEndian.Uint32(bs)), nil
}

func (d *Decoder) decodeUint32() (uint32, error) {
	bs := d.buf.Next(4)

	return binary.BigEndian.Uint32(bs), nil
}

func (d *Decoder) decodeInt64() (int64, error) {
	bs := d.buf.Next(8)

	return int64(binary.BigEndian.Uint64(bs)), nil
}

func (d *Decoder) decodeUint64() (uint64, error) {
	bs := d.buf.Next(8)
	v := binary.BigEndian.Uint64(bs)
	return v , nil
}

func (d *Decoder) decodeInt() (int, error) {
	bs := d.buf.Next(8)

	return int(binary.BigEndian.Uint64(bs)), nil
}

func (d *Decoder) decodeUint() (uint, error) {
	bs := d.buf.Next(8)

	return uint(binary.BigEndian.Uint64(bs)), nil
}

func (d *Decoder) decodeFloat32() (float32, error) {
	bs := d.buf.Next(4)
	v := binary.BigEndian.Uint32(bs)

	return math.Float32frombits(v), nil
}

func (d *Decoder) decodeFloat64() (float64, error) {
	bs := d.buf.Next(8)
	v := binary.BigEndian.Uint64(bs)

	return math.Float64frombits(v), nil
}

func decodeValue(d *Decoder, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Bool:
		return decodeBool(d, v)
	case reflect.String:
		return decodeString(d, v)
	case reflect.Int8:
		return decodeInt8(d, v)
	case reflect.Uint8:
		return decodeUint8(d, v)
	case reflect.Int16:
		return decodeInt16(d, v)
	case reflect.Uint16:
		return decodeUint16(d, v)
	case reflect.Int32:
		return decodeInt32(d, v)
	case reflect.Uint32:
		return decodeUint32(d, v)
	case reflect.Int64:
		return decodeInt64(d, v)
	case reflect.Uint64:
		return decodeUint64(d, v)
	case reflect.Int:
		return decodeInt(d, v)
	case reflect.Uint:
		return decodeUint(d, v)
	case reflect.Float32:
		return decodeFloat32(d, v)
	case reflect.Float64:
		return decodeFloat64(d, v)
	case reflect.Array:
		return decodeArray(d, v)
	case reflect.Slice:
		return decodeSlice(d, v)
	case reflect.Map:
		return decodeMap(d, v)
	case reflect.Struct:
		return decodeStruct(d, v)
	case reflect.Ptr:
		return decodePtr(d, v)
	}

	return nil
}

func decodeString(d *Decoder, v reflect.Value) (error) {
	vv, err := d.decodeString()
	if err != nil {
		return err
	}

	if err = canSet(v); err != nil {
		return err
	}

	v.SetString(vv)

	return nil
}

func decodeBool(d *Decoder, v reflect.Value) (error) {
	vv, err := d.decodeBool()
	if err != nil {
		return err
	}

	if err = canSet(v); err != nil {
		return err
	}

	v.SetBool(vv)

	return nil
}

func decodeInt8(d *Decoder, v reflect.Value) (error) {
	vv, err := d.decodeInt8()
	if err != nil {
		return err
	}

	if err = canSet(v); err != nil {
		return err
	}

	v.SetInt(int64(vv))

	return nil
}

func decodeUint8(d *Decoder, v reflect.Value) (error) {
	vv, err := d.decodeUint8()
	if err != nil {
		return err
	}

	if err = canSet(v); err != nil {
		return err
	}

	v.SetUint(uint64(vv))

	return nil
}

func decodeInt16(d *Decoder, v reflect.Value) (error) {
	vv, err := d.decodeInt16()
	if err != nil {
		return err
	}

	if err = canSet(v); err != nil {
		return err
	}

	v.SetInt(int64(vv))

	return nil
}

func decodeUint16(d *Decoder, v reflect.Value) (error) {
	vv, err := d.decodeUint16()
	if err != nil {
		return err
	}

	if err = canSet(v); err != nil {
		return err
	}

	v.SetUint(uint64(vv))

	return nil
}

func decodeInt32(d *Decoder, v reflect.Value) (error) {
	vv, err := d.decodeInt32()
	if err != nil {
		return err
	}

	if err = canSet(v); err != nil {
		return err
	}

	v.SetInt(int64(vv))

	return nil
}

func decodeUint32(d *Decoder, v reflect.Value) (error) {
	vv, err := d.decodeUint32()
	if err != nil {
		return err
	}

	if err = canSet(v); err != nil {
		return err
	}

	v.SetUint(uint64(vv))

	return nil
}

func decodeInt64(d *Decoder, v reflect.Value) (error) {
	vv, err := d.decodeInt64()
	if err != nil {
		return err
	}

	if err = canSet(v); err != nil {
		return err
	}

	v.SetInt(int64(vv))

	return nil
}

func decodeUint64(d *Decoder, v reflect.Value) (error) {
	vv, err := d.decodeUint64()
	if err != nil {
		return err
	}

	if err = canSet(v); err != nil {
		return err
	}

	v.SetUint(uint64(vv))

	return nil
}

func decodeInt(d *Decoder, v reflect.Value) (error) {
	vv, err := d.decodeInt64()
	if err != nil {
		return err
	}

	if err = canSet(v); err != nil {
		return err
	}

	v.SetInt(int64(vv))

	return nil
}

func decodeUint(d *Decoder, v reflect.Value) (error) {
	vv, err := d.decodeUint64()
	if err != nil {
		return err
	}

	if err = canSet(v); err != nil {
		return err
	}

	v.SetUint(uint64(vv))

	return nil
}

func decodeFloat32(d *Decoder, v reflect.Value) (error) {
	vv, err := d.decodeFloat32()
	if err != nil {
		return err
	}

	if err = canSet(v); err != nil {
		return err
	}

	v.SetFloat(float64(vv))

	return nil
}

func decodeFloat64(d *Decoder, v reflect.Value) (error) {
	vv, err := d.decodeFloat64()
	if err != nil {
		return err
	}

	if err = canSet(v); err != nil {
		return err
	}

	v.SetFloat(float64(vv))

	return nil
}

func decodeSlice(d *Decoder, v reflect.Value) error {
	vv, err := d.decodeUint8()
	if err != nil {
		return err
	}

	c := code(vv)
	if c == Nil {
		return nil
	}

	l, err := d.decodeUint16()
	if err != nil {
		return err
	}

	if l == 0 {
		return nil
	}

	if v.Cap() < int(l) {
		diff := int(l) - v.Cap()
		v = reflect.AppendSlice(v, reflect.MakeSlice(v.Type(), diff, diff))
	}

	for i := 0; i < int(l); i++ {
		sv := v.Index(i)
		if sv.IsNil() && sv.Kind() == reflect.Ptr{
			sv.Set(reflect.New(v.Type().Elem().Elem()))
		}

		if err := decodeValue(d, sv); err != nil {
			return err
		}
	}

	return nil
}

func decodeArray(d *Decoder, v reflect.Value) error {
	l, err := d.decodeUint16()
	if err != nil {
		return err
	}

	if int(l) > v.Len() {
		return fmt.Errorf("%s len is %d,but array has %d elements", v.Type(), v.Len(), l)
	}

	for i := 0; i < int(l); i++ {
		sv := v.Index(i)
		if err := decodeValue(d, sv); err != nil {
			return err
		}
	}

	return nil
}

func decodeMap(d *Decoder, v reflect.Value) error {
	vv, err := d.decodeUint8()
	if err != nil {
		return err
	}

	c := code(vv)
	if c == Nil {
		return nil
	}

	size, err := d.decodeUint16()
	if err != nil {
		return err
	}

	typ := v.Type()
	if v.IsNil() {
		v.Set(reflect.MakeMap(typ))
	}

	if size == 0 {
		return nil
	}

	return decodeMapValue(d, v, int(size))
}

func decodeMapValue(d *Decoder, v reflect.Value, size int) error {
	typ := v.Type()
	keyType := typ.Key()
	valueType := typ.Elem()

	for i := 0; i < size; i++ {
		mk := reflect.New(keyType).Elem()
		if err := decodeValue(d, mk); err != nil {
			return err
		}

		var mv reflect.Value
		if valueType.Kind() == reflect.Ptr {
			mv = reflect.New(valueType.Elem())
		} else {
			mv = reflect.New(valueType).Elem()
		}

		if err := decodeValue(d, mv); err != nil {
			return err
		}

		v.SetMapIndex(mk, mv)
	}

	return nil
}

func decodeStruct(d *Decoder, v reflect.Value) error {
	decoder := v.Interface().(CustomDecoder)

	return decoder.LoadDecoder(d)
}

func decodePtr(d *Decoder, v reflect.Value) error {
	vv, err := d.decodeUint8()
	if err != nil {
		return err
	}

	c := code(vv)
	if c == Nil {
		return nil
	}

	if v.IsNil() {
		return fmt.Errorf("buffer: Decode ptr is nil [%s]", v.Type())
	}

	if v.MethodByName("LoadDecoder").IsValid() {
		decoder := v.Interface().(CustomDecoder)
		return decoder.LoadDecoder(d)
	} else {
		return fmt.Errorf("buffer: Decode(unsupported %s)", v.Type())
	}
}

func canSet(v reflect.Value) error {
	if !v.CanSet() {
		return fmt.Errorf("decode(nonsettable %s)", v.Type())
	}
	return nil
}