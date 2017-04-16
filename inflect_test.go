package inflect

import (
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestGet(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Hour)
	t3 := t2.Add(time.Hour)

	d := &Data{}
	d.Code = 66
	d.CreatedAt = t1
	d.UpdatedAt = t2
	d.DeletedAt = &t3

	{
		iv, err := Get(d, `Code`)
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		v, ok := iv.(int)
		if !ok {
			t.Error(`not an int`)
			t.Fail()
			return
		}
		if v != 66 {
			t.Error(`!= 66`)
			t.Fail()
			return
		}
	}

	{
		iv, err := Get(d, `UpdatedAt`)
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		v, ok := iv.(time.Time)
		if !ok {
			t.Error(`not an int`)
			t.Fail()
			return
		}
		if v != t2 {
			t.Error(`!=`, t2)
			t.Fail()
			return
		}
	}

	{
		iv, err := Get(d, `CreatedAt`)
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		v, ok := iv.(time.Time)
		if !ok {
			t.Error(`not an int`)
			t.Fail()
			return
		}
		if v != t1 {
			t.Error(`!=`, t1)
			t.Fail()
			return
		}
	}

	{
		iv, err := Get(d, `DeletedAt`)
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		v, ok := iv.(*time.Time)
		if !ok {
			t.Error(`not an int`)
			t.Fail()
			return
		}
		if *v != t3 {
			t.Error(`!=`, t3)
			t.Fail()
			return
		}
		if v != &t3 {
			t.Error(`!=`, &t3)
			t.Fail()
			return
		}
	}

	{
		_, err := Get(d, `DeletedAtLast`)
		if err != ErrNotFound {
			t.Error(`expected`, ErrNotFound)
			t.Fail()
			return
		}
	}
}

func TestSet(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Hour)
	t3 := t2.Add(time.Hour)
	code := 66
	name := `test`

	d := &Data{}
	d.Code = code
	d.Name = name
	d.CreatedAt = t1
	d.UpdatedAt = t2
	d.DeletedAt = &t3

	{
		newValue := `new test`
		err := Set(d, `Name`, newValue)
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		if d.Name != newValue {
			t.Fail()
		}
	}

	{
		newValue := t3.Add(time.Hour * 9)
		err := Set(d, `DeletedAt`, &newValue)
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		if *d.DeletedAt != newValue {
			t.Fail()
		}
	}

	{
		newValue := t1.Add(time.Hour * 9)
		err := Set(d, `CreatedAt`, newValue)
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		if d.CreatedAt != newValue {
			t.Fail()
		}
	}
}

func TestTag(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Hour)
	t3 := t2.Add(time.Hour)
	code := 66
	name := `test`

	d := &Data{}
	d.Code = code
	d.Name = name
	d.CreatedAt = t1
	d.UpdatedAt = t2
	d.DeletedAt = &t3

	{
		s, err := Tag(d, `ID`, `bson`)
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		if s != "_id" {
			t.Error(err)
			t.Fail()
			return
		}
	}

	{
		s, err := Tag(d, `Model`, `bson`)
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		if s != ",inline" {
			t.Error(err)
			t.Fail()
			return
		}
	}

	{
		_, err := Tag(d, `Model`, `ison`)
		if err != ErrNotFound {
			t.Error(err)
			t.Fail()
			return
		}
	}

	{
		_, err := Tag(d, `ModelBase`, `json`)
		if err != ErrNotFound {
			t.Error(err)
			t.Fail()
			return
		}
	}
}

func TestCAS(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Hour)
	t3 := t2.Add(time.Hour)
	code := 66
	name := `test`

	d := &Data{}
	d.Code = code
	d.Name = name
	d.CreatedAt = t1
	d.UpdatedAt = t2
	d.DeletedAt = &t3

	{
		newValue := `new test`
		oldValue := name
		err := CAS(d, `Name`, oldValue, newValue)
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		if d.Name != newValue {
			t.Fail()
		}
	}

	{
		newValue := t3.Add(time.Hour * 9)
		oldValue := &t3
		err := CAS(d, `DeletedAt`, oldValue, &newValue)
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		if *d.DeletedAt != newValue {
			t.Fail()
		}
		if d.DeletedAt != &newValue {
			t.Fail()
		}
	}

	{
		newValue := t1.Add(time.Hour * 9)
		oldValue := t1
		err := CAS(d, `CreatedAt`, oldValue, newValue)
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		if d.CreatedAt != newValue {
			t.Fail()
		}
	}

	{
		newValue := t1.Add(time.Hour * 9)
		oldValue := t1.Add(time.Minute)
		err := CAS(d, `CreatedAt`, oldValue, newValue)
		if err != ErrNotMatched {
			t.Error(err)
			t.Fail()
			return
		}
	}
}

//-----------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------

type Data struct {
	MetaModel `json:",inline"`
	CreatedAt time.Time

	Name string
	Code int
}

//-----------------------------------------------------------------------------
