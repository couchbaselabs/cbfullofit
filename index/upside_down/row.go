//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package upside_down

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/couchbaselabs/cbfullofit/index"
)

const BYTE_SEPARATOR byte = 0xff

type UpsideDownCouchRowStream chan UpsideDownCouchRow

type UpsideDownCouchRow interface {
	Key() []byte
	Value() []byte
}

func ParseFromKeyValue(key, value []byte) UpsideDownCouchRow {
	switch key[0] {
	case 'v':
		return NewVersionRowKV(key, value)
	case 'f':
		return NewFieldRowKV(key, value)
	// case 'i':
	// 	return NewInverseFrequencyRowKV(key, value)
	case 't':
		return NewTermFrequencyRowKV(key, value)
	// case 'n':
	// 	return NewNormalizationRowKV(key, value)
	case 'b':
		return NewBackIndexRowKV(key, value)
	}
	return nil
}

// VERSION

type VersionRow struct {
	version uint8
}

func (v *VersionRow) Key() []byte {
	return []byte{'v'}
}

func (v *VersionRow) Value() []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, v.version)
	if err != nil {
		panic(fmt.Sprintf("binary.Write failed: %v", err))
	}
	return buf.Bytes()
}

func (v *VersionRow) String() string {
	return fmt.Sprintf("Version: %d", v.version)
}

func NewVersionRow(version uint8) *VersionRow {
	return &VersionRow{
		version: version,
	}
}

func NewVersionRowKV(key, value []byte) *VersionRow {
	rv := VersionRow{}
	buf := bytes.NewBuffer(value)
	err := binary.Read(buf, binary.LittleEndian, &rv.version)
	if err != nil {
		panic(fmt.Sprintf("binary.Read failed: %v", err))
	}
	return &rv
}

// FIELD definition

type FieldRow struct {
	index    uint16
	name     string
	path     string
	analyzer string
}

func (f *FieldRow) Key() []byte {
	buf := new(bytes.Buffer)
	err := buf.WriteByte('f')
	if err != nil {
		panic(fmt.Sprintf("Buffer.WriteByte failed: %v", err))
	}
	err = binary.Write(buf, binary.LittleEndian, f.index)
	if err != nil {
		panic(fmt.Sprintf("binary.Write failed: %v", err))
	}
	return buf.Bytes()
}

func (f *FieldRow) Value() []byte {
	buf := new(bytes.Buffer)
	_, err := buf.WriteString(f.name)
	if err != nil {
		panic(fmt.Sprintf("Buffer.WriteString failed: %v", err))
	}
	err = buf.WriteByte(BYTE_SEPARATOR)
	if err != nil {
		panic(fmt.Sprintf("Buffer.WriteByte failed: %v", err))
	}
	_, err = buf.WriteString(f.path)
	if err != nil {
		panic(fmt.Sprintf("Buffer.WriteString failed: %v", err))
	}
	err = buf.WriteByte(BYTE_SEPARATOR)
	if err != nil {
		panic(fmt.Sprintf("Buffer.WriteByte failed: %v", err))
	}
	_, err = buf.WriteString(f.analyzer)
	if err != nil {
		panic(fmt.Sprintf("Buffer.WriteString failed: %v", err))
	}
	return buf.Bytes()
}

func (f *FieldRow) Field() *index.Field {
	return &index.Field{
		Name:     f.name,
		Path:     f.path,
		Analyzer: f.analyzer,
	}
}

func (f *FieldRow) String() string {
	return fmt.Sprintf("Field: %d Name: %s Path: %s Analyzer: %s", f.index, f.name, f.path, f.analyzer)
}

func NewFieldRow(index uint16, name, path, analyzer string) *FieldRow {
	return &FieldRow{
		index:    index,
		name:     name,
		path:     path,
		analyzer: analyzer,
	}
}

func NewFieldRowKV(key, value []byte) *FieldRow {
	rv := FieldRow{}

	buf := bytes.NewBuffer(key)
	buf.ReadByte() // type
	err := binary.Read(buf, binary.LittleEndian, &rv.index)
	if err != nil {
		panic(fmt.Sprintf("binary.Read failed: %v", err))
	}

	buf = bytes.NewBuffer(value)
	rv.name, err = buf.ReadString(BYTE_SEPARATOR)
	if err != nil {
		panic(fmt.Sprintf("Buffer.ReadString failed: %v", err))
	}
	rv.name = rv.name[:len(rv.name)-1] // trim off separator byte
	rv.path, err = buf.ReadString(BYTE_SEPARATOR)
	if err != nil {
		panic(fmt.Sprintf("Buffer.ReadString failed: %v", err))
	}
	rv.path = rv.path[:len(rv.path)-1] // trim off separator byte
	rv.analyzer, err = buf.ReadString(BYTE_SEPARATOR)
	if err != io.EOF {
		panic(fmt.Sprintf("expected binary.ReadString to end in EOF: %v", err))
	}
	return &rv
}

// TERM FIELD FREQUENCY

type TermFrequencyRow struct {
	term  []byte
	field uint16
	doc   []byte
	freq  uint64
	norm  float32
}

func (tfr *TermFrequencyRow) Key() []byte {
	buf := new(bytes.Buffer)
	err := buf.WriteByte('t')
	if err != nil {
		panic(fmt.Sprintf("Buffer.WriteByte failed: %v", err))
	}
	_, err = buf.Write(tfr.term)
	if err != nil {
		panic(fmt.Sprintf("Buffer.Write failed: %v", err))
	}
	err = buf.WriteByte(BYTE_SEPARATOR)
	if err != nil {
		panic(fmt.Sprintf("Buffer.WriteByte failed: %v", err))
	}
	err = binary.Write(buf, binary.LittleEndian, tfr.field)
	if err != nil {
		panic(fmt.Sprintf("binary.Write failed: %v", err))
	}
	_, err = buf.Write(tfr.doc)
	if err != nil {
		panic(fmt.Sprintf("Buffer.Write failed: %v", err))
	}
	return buf.Bytes()
}

func (tfr *TermFrequencyRow) Value() []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, tfr.freq)
	if err != nil {
		panic(fmt.Sprintf("binary.Write failed: %v", err))
	}
	err = binary.Write(buf, binary.LittleEndian, tfr.norm)
	if err != nil {
		panic(fmt.Sprintf("binary.Write failed: %v", err))
	}
	return buf.Bytes()
}

func (tfr *TermFrequencyRow) String() string {
	return fmt.Sprintf("Term: `%s` Field: %d DocId: `%s` Frequency: %d Norm: %f", string(tfr.term), tfr.field, string(tfr.doc), tfr.freq, tfr.norm)
}

func NewTermFrequencyRow(term []byte, field uint16, doc []byte, freq uint64, norm float32) *TermFrequencyRow {
	return &TermFrequencyRow{
		term:  term,
		field: field,
		doc:   doc,
		freq:  freq,
		norm:  norm,
	}
}

func NewTermFrequencyRowKV(key, value []byte) *TermFrequencyRow {
	rv := TermFrequencyRow{}
	buf := bytes.NewBuffer(key)
	buf.ReadByte() // type

	var err error
	rv.term, err = buf.ReadBytes(BYTE_SEPARATOR)
	if err != nil {
		panic(fmt.Sprintf("Buffer.ReadString failed: %v", err))
	}
	rv.term = rv.term[:len(rv.term)-1] // trim off separator byte

	err = binary.Read(buf, binary.LittleEndian, &rv.field)
	if err != nil {
		panic(fmt.Sprintf("binary.Read failed: %v", err))
	}

	rv.doc, err = buf.ReadBytes(BYTE_SEPARATOR)
	if err != io.EOF {
		panic(fmt.Sprintf("expected binary.ReadString to end in EOF: %v", err))
	}

	buf = bytes.NewBuffer((value))
	err = binary.Read(buf, binary.LittleEndian, &rv.freq)
	if err != nil {
		panic(fmt.Sprintf("binary.Read failed: %v", err))
	}
	err = binary.Read(buf, binary.LittleEndian, &rv.norm)
	if err != nil {
		panic(fmt.Sprintf("binary.Read failed: %v", err))
	}

	return &rv

}

type BackIndexEntry struct {
	term  []byte
	field uint16
}

func (bie *BackIndexEntry) String() string {
	return fmt.Sprintf("Term: `%s` Field: %d", string(bie.term), bie.field)
}

type BackIndexRow struct {
	doc     []byte
	entries []*BackIndexEntry
}

func (br *BackIndexRow) Key() []byte {
	buf := new(bytes.Buffer)
	err := buf.WriteByte('b')
	if err != nil {
		panic(fmt.Sprintf("Buffer.WriteByte failed: %v", err))
	}
	err = binary.Write(buf, binary.LittleEndian, br.doc)
	if err != nil {
		panic(fmt.Sprintf("binary.Write failed: %v", err))
	}
	return buf.Bytes()
}

func (br *BackIndexRow) Value() []byte {
	buf := new(bytes.Buffer)
	for _, e := range br.entries {
		_, err := buf.Write(e.term)
		if err != nil {
			panic(fmt.Sprintf("Buffer.Write failed: %v", err))
		}
		err = buf.WriteByte(BYTE_SEPARATOR)
		if err != nil {
			panic(fmt.Sprintf("Buffer.WriteByte failed: %v", err))
		}
		err = binary.Write(buf, binary.LittleEndian, e.field)
		if err != nil {
			panic(fmt.Sprintf("binary.Write failed: %v", err))
		}
	}
	return buf.Bytes()
}

func (br *BackIndexRow) String() string {
	return fmt.Sprintf("Backindex DocId: `%s` Entries: %v", string(br.doc), br.entries)
}

func NewBackIndexRow(doc []byte, entries []*BackIndexEntry) *BackIndexRow {
	return &BackIndexRow{
		doc:     doc,
		entries: entries,
	}
}

func NewBackIndexRowKV(key, value []byte) *BackIndexRow {
	rv := BackIndexRow{}

	buf := bytes.NewBuffer(key)
	buf.ReadByte() // type

	var err error
	rv.doc, err = buf.ReadBytes(BYTE_SEPARATOR)
	if err != io.EOF {
		panic(fmt.Sprintf("expected binary.ReadString to end in EOF: %v", err))
	}

	buf = bytes.NewBuffer(value)
	rv.entries = make([]*BackIndexEntry, 0)

	var term []byte
	term, err = buf.ReadBytes(BYTE_SEPARATOR)
	if err != nil && err != io.EOF {
		panic(fmt.Sprintf("Buffer.ReadString failed: %v", err))
	}
	for err != io.EOF {
		ent := BackIndexEntry{}
		ent.term = term[:len(term)-1] // trim off separator byte

		err = binary.Read(buf, binary.LittleEndian, &ent.field)
		if err != nil {
			panic(fmt.Sprintf("binary.Read failed: %v", err))
		}
		rv.entries = append(rv.entries, &ent)

		term, err = buf.ReadBytes(BYTE_SEPARATOR)
		if err != nil && err != io.EOF {
			panic(fmt.Sprintf("Buffer.ReadString failed: %v", err))
		}
	}

	return &rv
}
