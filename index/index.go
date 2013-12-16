//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package index

import (
	"fmt"
)

type Index interface {
	Open() error
	Close()

	Update(id []byte, doc []byte) error
	Delete(id []byte) error

	TermFieldReader(term []byte, field string) (TermFieldReader, error)

	DocCount() uint64
}

type TermFieldDoc struct {
	ID   string
	Freq uint64
	Norm float64
}

type TermFieldReader interface {
	Next() (*TermFieldDoc, error)
	Advance(ID []byte) (*TermFieldDoc, error)
	Count() uint64
	Close()
}

type Field struct {
	Name     string
	Path     string
	Analyzer string
}

func (f *Field) String() string {
	return fmt.Sprintf("Field[name=%s, path=%s, analyzer=%s]", f.Name, f.Path, f.Analyzer)
}
