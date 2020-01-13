/*
 * Copyright 2019, Offchain Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package value

import (
	"bytes"
	"io"

	"github.com/offchainlabs/arbitrum/packages/arb-util/common"
)

const (
	TypeCodeInt       uint8 = 0
	TypeCodeCodePoint uint8 = 1
	TypeCodeHashOnly  uint8 = 2
	TypeCodeTuple     uint8 = 3
)

func TypeCodeName(code uint8) string {
	switch code {
	case TypeCodeInt:
		return "Int"
	case TypeCodeTuple:
		return "Tuple"
	case TypeCodeHashOnly:
		return "HashOnly"
	case TypeCodeCodePoint:
		return "CodePoint"
	default:
		return "Unknown"
	}
}

type Value interface {
	TypeCode() uint8
	InternalTypeCode() uint8
	Clone() Value
	CloneShallow() Value
	Equal(Value) bool
	Hash() common.Hash
	Size() int64
	Marshal(io.Writer) error
	MarshalForProof(io.Writer) error
}

func Eq(x, y Value) bool {
	return x.Equal(y)
}

type UnmarshalError struct {
	str string
}

func (e UnmarshalError) Error() string {
	return e.str
}

func MarshalValue(v Value, w io.Writer) error {
	_, err := w.Write([]byte{v.InternalTypeCode()})
	if err != nil {
		return err
	}
	return v.Marshal(w)
}

func MarshalValueToBytes(val Value) []byte {
	var buf bytes.Buffer
	_ = MarshalValue(val, &buf)
	return buf.Bytes()
}

func MarshalValueForProof(v Value, w io.Writer) error {
	_, err := w.Write([]byte{v.InternalTypeCode()})
	if err != nil {
		return err
	}
	return v.MarshalForProof(w)
}

func UnmarshalValueWithType(tipe byte, r io.Reader) (Value, error) {
	switch {
	case tipe == TypeCodeInt:
		return NewIntValueFromReader(r)
	case tipe == TypeCodeCodePoint:
		return NewCodePointValueFromReader(r)
	case tipe == TypeCodeHashOnly:
		return NewHashOnlyValueFromReader(r)
	case tipe <= TypeCodeTuple+MaxTupleSize:
		return NewSizedTupleFromReader(r, tipe-TypeCodeTuple)
	default:
		return NewEmptyTuple(), UnmarshalError{"Unmarshal: invalid value type"}
	}
}

func UnmarshalValue(r io.Reader) (Value, error) {
	tipe := make([]byte, 1)
	_, err := io.ReadFull(r, tipe)
	if err != nil {
		return NewEmptyTuple(), err
	}
	return UnmarshalValueWithType(tipe[0], r)
}

func UnmarshalValueFromBytes(val []byte) (Value, error) {
	buf := bytes.NewReader(val)
	return UnmarshalValue(buf)
}
