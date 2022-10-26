package tinyserializer

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

// Serializer is a struct that can serialize and deserialize data
type Serializer struct {
	// The buffer to serialize to or deserialize from
	buffer   *bytes.Buffer
	compress bool
}

// Now, all fields will be stored along
// with their size in the following format:
// Cannot serialize maps yet
// [field size][field data][field size][field data]

// NewSerializer creates a new serializer
func NewSerializer() *Serializer {
	return &Serializer{
		buffer:   new(bytes.Buffer),
		compress: false,
	}
}

func (s *Serializer) SetCompress(compress bool) *Serializer {
	s.compress = compress
	return s
}

func (s *Serializer) SetData(data []byte) *Serializer {
	s.buffer = bytes.NewBuffer(data)
	return s
}

// Serialize serializes the given data
func (s *Serializer) Serialize(data interface{}) ([]byte, error) {
	// Reset the buffer
	s.buffer.Reset()

	// Serialize the data
	err := s.serialize(data)
	if err != nil {
		return nil, err
	}

	if s.compress {
		return Compress(s.buffer.Bytes())
	}
	// Return the serialized data
	return s.buffer.Bytes(), nil
}

// Deserialize deserializes the given data
func (s *Serializer) Deserialize(data []byte, out interface{}) error {
	// Create a new serializer
	if s.compress {
		var err error
		data, err = Decompress(data)
		if err != nil {
			return err
		}
	}

	// Set the buffer to the given data
	s.buffer = bytes.NewBuffer(data)

	// Deserialize the data
	return s.deserialize(out)
}

// serialize serializes the given data
func (s *Serializer) serialize(data interface{}) error {
	// Get the value of the data
	value := reflect.ValueOf(data)

	// Check if the data is a pointer or struct
	if value.Kind() != reflect.Ptr && value.Kind() != reflect.Struct && value.Kind() != reflect.Slice && value.Kind() != reflect.Map {
		return fmt.Errorf("data must be a pointer, struct, map or slice")
	}
	// Get the value of the data
	value = GetValue(value)

	dataType := value.Type()
	// Check if the data is a struct
	if value.Kind() == reflect.Struct {
		return s.WriteStruct(value, dataType)
	} else if value.Kind() == reflect.Slice {
		value = value.Index(0)
		for i := 0; i < value.Len(); i++ {
			err := s.WriteField(value.Index(i), value.Index(i).Kind())
			if err != nil {
				return err
			}
		}
	} else if value.Kind() == reflect.Map {
		return s.serializeMapFields(value)
	}
	return nil
}

func (s *Serializer) WriteStruct(value reflect.Value, dataType reflect.Type) error {
	numFields := value.NumField()
	// Loop through all fields
	for i := 0; i < numFields; i++ {
		// Get the field
		field := value.Field(i)
		kind := field.Kind()
		// Check field tags
		if dataType.Field(i).Tag.Get("tiny") == "" {
			continue
		} else if dataType.Field(i).Tag.Get("tiny") == "-" {
			continue
		} else if dataType.Field(i).Tag.Get("tiny") == "omitempty" && field.IsZero() {
			continue
		}
		// Get the field type
		fieldType := dataType.Field(i)
		// Check if the field is exported
		if fieldType.PkgPath != "" {
			continue
		}
		// Check if the field is a slice or a map
		s.WriteField(field, kind)
	}
	return nil
}

func (s *Serializer) WriteField(field reflect.Value, kind reflect.Kind) error {
	var err error
	if kind == reflect.Slice {
		// Serialize the field
		length := field.Len()

		// Write the length of the slice
		if err = binary.Write(s.buffer, binary.LittleEndian, uint32(length)); err != nil {
			return fmt.Errorf("failed to write slice length: " + err.Error())
		}

		// Loop through all elements in the slice
		for i := 0; i < length; i++ {
			// Get the element
			element := field.Index(i)
			// Serialize the element
			if err = s.serializeField(element); err != nil {
				return err
			}
		}
	} else if kind == reflect.Map {
		// Serialize the field
		length := field.Len()

		// Write the length of the map
		err = binary.Write(s.buffer, binary.LittleEndian, uint32(length))
		if err != nil {
			return fmt.Errorf("failed to write map length: " + err.Error())
		}

		// Loop through all elements in the map
		return s.serializeMapFields(field)
	} else {
		// Serialize the field
		return s.serializeField(field)
	}
	return nil
}

func (s *Serializer) serializeMapFields(value reflect.Value) error {
	// Get the number of fields
	numFields := value.MapRange()

	// Loop through all fields
	for numFields.Next() {
		// Get the field
		// Serialize the field
		if err := s.serializeField(numFields.Key()); err != nil {
			return fmt.Errorf("failed to serialize field: " + err.Error())
		}
		// Serialize the field
		if err := s.serializeField(numFields.Value()); err != nil {
			return fmt.Errorf("failed to serialize field: " + err.Error())
		}
	}
	return nil
}

func GetValue(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Ptr {
		return value.Elem()
	}
	return value
}

// serializeField serializes the given field
func (s *Serializer) serializeField(field reflect.Value) error {
	// Get the field kind
	kind := field.Kind()

	// Check if the field is a struct
	if kind == reflect.Struct {
		// Serialize the struct
		return s.serialize(field.Interface())
	}

	// Get the field size
	usize := field.Type().Size()

	// Write the field data
	// If the field is a string, we need to convert it to a byte slice
	var ndata []byte
	switch kind {
	case reflect.String:
		ndata = []byte(field.String())
	case reflect.Bool:
		ndata = []byte{0}
		if field.Bool() {
			ndata[0] = 1
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ndata = make([]byte, usize)
		binary.LittleEndian.PutUint64(ndata, uint64(field.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ndata = make([]byte, usize)
		binary.LittleEndian.PutUint64(ndata, field.Uint())
	case reflect.Float32, reflect.Float64:
		ndata = make([]byte, usize)
		binary.LittleEndian.PutUint64(ndata, math.Float64bits(field.Float()))
	case reflect.Complex64, reflect.Complex128:
		ndata = make([]byte, usize)
		binary.LittleEndian.PutUint64(ndata, math.Float64bits(real(field.Complex())))
		binary.LittleEndian.PutUint64(ndata[usize/2:], math.Float64bits(imag(field.Complex())))
	case reflect.Slice:
		return s.WriteField(field, kind)
	case reflect.Map:
		return s.WriteField(field, kind)
	case reflect.Ptr:
		return s.serializeField(field.Elem())
	default:
		ndata = field.Bytes()
	}
	// Get data
	data, err := GetBytes(ndata)
	if err != nil {
		return err
	}
	// Write the field data
	_, err = s.buffer.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write field data: " + err.Error())
	}

	return nil
}

func GetBytes(ndata []byte) ([]byte, error) {
	// Convert uinptr to uint64
	size := uint16(len(ndata))

	var b []byte = make([]byte, 2+size)
	binary.LittleEndian.PutUint16(b, size)
	copy(b[2:], ndata)

	return b, nil
}

func (s *Serializer) CheckTag(dataType reflect.Type, field reflect.Value, i int) bool {
	dt_field := dataType.Field(i)
	fieldtag := dt_field.Tag.Get("tiny")
	if fieldtag == "-" {
		return false
	} else if fieldtag == "omitempty" && field.IsZero() {
		return false
	} else if fieldtag != "" {
		// Perform noop
		return true
	} else {
		return false
	}
}

// deserialize deserializes the given data
func (s *Serializer) deserialize(data interface{}) error {
	// Get the value of the data
	value := reflect.ValueOf(data)

	// Check if the data is a pointer
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("data is not a pointer %s", value.Kind())
	}

	// Get the value of the data
	value = value.Elem()

	// Check if the data is a struct
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("deserialization error; data is not a struct [%s]", value.Kind())
	}

	// Get the type of the data
	dataType := value.Type()

	// Get the number of fields
	numFields := value.NumField()

	// Loop through all fields
	for i := 0; i < numFields; i++ {
		// Get the field
		field := value.Field(i)

		// Check field tags
		if !s.CheckTag(dataType, field, i) {
			continue
		}
		// Get the field type
		fieldType := dataType.Field(i)

		// Check if the field is exported
		if fieldType.PkgPath != "" {
			continue
		}

		// Check if the field is a slice or a map
		if field.Kind() == reflect.Slice {
			s.deserializeSlice(field)
		} else if field.Kind() == reflect.Map {
			s.deserializeMap(field)
		} else {
			// Deserialize the field
			err := s.deserializeField(field)
			if err != nil {
				return fmt.Errorf("failed to deserialize field: " + err.Error())
			}
		}
	}

	return nil
}

func (s *Serializer) deserializeSlice(field reflect.Value) error {
	// Deserialize the field
	// Get the length of the slice
	var length uint32
	err := binary.Read(s.buffer, binary.LittleEndian, &length)
	if err != nil {
		return fmt.Errorf("failed to read slice length: " + err.Error())
	}

	// Create a new slice
	field.Set(reflect.MakeSlice(field.Type(), int(length), int(length)))
	// Loop through all elements
	for j := 0; j < int(length); j++ {
		// Deserialize the element
		err := s.deserializeField(field.Index(j))
		if err != nil {
			return fmt.Errorf("failed to deserialize slice element: " + err.Error())
		}
	}
	return nil
}

func (s *Serializer) deserializeMap(field reflect.Value) error {
	// Deserialize the field
	// Get the length of the map
	var length uint32
	err := binary.Read(s.buffer, binary.LittleEndian, &length)
	if err != nil {
		return fmt.Errorf("failed to read map length: " + err.Error())
	}
	// Create a new map
	field.Set(reflect.MakeMap(field.Type()))
	// Loop through all elements
	for j := 0; j < int(length); j++ {
		// Deserialize the key
		key := reflect.New(field.Type().Key()).Elem()
		err := s.deserializeField(key)
		if err != nil {
			return fmt.Errorf("failed to deserialize map key: " + err.Error())
		}
		// Deserialize the value
		value := reflect.New(field.Type().Elem()).Elem()
		err = s.deserializeField(value)
		if err != nil {
			return fmt.Errorf("failed to deserialize map value: " + err.Error())
		}
		// Set the map value
		field.SetMapIndex(key, value)
	}
	return nil
}

// deserializeField deserializes the given field
func (s *Serializer) deserializeField(field reflect.Value) error {
	// Get the field kind
	kind := field.Kind()

	// Check if the field is a struct
	switch kind {
	case reflect.Struct:
		return s.deserialize(field.Addr().Interface())
	case reflect.Ptr:
		// Get the pointer type
		// Check if the pointer is nil
		if ptrType := field.Type(); field.IsNil() {
			// Create a new pointer
			field.Set(reflect.New(ptrType.Elem()))
		}
		// Deserialize the pointer
		return s.deserialize(field.Interface())
	case reflect.Slice:
		return s.deserializeSlice(field)
	case reflect.Map:
		return s.deserializeMap(field)
	}

	// Get the field size
	var size uint16
	err := binary.Read(s.buffer, binary.LittleEndian, &size)
	if err != nil {
		return fmt.Errorf("failed to read field size: %v", err)
	}

	// Read field data for the given size
	data := make([]byte, size)
	err = binary.Read(s.buffer, binary.LittleEndian, data)
	if err != nil {
		return fmt.Errorf("failed to read field data: %v", err)
	}

	// Set the field data
	// If the field is a string, we need to convert it to a string
	switch kind {
	case reflect.String:
		field.SetString(string(data))
	case reflect.Bool:
		field.SetBool(data[0] == 1)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(int64(binary.LittleEndian.Uint64(data)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.SetUint(binary.LittleEndian.Uint64(data))
	case reflect.Float32, reflect.Float64:
		field.SetFloat(math.Float64frombits(binary.LittleEndian.Uint64(data)))
	case reflect.Complex64, reflect.Complex128:
		field.SetComplex(complex(math.Float64frombits(binary.LittleEndian.Uint64(data)), math.Float64frombits(binary.LittleEndian.Uint64(data[8:]))))
	default:
		field.SetBytes(data)
	}

	return nil
}

func Compress(data []byte) ([]byte, error) {
	// Create a new buffer
	buffer := new(bytes.Buffer)

	// Create a new gzip writer
	writer := gzip.NewWriter(buffer)

	// Write the data to the writer
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	// Create a new buffer
	buffer := bytes.NewBuffer(data)

	// Create a new gzip reader
	reader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}

	// Read the data from the reader
	result, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Close the reader
	err = reader.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}
