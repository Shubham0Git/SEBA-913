/*
 * Copyright (c) 2018 - present.  Boling Consulting Solutions (bcsw.net)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*
 * NOTE: This file was generated, manual edits will be overwritten!
 *
 * Generated by 'goCodeGenerator.py':
 *              https://github.com/cboling/OMCI-parser/README.md
 */
package generated

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/google/gopacket"
	"sort"
	"strings"
)

type AttributeDefinitionMap map[uint]*AttributeDefinition

// AttributeDefinition defines a single specific Managed Entity's attributes
type AttributeDefinition struct {
	Name         string
	Index        uint
	DefValue     interface{} // Note: Not supported yet
	Size         int
	Access       mapset.Set // AttributeAccess...
	Constraint   func(interface{}) *ParamError
	Avc          bool // If true, an AVC notification can occur for the attribute
	Tca          bool // If true, a threshold crossing alert alarm notification can occur for the attribute
	Counter      bool // If true, this attribute is a PM counter
	Optional     bool // If true, attribute is option, else mandatory
	TableSupport bool // If true, attribute is a table
	Deprecated   bool // If true, this attribute is deprecated and only 'read' operations (if-any) performed
}

func (attr *AttributeDefinition) String() string {
	return fmt.Sprintf("AttributeDefinition: %v (%v): Size: %v, Default: %v, Access: %v",
		attr.GetName(), attr.GetIndex(), attr.GetSize(), attr.GetDefault(), attr.GetAccess())
}
func (attr *AttributeDefinition) GetName() string         { return attr.Name }
func (attr *AttributeDefinition) GetIndex() uint          { return attr.Index }
func (attr *AttributeDefinition) GetDefault() interface{} { return attr.DefValue }
func (attr *AttributeDefinition) GetSize() int            { return attr.Size }
func (attr *AttributeDefinition) GetAccess() mapset.Set   { return attr.Access }
func (attr *AttributeDefinition) GetConstraints() func(interface{}) *ParamError {
	return attr.Constraint
}
func (attr *AttributeDefinition) IsTableAttribute() bool {
	return attr.TableSupport
}

func (attr *AttributeDefinition) Decode(data []byte, df gopacket.DecodeFeedback, msgType byte) (interface{}, error) {
	if attr.IsTableAttribute() {
		value, err := attr.tableAttributeDecode(data, df, msgType)
		if err != nil {
			return nil, err
		}
		if attr.GetConstraints() != nil {
			if omciErr := attr.GetConstraints()(value); omciErr != nil {
				return nil, omciErr.GetError()
			}
		}
		return value, nil
	}
	size := attr.GetSize()

	if len(data) < size {
		df.SetTruncated()
		return nil, NewMessageTruncatedError("packet too small for field")
	}
	switch attr.GetSize() {
	default:
		value := make([]byte, size)
		copy(value, data[:size])
		if attr.GetConstraints() != nil {
			if omciErr := attr.GetConstraints()(value); omciErr != nil {
				return nil, omciErr.GetError()
			}
		}
		return value, nil
	case 1:
		value := data[0]
		if attr.GetConstraints() != nil {
			if omciErr := attr.GetConstraints()(value); omciErr != nil {
				return nil, omciErr.GetError()
			}
		}
		return value, nil
	case 2:
		value := binary.BigEndian.Uint16(data[0:2])
		if attr.GetConstraints() != nil {
			if omciErr := attr.GetConstraints()(value); omciErr != nil {
				return nil, omciErr.GetError()
			}
		}
		return value, nil
	case 4:
		value := binary.BigEndian.Uint32(data[0:4])
		if attr.GetConstraints() != nil {
			if omciErr := attr.GetConstraints()(value); omciErr != nil {
				return nil, omciErr.GetError()
			}
		}
		return value, nil
	case 8:
		value := binary.BigEndian.Uint64(data[0:8])
		if attr.GetConstraints() != nil {
			omciErr := attr.GetConstraints()(value)
			if omciErr != nil {
				return nil, omciErr.GetError()
			}
		}
		return value, nil
	}
}

func (attr *AttributeDefinition) SerializeTo(value interface{}, b gopacket.SerializeBuffer,
	msgType byte, bytesAvailable int) (int, error) {
	if attr.IsTableAttribute() {
		return attr.tableAttributeSerializeTo(value, b, msgType, bytesAvailable)
	}
	size := attr.GetSize()
	if bytesAvailable < size {
		return 0, NewMessageTruncatedError(fmt.Sprintf("not enough space for attribute: %v", attr.Name))
	}
	bytes, err := b.AppendBytes(size)
	if err != nil {
		return 0, err
	}
	switch size {
	default:
		copy(bytes, value.([]byte))
	case 1:
		switch value.(type) {
		case int:
			bytes[0] = byte(value.(int))
		default:
			bytes[0] = value.(byte)
		}
	case 2:
		switch value.(type) {
		case int:
			binary.BigEndian.PutUint16(bytes, uint16(value.(int)))
		default:
			binary.BigEndian.PutUint16(bytes, value.(uint16))
		}
	case 4:
		switch value.(type) {
		case int:
			binary.BigEndian.PutUint32(bytes, uint32(value.(int)))
		default:
			binary.BigEndian.PutUint32(bytes, value.(uint32))
		}
	case 8:
		switch value.(type) {
		case int:
			binary.BigEndian.PutUint64(bytes, uint64(value.(int)))
		default:
			binary.BigEndian.PutUint64(bytes, value.(uint64))
		}
	}
	return size, nil
}

// BufferToTableAttributes takes the reconstructed octet buffer transmitted for
// a table attribute (over many GetNextResponses) and converts it into the desired
// format for each table row
func (attr *AttributeDefinition) BufferToTableAttributes(data []byte) (interface{}, error) {
	// Source is network byte order octets. Convert to proper array of slices
	rowSize := attr.GetSize()
	dataSize := len(data)
	index := 0

	switch rowSize {
	default:
		value := make([][]byte, dataSize/rowSize)
		for offset := 0; offset < dataSize; offset += rowSize {
			value[index] = make([]byte, rowSize)
			copy(value[index], data[offset:])
			index++
		}
		return value, nil
	case 1:
		value := make([]byte, dataSize)
		copy(value, data)
		return value, nil
	case 2:
		value := make([]uint16, dataSize/2)
		for offset := 0; offset < dataSize; offset += rowSize {
			value[offset] = binary.BigEndian.Uint16(data[offset:])
			index++
		}
		return value, nil
	case 4:
		value := make([]uint32, dataSize/4)
		for offset := 0; offset < dataSize; offset += rowSize {
			value[offset] = binary.BigEndian.Uint32(data[offset:])
			index++
		}
		return value, nil
	case 8:
		value := make([]uint64, dataSize/8)
		for offset := 0; offset < dataSize; offset += rowSize {
			value[offset] = binary.BigEndian.Uint64(data[offset:])
			index++
		}
		return value, nil
	}
}

func (attr *AttributeDefinition) tableAttributeDecode(data []byte, df gopacket.DecodeFeedback, msgType byte) (interface{}, error) {
	// Serialization of a table depends on the type of message. A
	// Review of ITU-T G.988 shows that access on tables are
	// either Read and/or Write, never Set-by-Create
	switch msgType {
	default:
		return nil, errors.New(fmt.Sprintf("unsupported Message Type '%v' for table serialization", msgType))

	case byte(Get) | AK: // Get Response
		// Size
		value := binary.BigEndian.Uint32(data[0:4])
		return value, nil

	case byte(GetNext) | AK: // Get Next Response
		// Block of data (octets) that need to be reassembled before conversion
		// to table/row-data.  If table attribute is not explicitly given a value
		// we have to assume the entire data buffer is the value. The receiver of
		// this frame will need to trim off any addtional information at the end
		// of the last frame sequence since they (and the ONU) are the only ones
		// who know how long the data really is.
		size := attr.GetSize()
		if size != 0 && len(data) < attr.GetSize() {
			df.SetTruncated()
			return nil, NewMessageTruncatedError("packet too small for field")
		} else if size == 0 {
			return nil, NewProcessingError("table attributes with no size are not supported: %v", attr.Name)
		}
		return data, nil

	case byte(Set) | AR: // Set Request
		fmt.Println("TODO")
		return nil, errors.New("TODO")

	case byte(SetTable) | AR: // Set Table Request
		// TODO: Only baseline supported at this time
		return nil, errors.New("attribute encode for set-table-request not yet supported")
	}
	return nil, errors.New("TODO")
}

func (attr *AttributeDefinition) tableAttributeSerializeTo(value interface{}, b gopacket.SerializeBuffer, msgType byte,
	bytesAvailable int) (int, error) {
	// Serialization of a table depends on the type of message. A
	// Review of ITU-T G.988 shows that access on tables are
	// either Read and/or Write, never Set-by-Create
	switch msgType {
	default:
		return 0, errors.New(fmt.Sprintf("unsupported Message Type '%v' for table serialization", msgType))

	case byte(Get) | AK: // Get Response
		// Size
		if bytesAvailable < 4 {
			return 0, NewMessageTruncatedError(fmt.Sprintf("not enough space for attribute: %v", attr.Name))
		}
		if dwordSize, ok := value.(uint32); ok {
			bytes, err := b.AppendBytes(4)
			if err != nil {
				return 0, err
			}
			binary.BigEndian.PutUint32(bytes, dwordSize)
			return 4, nil
		}
		return 0, errors.New("unexpected type for table serialization")

	case byte(GetNext) | AK: // Get Next Response
		// Values are already in network by order form
		if data, ok := value.([]byte); ok {
			if bytesAvailable < len(data) {
				return 0, NewMessageTruncatedError(fmt.Sprintf("not enough space for attribute: %v", attr.Name))
			}
			bytes, err := b.AppendBytes(len(data))
			if err != nil {
				return 0, err
			}
			copy(bytes, data)
			return len(data), nil
		}
		return 0, errors.New("unexpected type for table serialization")

	case byte(Set) | AR: // Set Request
		fmt.Println("TODO")

	case byte(SetTable) | AR: // Set Table Request
		// TODO: Only baseline supported at this time
		return 0, errors.New("attribute encode for set-table-request not yet supported")
	}
	size := attr.GetSize()
	if bytesAvailable < size {
		return 0, NewMessageTruncatedError(fmt.Sprintf("not enough space for attribute: %v", attr.Name))
	}
	bytes, err := b.AppendBytes(size)
	if err != nil {
		return 0, err
	}
	switch attr.GetSize() {
	default:
		copy(bytes, value.([]byte))
	case 1:
		switch value.(type) {
		case int:
			bytes[0] = byte(value.(int))
		default:
			bytes[0] = value.(byte)
		}
	case 2:
		switch value.(type) {
		case int:
			binary.BigEndian.PutUint16(bytes, uint16(value.(int)))
		default:
			binary.BigEndian.PutUint16(bytes, value.(uint16))
		}
	case 4:
		switch value.(type) {
		case int:
			binary.BigEndian.PutUint32(bytes, uint32(value.(int)))
		default:
			binary.BigEndian.PutUint32(bytes, value.(uint32))
		}
	case 8:
		switch value.(type) {
		case int:
			binary.BigEndian.PutUint64(bytes, uint64(value.(int)))
		default:
			binary.BigEndian.PutUint64(bytes, value.(uint64))
		}
	}
	return size, nil
}

// GetAttributeDefinitionByName searches the attribute definition map for the
// attribute with the specified name (case insensitive)
func GetAttributeDefinitionByName(attrMap *AttributeDefinitionMap, name string) (*AttributeDefinition, OmciErrors) {
	nameLower := strings.ToLower(name)
	for _, attrVal := range *attrMap {
		if nameLower == strings.ToLower(attrVal.GetName()) {
			return attrVal, nil
		}
	}
	return nil, NewAttributeFailureError(fmt.Sprintf("attribute '%s' not found", name))
}

// GetAttributeDefinitionMapKeys is a convenience functions since we may need to
// iterate a map in key index order. Maps in Go since v1.0 the iteration order
// of maps have been randomized.
func GetAttributeDefinitionMapKeys(attrMap AttributeDefinitionMap) []uint {
	var keys []uint
	for k := range attrMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

// GetAttributeBitmap is a convenience functions to scan a list of attributes
// and return the bitmask that represents them
func GetAttributeBitmap(attrMap AttributeDefinitionMap, attributes mapset.Set) (uint16, error) {
	var mask uint16
	for k, def := range attrMap {
		if attributes.Contains(def.Name) {
			mask |= 1 << uint16(16-k)
			attributes.Remove(def.Name)
		}
	}
	if attributes.Cardinality() > 0 {
		return 0, errors.New(fmt.Sprintf("unsupported attributes: %v", attributes))
	}
	return mask, nil
}

///////////////////////////////////////////////////////////////////////
// Packet definitions for attributes of various types/sizes

func ByteField(name string, defVal uint8, access mapset.Set, avc bool,
	counter bool, optional bool, deprecated bool, index uint) *AttributeDefinition {
	return &AttributeDefinition{
		Name:         name,
		Index:        index,
		DefValue:     defVal,
		Size:         1,
		Access:       access,
		Avc:          avc,
		Counter:      counter,
		TableSupport: false,
		Optional:     optional,
		Deprecated:   deprecated,
	}
}

func Uint16Field(name string, defVal uint16, access mapset.Set, avc bool,
	counter bool, optional bool, deprecated bool, index uint) *AttributeDefinition {
	return &AttributeDefinition{
		Name:         name,
		Index:        index,
		DefValue:     defVal,
		Size:         2,
		Access:       access,
		Avc:          avc,
		Counter:      counter,
		TableSupport: false,
		Optional:     optional,
		Deprecated:   deprecated,
	}
}

func Uint32Field(name string, defVal uint32, access mapset.Set, avc bool,
	counter bool, optional bool, deprecated bool, index uint) *AttributeDefinition {
	return &AttributeDefinition{
		Name:         name,
		Index:        index,
		DefValue:     defVal,
		Size:         4,
		Access:       access,
		Avc:          avc,
		Counter:      counter,
		TableSupport: false,
		Optional:     optional,
		Deprecated:   deprecated,
	}
}

func Uint64Field(name string, defVal uint64, access mapset.Set, avc bool,
	counter bool, optional bool, deprecated bool, index uint) *AttributeDefinition {
	return &AttributeDefinition{
		Name:         name,
		Index:        index,
		DefValue:     defVal,
		Size:         8,
		Access:       access,
		Avc:          avc,
		Counter:      counter,
		TableSupport: false,
		Optional:     optional,
		Deprecated:   deprecated,
	}
}

func MultiByteField(name string, size uint, defVal []byte, access mapset.Set, avc bool,
	counter bool, optional bool, deprecated bool, index uint) *AttributeDefinition {
	return &AttributeDefinition{
		Name:         name,
		Index:        index,
		DefValue:     defVal,
		Size:         int(size),
		Access:       access,
		Avc:          avc,
		Counter:      counter,
		TableSupport: false,
		Optional:     optional,
		Deprecated:   deprecated,
	}
}

// Notes on various OMCI ME Table attribute fields.  This comment will eventually be
// removed once a good table solution is implemented.  These are not all the MEs with
// table attributes, but probably ones I care about to support initially.
//
//   ME                     Notes
//  --------------------------------------------------------------------------------------------
//	Port-mapping package -> Combined Port table -> N * 25 sized rows (port (1) + ME(2) * 12)
//  ONU Remote Debug     -> Reply table (N bytes)
//  ONU3-G               -> Status snapshot recordtable M x N bytes
//  MCAST Gem interworkTP-> IPv4 multicast adress table (12*n) (two 2 byte fields, two 4 byte fields)
//                          IPv6 multicast adress table (24*n) (various sub-fields)
//  L2 mcast gem TP      -> MCAST MAC addr filtering table (11 * n) (various sub-fields)
//  MAC Bridge Port Filt -> MAC Filter table (8 * n) (3 fields, some are bits)      *** BITS ***
//  MAC Bridge Port data -> Bridge Table (8*M) (vaius fields, some are bits)        *** BITS ***
//  VLAN tagging filter  -> Rx Vlan tag op table (16 * n) Lots of bit fields        *** BITS ***
//  MCAST operations profile
//  MCAST Subscriber config info
//  MCAST subscriber monitor
//  OMCI                -> Two tables (N bytes and 2*N bytes)
//  General pupose buffer   -> N bytes
//  Enhanced security control (17 * N bytes), (16 * P Bytes) , (16 * Q bytes), and more...
//
// An early example of info to track
//
type TableInfo struct {
	Value interface{}
	Size  int
}

func (t *TableInfo) String() string {
	return fmt.Sprintf("TableInfo: Size: %d, Value(s): %v", t.Size, t.Value)
}

// Now the field
func TableField(name string, tableInfo TableInfo, access mapset.Set,
	avc bool, optional bool, deprecated bool, index uint) *AttributeDefinition {
	return &AttributeDefinition{
		Name:         name,
		Index:        index,
		DefValue:     tableInfo.Value,
		Size:         tableInfo.Size, //Number of elements
		Access:       access,
		Avc:          avc,
		Counter:      false,
		TableSupport: true,
		Optional:     optional,
		Deprecated:   deprecated,
	}
}

func UnknownField(name string, defVal uint64, access mapset.Set, avc bool,
	counter bool, optional bool, deprecated bool, index uint) *AttributeDefinition {
	return &AttributeDefinition{
		Name:         name,
		Index:        index,
		DefValue:     defVal,
		Size:         99999999,
		Access:       access,
		Avc:          avc,
		Counter:      counter,
		TableSupport: false,
		Optional:     optional,
		Deprecated:   deprecated,
	}
}

///////////////////////////////////////////////////////////////////////
// Attribute Name to Value    (Interfaced defined in generated subdirectory)

type AttributeValueMap map[string]interface{}
