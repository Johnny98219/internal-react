// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package immutable

import (
	"go/ast"
	"go/types"
	"strings"
	"sync"

	"golang.org/x/tools/go/types/typeutil"
)

type ImmType interface {
	isImmType()
}

type (
	ImmTypeUnknown struct{}
	ImmTypeStruct  struct{}
	ImmTypeMap     struct {
		Key  types.Type
		Elem types.Type
	}
	ImmTypeSlice struct {
		Elem types.Type
	}
)

func (i ImmTypeUnknown) isImmType() {}
func (i ImmTypeStruct) isImmType()  {}
func (i ImmTypeMap) isImmType()     {}
func (i ImmTypeSlice) isImmType()   {}

const (
	// ImmTypeTmplPrefix is the prefix used to identify immutable type templates
	ImmTypeTmplPrefix = "_Imm_"

	// Pkg is the import path of this package
	PkgImportPath = "github.com/myitcv/immutable"
)

// Immutable is the interface implemented by all immutable types. If Go had generics the interface would
// be defined, assuming a generic type parameter T, as follows:
//
// 	type Immutable<T> interface {
// 		AsMutable() T
// 		AsImmutable() T
// 		WithMutations(f func(v T)) T
// 		Mutable() bool
// 	}
//
// Because we do not have such a type parameter we can only define the Mutable() method in the interface
type Immutable interface {
	Mutable() bool
}

// IsImmTmpl determines whether the supplied type spec is an immutable template type (either a struct,
// slice or map), returning the name of the type with the ImmTypeTmplPrefix removed in that case
func IsImmTmpl(ts *ast.TypeSpec) (string, bool) {
	typName := ts.Name.Name

	if !strings.HasPrefix(typName, ImmTypeTmplPrefix) {
		return "", false
	}

	valid := false

	switch typ := ts.Type.(type) {
	case *ast.MapType:
		valid = true
	case *ast.ArrayType:
		if typ.Len == nil {
			valid = true
		}
	case *ast.StructType:
		valid = true
	}

	if !valid {
		return "", false
	}

	name := strings.TrimPrefix(typName, ImmTypeTmplPrefix)

	return name, true
}

var ic immCache

type immCache struct {
	mu      sync.Mutex
	msCache typeutil.MethodSetCache
	res     map[*types.Named]ImmType
}

func (i *immCache) lookup(tt types.Type) ImmType {
	pt, ok := tt.(*types.Pointer)
	if !ok {
		return nil
	}

	nt, ok := pt.Elem().(*types.Named)
	if !ok {
		return nil
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	v, ok := i.res[nt]
	if ok {
		return v
	}

	if i.res == nil {
		i.res = make(map[*types.Named]ImmType)
	}

	ms := i.msCache.MethodSet(pt)

	foundMutable := false
	foundAsMutable := false
	foundAsImmutable := false
	foundWithMutable := false
	foundWithImmutable := false

	for i := 0; i < ms.Len(); i++ {
		f := ms.At(i).Obj().(*types.Func)
		t := f.Type().(*types.Signature)

		switch mn := f.Name(); mn {
		case "Mutable":
			if t.Params().Len() != 0 {
				break
			}

			if t.Results().Len() != 1 {
				break
			}

			tres := t.Results().At(0)

			if b, ok := tres.Type().(*types.Basic); ok {
				foundMutable = b.Kind() == types.Bool
			}
		case "AsMutable":
			if t.Params().Len() != 0 {
				break
			}

			if t.Results().Len() != 1 {
				break
			}

			foundAsMutable = isPtrToNamedTyp(t.Results().At(0).Type(), nt)

		case "AsImmutable":
			if t.Params().Len() != 1 {
				break
			}

			if !isPtrToNamedTyp(t.Params().At(0).Type(), nt) {
				break
			}

			if t.Results().Len() != 1 {
				break
			}

			foundAsImmutable = isPtrToNamedTyp(t.Results().At(0).Type(), nt)

		case "WithMutable", "WithImmutable":
			if t.Params().Len() != 1 {
				break
			}

			st, ok := t.Params().At(0).Type().(*types.Signature)
			if !ok {
				break
			}

			if st.Params().Len() != 1 {
				break
			}

			if !isPtrToNamedTyp(st.Params().At(0).Type(), nt) {
				break
			}

			if st.Results().Len() != 0 {
				break
			}

			if t.Results().Len() != 1 {
				break
			}

			valid := isPtrToNamedTyp(t.Results().At(0).Type(), nt)

			switch mn {
			case "WithMutable":
				foundWithMutable = valid
			case "WithImmutable":
				foundWithImmutable = valid
			}
		}

	}

	isImm := foundMutable && foundAsMutable && foundAsImmutable &&
		foundWithMutable && foundWithImmutable

	if !isImm {
		i.res[nt] = nil
		return nil
	}

	v = ImmTypeUnknown{}

	// now we work out whether it's a struct, slice of map... else
	// it's unknown to this package

	st, ok := nt.Underlying().(*types.Struct)
	if !ok {
		return v
	}

	hasTmpl := false

	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)

		switch f.Name() {
		case "__tmpl":
			hasTmpl = true
		case "theMap":
			m := f.Type().(*types.Map)

			v = ImmTypeMap{
				Key:  m.Key(),
				Elem: m.Elem(),
			}
		case "theSlice":
			s := f.Type().(*types.Slice)

			v = ImmTypeSlice{
				Elem: s.Elem(),
			}
		}
	}

	if v == (ImmTypeUnknown{}) && hasTmpl {
		v = ImmTypeStruct{}
	}

	i.res[nt] = v
	return v
}

func isPtrToNamedTyp(t types.Type, nt *types.Named) bool {
	pt, ok := t.(*types.Pointer)
	if !ok {
		return false
	}

	n, ok := pt.Elem().(*types.Named)
	if !ok {
		return false
	}

	return nt == n
}

// IsImmType determines whether the supplied type is an immutable type. In case a type is
// immutable, a value of type ImmTypeStruct, ImmTypeSlice or ImmTypeMap is returned. In case the type is
// immutable but neither of the aforementioned instances, ImmTypeUnknown is returned. If a type
// is not immutable then nil is returned
func IsImmType(t types.Type) ImmType {
	return ic.lookup(t)
}
