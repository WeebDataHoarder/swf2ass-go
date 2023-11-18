package types

import (
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/japanese"
	"io"
	"math"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"unsafe"
)

type DataReader interface {
	io.ByteReader
	io.Reader
	ReadBits(n uint8) (u uint64, err error)
	Align() (skipped uint8)
}

type TypeReader interface {
	SWFRead(reader DataReader, swfVersion uint8) error
}

type TypeDefault interface {
	SWFDefault(swfVersion uint8)
}

type TypeFuncConditional func(swfVersion uint8) bool

type TypeFuncNumber func(swfVersion uint8) uint64

func ReadBool(r DataReader) (d bool, err error) {
	v, err := r.ReadBits(1)
	if err != nil {
		return false, err
	}
	return v == 1, nil
}

func ReadUB[T ~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8](r DataReader, n uint64) (d T, err error) {
	v, err := r.ReadBits(uint8(n))
	if err != nil {
		return 0, err
	}
	return T(v), nil
}

func ReadSB[T ~int | ~int64 | ~int32 | ~int16 | ~int8](r DataReader, n uint64) (d T, err error) {
	v, err := r.ReadBits(uint8(n))
	if err != nil {
		return 0, err
	}
	//TODO: check
	//sign bit is set
	if v&(1<<(n-1)) > 0 {
		v |= (math.MaxUint64 >> (64 - (n - 1))) << (64 - (n - 1))
	}
	return T(v), nil
}

func ReadFB(r DataReader, n uint64) (d Fixed, err error) {
	v, err := r.ReadBits(uint8(n))
	if err != nil {
		return 0, err
	}
	return Fixed(v), nil
}

func ReadU64[T ~uint64](r DataReader, d *T) (err error) {
	r.Align()
	var buf [8]byte
	_, err = io.ReadFull(r, buf[:])
	if err != nil {
		return err
	}

	*d = T(binary.LittleEndian.Uint16(buf[:]))

	return nil
}

func ReadU32[T ~uint32](r DataReader, d *T) (err error) {
	r.Align()
	var buf [4]byte
	_, err = io.ReadFull(r, buf[:])
	if err != nil {
		return err
	}

	*d = T(binary.LittleEndian.Uint32(buf[:]))

	return nil
}

func ReadU24[T ~uint32](r DataReader, d *T) (err error) {
	r.Align()
	var buf [4]byte
	_, err = io.ReadFull(r, buf[:3])
	if err != nil {
		return err
	}

	*d = T(binary.LittleEndian.Uint32(buf[:]))

	return nil
}

func ReadU16[T ~uint16](r DataReader, d *T) (err error) {
	r.Align()
	var buf [2]byte
	_, err = io.ReadFull(r, buf[:])
	if err != nil {
		return err
	}

	*d = T(binary.LittleEndian.Uint16(buf[:]))

	return nil
}

func ReadU8[T ~uint8](r DataReader, d *T) (err error) {
	r.Align()
	var buf [1]byte
	_, err = io.ReadFull(r, buf[:])
	if err != nil {
		return err
	}

	*d = T(buf[0])

	return nil
}

func ReadSI64[T ~int64](r DataReader, d *T) (err error) {
	var v uint64
	err = ReadU64(r, &v)
	*d = T(v)
	return err
}

func ReadSI32[T ~int32](r DataReader, d *T) (err error) {
	var v uint32
	err = ReadU32(r, &v)
	*d = T(v)
	return err
}

func ReadSI16[T ~int16](r DataReader, d *T) (err error) {
	var v uint16
	err = ReadU16(r, &v)
	*d = T(v)
	return err
}

func ReadSI8[T ~int8](r DataReader, d *T) (err error) {
	var v uint8
	err = ReadU8(r, &v)
	*d = T(v)
	return err
}

func ReadArraySI8[T ~int8](r DataReader, n int) (d []T, err error) {
	d = make([]T, n)
	for i := range d {
		err = ReadSI8(r, &d[i])
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func ReadArraySI16[T ~int16](r DataReader, n int) (d []T, err error) {
	d = make([]T, n)
	for i := range d {
		err = ReadSI16(r, &d[i])
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func ReadArrayU8[T ~uint8](r DataReader, n int) (d []T, err error) {
	d = make([]T, n)
	for i := range d {
		err = ReadU8(r, &d[i])
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func ReadArrayU16[T ~uint16](r DataReader, n int) (d []T, err error) {
	d = make([]T, n)
	for i := range d {
		err = ReadU16(r, &d[i])
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func ReadArrayU24[T ~uint32](r DataReader, n int) (d []T, err error) {
	d = make([]T, n)
	for i := range d {
		err = ReadU24(r, &d[i])
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func ReadArrayU32[T ~uint32](r DataReader, n int) (d []T, err error) {
	d = make([]T, n)
	for i := range d {
		err = ReadU32(r, &d[i])
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func ReadArrayU64[T ~uint64](r DataReader, n int) (d []T, err error) {
	d = make([]T, n)
	for i := range d {
		err = ReadU64(r, &d[i])
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func ReadEncodedU32[T ~uint32](r DataReader) (d T, err error) {
	//TODO: verify
	r.Align()
	v, err := binary.ReadUvarint(r)
	if err != nil {
		return 0, err
	}

	return T(v), nil
}

func ReadString(r DataReader, swfVersion uint8) (d string, err error) {
	var v uint8
	for {
		err = ReadU8[uint8](r, &v)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return "", io.ErrUnexpectedEOF
			}
			return "", err
		}
		if v == 0 {
			break
		}
		d += string(v)
	}

	if swfVersion >= 6 {
		//always utf-8
		return d, nil
	}
	//TODO: detect
	decoder := japanese.ShiftJIS.NewDecoder()
	newD, err := decoder.String(d)
	if err != nil {
		//probably ascii?
		return d, nil
	}
	return newD, nil
}

var typeReaderReflect = reflect.TypeOf((*TypeReader)(nil)).Elem()

var typeDefaultReflect = reflect.TypeOf((*TypeDefault)(nil)).Elem()

var typeFuncConditionalReflect = reflect.TypeOf((*TypeFuncConditional)(nil)).Elem()

var typeFuncNumberReflect = reflect.TypeOf((*TypeFuncNumber)(nil)).Elem()

func getNestedType(el reflect.Value, fields ...string) reflect.Value {
	for len(fields) > 0 && fields[0] != "" {
		if strings.HasSuffix(fields[0], "()") {
			n := strings.TrimSuffix(fields[0], "()")
			m := el.MethodByName(n)
			if !m.IsValid() {
				m = el.Addr().MethodByName(n)
			}
			el = m
		} else {
			el = el.FieldByName(fields[0])
		}
		fields = fields[1:]
	}
	return el
}

func ReadType(r DataReader, swfVersion uint8, data any) (err error) {
	if tr, ok := data.(TypeDefault); ok {
		tr.SWFDefault(swfVersion)
	}
	if tr, ok := data.(TypeReader); ok {
		return tr.SWFRead(r, swfVersion)
	}

	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Pointer {
		return errors.New("not a pointer")
	}

	return ReadTypeInner(r, swfVersion, dataValue.Elem(), "", data)
}

func ReadTypeInner(r DataReader, swfVersion uint8, root reflect.Value, tags reflect.StructTag, data any) (err error) {
	if tr, ok := data.(TypeDefault); ok {
		tr.SWFDefault(swfVersion)
	}
	if tr, ok := data.(TypeReader); ok {
		return tr.SWFRead(r, swfVersion)
	}

	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()

	if dataType.Kind() != reflect.Pointer {
		return errors.New("not a pointer")
	}

	dataElement := dataValue.Elem()
	dataElementType := dataElement.Type()

	switch dataElementType.Kind() {
	case reflect.Struct:
		//get struct flags
		var flags []string
		flagsField, ok := dataElementType.FieldByName("_")
		if ok {
			flags = strings.Split(flagsField.Tag.Get("swfFlags"), ",")
		}

		structRoot := root

		if slices.Contains(flags, "align") {
			r.Align()
		}

		if slices.Contains(flags, "root") {
			structRoot = dataElement
		}

		n := dataElementType.NumField()
		for i := 0; i < n; i++ {
			fieldValue := dataElement.Field(i)
			fieldType := dataElementType.Field(i)

			if _, ok := fieldType.Tag.Lookup("swfIgnore"); ok || fieldType.Name == "_" {
				continue
			}

			//Check if we should read this entry or not
			if swfTag := fieldType.Tag.Get("swfCondition"); swfTag != "" {
				splits := strings.Split(swfTag, ".")
				el := getNestedType(structRoot, splits...)

				switch el.Kind() {
				case reflect.Bool:
					if !el.Bool() {
						continue
					}
				case reflect.Func:
					if el.Type().AssignableTo(typeFuncConditionalReflect) {
						values := el.Call([]reflect.Value{reflect.ValueOf(swfVersion)})
						if !values[0].Bool() {
							continue
						}
					} else {
						return fmt.Errorf("invalid conditional method %s", swfTag)
					}
				default:
					return fmt.Errorf("invalid conditional type %s", swfTag)
				}
			}

			if swfTag := fieldType.Tag.Get("swfBits"); swfTag != "" {
				entries := strings.Split(swfTag, ",")
				splits := strings.Split(entries[0], ".")
				el := getNestedType(structRoot, splits...)

				bitFlags := entries[1:]

				var nbits uint64

				if len(splits) == 1 && len(splits[0]) == 0 && len(bitFlags) > 0 {
					//numerical fixed
					nbits, err = strconv.ParseUint(bitFlags[0], 10, 0)
					if err != nil {
						return err
					}
				} else {
					switch el.Kind() {
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						nbits = el.Uint()
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						nbits = uint64(el.Int())
					case reflect.Func:
						if el.Type().AssignableTo(typeFuncNumberReflect) {
							values := el.Call([]reflect.Value{reflect.ValueOf(swfVersion)})
							nbits = values[0].Uint()
						} else {
							return fmt.Errorf("invalid nbits method %s", swfTag)
						}
					default:
						return fmt.Errorf("invalid nbits type %s", swfTag)
					}
				}

				if slices.Contains(bitFlags, "signed") {
					value, err := ReadSB[int64](r, nbits)
					if err != nil {
						return err
					}
					fieldValue.SetInt(value)
				} else {
					value, err := ReadUB[uint64](r, nbits)
					if err != nil {
						return err
					}
					fieldValue.SetUint(value)
				}
				continue
			}

			err = ReadTypeInner(r, swfVersion, structRoot, fieldType.Tag, fieldValue.Addr().Interface())
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		var number uint64
		readMoreRecords := func() bool {
			more := number > 0
			number--
			return more
		}

		sliceType := dataElementType.Elem()

		if swfTag := tags.Get("swfCount"); swfTag != "" {
			splits := strings.Split(swfTag, ".")
			el := getNestedType(root, splits...)

			switch el.Kind() {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				number = el.Uint()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				number = uint64(el.Int())
			case reflect.Func:
				if el.Type().AssignableTo(typeFuncNumberReflect) {
					values := el.Call([]reflect.Value{reflect.ValueOf(swfVersion)})
					number = values[0].Uint()
				} else {
					return fmt.Errorf("invalid count method %s", swfTag)
				}
			default:
				return fmt.Errorf("invalid count type %s", swfTag)
			}
		} else if swfTag := tags.Get("swfMore"); swfTag != "" {
			splits := strings.Split(swfTag, ".")
			el := getNestedType(root, splits...)

			switch el.Kind() {
			case reflect.Func:
				if el.Type().AssignableTo(typeFuncConditionalReflect) {
					readMoreRecords = func() bool {
						values := el.Call([]reflect.Value{reflect.ValueOf(swfVersion)})
						return values[0].Bool()
					}
				} else {
					return fmt.Errorf("invalid more method %s", swfTag)
				}
			default:
				return fmt.Errorf("invalid more type %s", swfTag)
			}
		}

		if sliceType.Kind() == reflect.Pointer {
			return errors.New("unsupported pointer in slice")
		}

		newSlice := reflect.MakeSlice(dataElementType, 0, int(number))
		for readMoreRecords() {
			value := reflect.New(sliceType)
			err = ReadTypeInner(r, swfVersion, root, tags, value.Interface())
			if err != nil {
				return err
			}
			newSlice = reflect.Append(newSlice, value.Elem())
		}
		//TODO: check this
		dataElement.Set(newSlice)
	case reflect.Array:
		if dataElementType.Elem().Kind() == reflect.Pointer {
			return errors.New("unsupported pointer in slice")
		}

		for i := 0; i < dataElement.Len(); i++ {
			err = ReadTypeInner(r, swfVersion, root, tags, dataElement.Index(i).Addr().Interface())
			if err != nil {
				return err
			}
		}
	case reflect.Bool:
		value, err := ReadBool(r)
		if err != nil {
			return err
		}
		dataElement.SetBool(value)
	case reflect.Uint8:
		var value uint8
		err = ReadU8(r, &value)
		if err != nil {
			return err
		}
		dataElement.SetUint(uint64(value))
	case reflect.Uint16:
		var value uint16
		err = ReadU16(r, &value)
		if err != nil {
			return err
		}
		dataElement.SetUint(uint64(value))
	case reflect.Uint32:
		var value uint32
		err = ReadU32(r, &value)
		if err != nil {
			return err
		}
		dataElement.SetUint(uint64(value))
	case reflect.Uint64:
		var value uint64
		err = ReadU64(r, &value)
		if err != nil {
			return err
		}
		dataElement.SetUint(value)
	case reflect.Int8:
		var value int8
		err = ReadSI8(r, &value)
		if err != nil {
			return err
		}
		dataElement.SetInt(int64(value))
	case reflect.Int16:
		var value int16
		err = ReadSI16(r, &value)
		if err != nil {
			return err
		}
		dataElement.SetInt(int64(value))
	case reflect.Int32:
		var value int32
		err = ReadSI32(r, &value)
		if err != nil {
			return err
		}
		dataElement.SetInt(int64(value))
	case reflect.Int64:
		var value int64
		err = ReadSI64(r, &value)
		if err != nil {
			return err
		}
		dataElement.SetInt(value)
	case reflect.String:
		value, err := ReadString(r, swfVersion)
		if err != nil {
			return err
		}
		dataElement.SetString(value)
	case reflect.Float32:
		var value uint32
		err = ReadU32(r, &value)
		if err != nil {
			return err
		}
		dataElement.SetFloat(float64(math.Float32frombits(value)))
	case reflect.Float64:
		var value uint64
		err = ReadU64(r, &value)
		if err != nil {
			return err
		}
		dataElement.SetFloat(float64(math.Float64frombits(value)))
	default:
		return fmt.Errorf("unsupported type: %s", dataElementType.Name())
	}

	return nil
}

var basicTypes = []reflect.Type{
	reflect.Bool:          reflect.TypeOf(false),
	reflect.Int:           reflect.TypeOf(int(0)),
	reflect.Int8:          reflect.TypeOf(int8(0)),
	reflect.Int16:         reflect.TypeOf(int16(0)),
	reflect.Int32:         reflect.TypeOf(int32(0)),
	reflect.Int64:         reflect.TypeOf(int64(0)),
	reflect.Uint:          reflect.TypeOf(uint(0)),
	reflect.Uint8:         reflect.TypeOf(uint8(0)),
	reflect.Uint16:        reflect.TypeOf(uint16(0)),
	reflect.Uint32:        reflect.TypeOf(uint32(0)),
	reflect.Uint64:        reflect.TypeOf(uint64(0)),
	reflect.Uintptr:       reflect.TypeOf(uintptr(0)),
	reflect.Float32:       reflect.TypeOf(float32(0)),
	reflect.Float64:       reflect.TypeOf(float64(0)),
	reflect.Complex64:     reflect.TypeOf(complex64(0)),
	reflect.Complex128:    reflect.TypeOf(complex128(0)),
	reflect.String:        reflect.TypeOf(""),
	reflect.UnsafePointer: reflect.TypeOf(unsafe.Pointer(nil)),
}

func underlyingType(t reflect.Type) (ret reflect.Type) {
	if t.Name() == "" {
		// t is an unnamed type. the underlying type is t itself
		return t
	}
	kind := t.Kind()
	if ret = basicTypes[kind]; ret != nil {
		return ret
	}
	switch kind {
	case reflect.Array:
		ret = reflect.ArrayOf(t.Len(), t.Elem())
	case reflect.Chan:
		ret = reflect.ChanOf(t.ChanDir(), t.Elem())
	case reflect.Map:
		ret = reflect.MapOf(t.Key(), t.Elem())
	case reflect.Func:
		n_in := t.NumIn()
		n_out := t.NumOut()
		in := make([]reflect.Type, n_in)
		out := make([]reflect.Type, n_out)
		for i := 0; i < n_in; i++ {
			in[i] = t.In(i)
		}
		for i := 0; i < n_out; i++ {
			out[i] = t.Out(i)
		}
		ret = reflect.FuncOf(in, out, t.IsVariadic())
	case reflect.Interface:
		// not supported
	case reflect.Ptr:
		ret = reflect.PtrTo(t.Elem())
	case reflect.Slice:
		ret = reflect.SliceOf(t.Elem())
	case reflect.Struct:
		// only partially supported: embedded fields
		// and unexported fields may cause panic in reflect.StructOf()
		defer func() {
			// if a panic happens, return t unmodified
			if recover() != nil && ret == nil {
				ret = t
			}
		}()
		n := t.NumField()
		fields := make([]reflect.StructField, n)
		for i := 0; i < n; i++ {
			fields[i] = t.Field(i)
		}
		ret = reflect.StructOf(fields)
	}
	return ret
}
