package inflectlab

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//-----------------------------------------------------------------------------
// test specific types

type Model struct {
	// ID bson.ObjectId `json:"id" bson:"_id"`

	ID        string `json:"id" bson:"_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `json:",omitempty" bson:",omitempty"`
}

// MetaData .
type MetaData map[string]interface{}

// MetaModel .
type MetaModel struct {
	Model    `json:",inline" bson:",inline"`
	MetaData `json:",omitempty" bson:",omitempty"`
}

type Data struct {
	MetaModel `json:",inline"`
	CreatedAt time.Time

	Name string
	Code int
}

//-----------------------------------------------------------------------------

var (
	t1   = time.Now()
	t2   = t1.Add(time.Hour)
	t3   = t2.Add(time.Hour)
	code = 66
	name = "test"
)

func data() *Data {
	d := &Data{
		Code: code,
		Name: name,
	}
	d.CreatedAt = t1
	d.UpdatedAt = t2
	d.DeletedAt = &t3
	return d
}

func Test01(t *testing.T) {
	d := data()
	json.Marshal(d)
	fields, err := GetFields(d)
	if err != nil {
		t.Fatal(err)
	}

	for k := range fields {
		switch k {
		case "MetaModel":
			for k1 := range fields[k].Children {
				switch k1 {
				case "Model":
					fields1 := fields[k].Children[k1].Children
					for k2 := range fields1 {
						switch k2 {
						case "CreatedAt", "UpdatedAt", "DeletedAt", "ID":
						default:
							t.Failed()
						}
					}
				case "MetaData":
					assert.Equal(t, 0, len(fields[k].Children[k1].Children))
				default:
					t.Failed()
				}
			}
		case "CreatedAt", "Name", "Code":
			assert.Equal(t, 0, len(fields[k].Children))
		default:
			t.Failed()
		}
	}

	nid := "IDNEW"
	fields["MetaModel"].Children["Model"].Children["ID"].Ptr.Set(reflect.ValueOf(nid))
	assert.Equal(t, nid, d.ID)
	assert.Equal(t, nid, fields["MetaModel"].Children["Model"].Children["ID"].Ptr.Interface())
}
