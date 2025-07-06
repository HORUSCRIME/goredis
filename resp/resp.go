package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type ValueType int

const (
	SimpleStringType ValueType = iota
	ErrorType
	IntegerType
	BulkStringType
	ArrayType
	NullType
	BooleanType
	DoubleType
	BigNumberType
)

type Value struct {
	Type    ValueType
	Str     string
	Num     int64
	Bulk    []byte
	Array   []Value
	Null    bool
	Boolean bool
	Double  float64
}

func NewSimpleString(s string) Value {
	return Value{Type: SimpleStringType, Str: s}
}

func NewError(s string) Value {
	return Value{Type: ErrorType, Str: s}
}

func NewInteger(i int64) Value {
	return Value{Type: IntegerType, Num: i}
}

func NewBulkString(b []byte) Value {
	return Value{Type: BulkStringType, Bulk: b}
}

func NewNullBulkString() Value {
	return Value{Type: BulkStringType, Bulk: nil, Null: true}
}

func NewArray(arr []Value) Value {
	return Value{Type: ArrayType, Array: arr}
}

func NewNullArray() Value {
	return Value{Type: ArrayType, Array: nil, Null: true}
}

func NewNull() Value {
	return Value{Type: NullType, Null: true}
}

func NewBoolean(b bool) Value {
	return Value{Type: BooleanType, Boolean: b}
}

func NewDouble(d float64) Value {
	return Value{Type: DoubleType, Double: d}
}

func Decode(reader *bufio.Reader) (Value, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return Value{}, err
	}
	if len(line) < 2 || line[len(line)-2] != '\r' || line[len(line)-1] != '\n' {
		return Value{}, fmt.Errorf("malformed RESP line ending: %q", line)
	}
	line = line[:len(line)-2]

	switch line[0] {
	case '+': 
		return NewSimpleString(string(line[1:])), nil
	case '-': 
		return NewError(string(line[1:])), nil
	case ':': 
		num, err := strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return Value{}, fmt.Errorf("invalid integer: %w", err)
		}
		return NewInteger(num), nil
	case '$': 
		length, err := strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return Value{}, fmt.Errorf("invalid bulk string length: %w", err)
		}
		if length == -1 {
			return NewNullBulkString(), nil
		}

		bulk := make([]byte, length)
		_, err = io.ReadFull(reader, bulk)
		if err != nil {
			return Value{}, fmt.Errorf("failed to read bulk string: %w", err)
		}
		_, err = reader.ReadBytes('\n') 
		if err != nil {
			return Value{}, fmt.Errorf("failed to read bulk string trailing CRLF: %w", err)
		}
		return NewBulkString(bulk), nil
	case '*': 
		length, err := strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return Value{}, fmt.Errorf("invalid array length: %w", err)
		}
		if length == -1 {
			return NewNullArray(), nil
		}

		arr := make([]Value, length)
		for i := 0; i < int(length); i++ {
			val, err := Decode(reader)
			if err != nil {
				return Value{}, fmt.Errorf("failed to decode array element: %w", err)
			}
			arr[i] = val
		}
		return NewArray(arr), nil
	case '_':
		return NewNull(), nil
	case '#':
		if string(line[1:]) == "t" {
			return NewBoolean(true), nil
		} else if string(line[1:]) == "f" {
			return NewBoolean(false), nil
		}
		return Value{}, fmt.Errorf("invalid boolean value: %q", line)
	case ',': 
		d, err := strconv.ParseFloat(string(line[1:]), 64)
		if err != nil {
			return Value{}, fmt.Errorf("invalid double value: %w", err)
		}
		return NewDouble(d), nil
	default:
		return Value{}, fmt.Errorf("unknown RESP type: %q", line)
	}
}

func Encode(writer *bufio.Writer, val Value) error {
	var err error 

	switch val.Type {
	case SimpleStringType:
		_, err = writer.WriteString(fmt.Sprintf("+%s\r\n", val.Str))
		return err
	case ErrorType:
		_, err = writer.WriteString(fmt.Sprintf("-%s\r\n", val.Str))
		return err
	case IntegerType:
		_, err = writer.WriteString(fmt.Sprintf(":%d\r\n", val.Num))
		return err
	case BulkStringType:
		if val.Null {
			_, err = writer.WriteString("$-1\r\n")
			return err
		}
		_, err = writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(val.Bulk), val.Bulk))
		return err
	case ArrayType:
		if val.Null {
			_, err = writer.WriteString("*-1\r\n")
			return err
		}
		_, err = writer.WriteString(fmt.Sprintf("*%d\r\n", len(val.Array)))
		if err != nil {
			return err
		}
		for _, elem := range val.Array {
			if err := Encode(writer, elem); err != nil {
				return err
			}
		}
		return nil
	case NullType: 
		_, err = writer.WriteString("_\r\n")
		return err
	case BooleanType: 
		if val.Boolean {
			_, err = writer.WriteString("#t\r\n")
		} else {
			_, err = writer.WriteString("#f\r\n")
		}
		return err
	case DoubleType: 
		_, err = writer.WriteString(fmt.Sprintf(",%f\r\n", val.Double))
		return err
	default:
		return fmt.Errorf("unsupported RESP type for encoding: %v", val.Type)
	}
}