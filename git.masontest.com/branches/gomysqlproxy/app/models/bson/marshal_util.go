// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author Tony Xu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

// Utility functions for custom encoders

package bson

import (
	"time"

    "git.masontest.com/branches/gomysqlproxy/app/models/bytes2"
)

// EncodeInterface bson encodes an interface{}. Elements
// can be basic bson encodable types, or []interface{},
// or map[string]interface{}, whose elements have to in
// turn be bson encodable.
func EncodeInterface(buf *bytes2.ChunkedWriter, key string, val interface{}) {
	if val == nil {
		EncodePrefix(buf, Null, key)
		return
	}
	switch val := val.(type) {
	case string:
		EncodeString(buf, key, val)
	case []byte:
		EncodeBinary(buf, key, val)
	case int64:
		EncodeInt64(buf, key, val)
	case int32:
		EncodeInt32(buf, key, val)
	case int:
		EncodeInt(buf, key, val)
	case uint64:
		EncodeUint64(buf, key, val)
	case uint32:
		EncodeUint32(buf, key, val)
	case uint:
		EncodeUint(buf, key, val)
	case float64:
		EncodeFloat64(buf, key, val)
	case bool:
		EncodeBool(buf, key, val)
	case map[string]interface{}:
		if val == nil {
			EncodePrefix(buf, Null, key)
			return
		}
		EncodePrefix(buf, Object, key)
		lenWriter := NewLenWriter(buf)
		for k, v := range val {
			EncodeInterface(buf, k, v)
		}
		lenWriter.Close()
	case []interface{}:
		if val == nil {
			EncodePrefix(buf, Null, key)
			return
		}
		EncodePrefix(buf, Array, key)
		lenWriter := NewLenWriter(buf)
		for i, v := range val {
			EncodeInterface(buf, Itoa(i), v)
		}
		lenWriter.Close()
	case time.Time:
		EncodeTime(buf, key, val)
	default:
		panic(NewBsonError("don't know how to marshal %T", val))
	}
}

// EncodeStringArray bson encodes a []string
func EncodeStringArray(buf *bytes2.ChunkedWriter, name string, values []string) {
	if values == nil {
		EncodePrefix(buf, Null, name)
		return
	}
	EncodePrefix(buf, Array, name)
	lenWriter := NewLenWriter(buf)
	for i, val := range values {
		EncodeString(buf, Itoa(i), val)
	}
	lenWriter.Close()
}
