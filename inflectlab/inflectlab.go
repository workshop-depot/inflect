//Package inflectlab is an incubator for ideas/tools
package inflectlab

import (
	"reflect"
	"time"

	"github.com/dc0d/inflect"
)

// Field provides access to a field, if a struct, it's children
// and it's metadata/tags
type Field struct {
	Ptr  *reflect.Value
	Tags map[string]string

	Children map[string]Field
}

func tags(s string) map[string]string {
	res := make(map[string]string)

	insideTag := false

	for s != "" {
		switch {
		case !insideTag:
			i := 0
			for i < len(s) && s[i] == ' ' {
				i++
			}
			s = s[i:]
			if s == "" {
				break
			}
			insideTag = true
		case insideTag:
			i := 0
			for i < len(s) && s[i] != ':' && s[i] != ' ' {
				i++
			}
			key := s[:i]
			switch s[i] {
			case ' ':
				res[key] = ""
				insideTag = false
				s = s[i:]
			case ':':
				s = s[(i + 1):]
				if s == "" {
					break
				}
				if s[0] != '"' {
					panic(`expected double quote "`)
				}
				s = s[1:]
				i := 0
				for i < len(s) && s[i] != '"' {
					i++
				}
				res[key] = s[:i]
				insideTag = false
				s = s[i+1:]
			}
		}
	}

	return res
}

func numField(ptr *reflect.Value) (res int) {
	defer func() {
		if e := recover(); e != nil {
			res = -1
		}
	}()
	res = ptr.NumField()
	return
}

func ptrOf(data interface{}) *reflect.Value {
	var val reflect.Value

	switch reflect.TypeOf(data).Kind() {
	case reflect.Ptr:
		val = reflect.ValueOf(data).Elem()
	default:
		return nil
	}

	return &val
}

// GetFields extracts field info from a pointer to a struct
func GetFields(data interface{}) (map[string]Field, error) {
	var ptr *reflect.Value
	switch x := data.(type) {
	case *reflect.Value:
		ptr = x
	default:
		ptr = ptrOf(data)
		if ptr == nil {
			return nil, inflect.ErrNonPointer
		}
	}

	typ := ptr.Type()
	res := make(map[string]Field)
	for i := 0; i < numField(ptr); i++ {
		fval := ptr.Field(i)
		ftyp := typ.Field(i)

		f := Field{
			Ptr:  &fval,
			Tags: tags(string(ftyp.Tag)),
		}

		switch fval.Kind() {
		case reflect.Invalid,
			reflect.Bool,
			reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
			reflect.Uintptr,
			reflect.Float32,
			reflect.Float64,
			reflect.Complex64,
			reflect.Complex128,
			reflect.Array,
			reflect.Chan,
			reflect.Func,
			reflect.Interface,
			reflect.Map,
			// reflect.Ptr,
			reflect.Slice,
			reflect.String,
			// reflect.Struct,
			reflect.UnsafePointer:
		default:
			switch ftyp.Type {
			case reflect.TypeOf(time.Time{}):
			default:
				cl, err := GetFields(&fval)
				if err != nil {
					// TODO:
					// fmt.Printf("ERR %v\n", err)
				} else {
					f.Children = cl
				}
			}
		}

		res[ftyp.Name] = f
	}
	return res, nil
}
