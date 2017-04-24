package inflect

import (
	"fmt"
	"testing"
	"time"
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

func checkValue(t *testing.T, value, expected interface{}) {
	vType := fmt.Sprintf("%T", value)
	eType := fmt.Sprintf("%T", expected)
	if vType != eType {
		t.Fatal("expected type:", eType, "got:", vType)
	}
	if value != expected {
		t.Fatal("expected value:", expected, "got:", value)
	}
}

func TestGet1(t *testing.T) {
	d := data()
	iv, err := Get(d, `Code`)
	if err != nil {
		t.Fatal(err)
	}
	checkValue(t, iv, code)
}

func TestGet2(t *testing.T) {
	d := data()
	iv, err := Get(d, `UpdatedAt`)
	if err != nil {
		t.Fatal(err)
	}
	checkValue(t, iv, t2)
}

func TestGet3(t *testing.T) {
	d := data()
	iv, err := Get(d, `CreatedAt`)
	if err != nil {
		t.Fatal(err)
	}
	checkValue(t, iv, t1)
}

func TestGet4(t *testing.T) {
	d := data()
	iv, err := Get(d, `DeletedAt`)
	if err != nil {
		t.Fatal(err)
	}
	checkValue(t, iv, &t3)
	v, _ := iv.(*time.Time)
	checkValue(t, *v, t3)
}

func TestGet5(t *testing.T) {
	d := data()
	_, err := Get(d, `DeletedAtLast`)
	if err != ErrNotFound {
		t.Fatal(`expected error:`, ErrNotFound)
	}
}

func TestSet1(t *testing.T) {
	d := data()
	newValue := `new test`
	err := Set(d, `Name`, newValue)
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != newValue {
		t.Fatal()
	}
}

func TestSet2(t *testing.T) {
	d := data()
	newValue := t3.Add(time.Hour * 9)
	err := Set(d, `DeletedAt`, &newValue)
	if err != nil {
		t.Fatal(err)
	}
	if *d.DeletedAt != newValue {
		t.Fatal()
	}
}

func TestSet3(t *testing.T) {
	d := data()
	newValue := t1.Add(time.Hour * 9)
	err := Set(d, `CreatedAt`, newValue)
	if err != nil {
		t.Fatal(err)
	}
	if d.CreatedAt != newValue {
		t.Fatal()
	}
}

func TestTag1(t *testing.T) {
	d := data()
	s, err := Tag(d, `ID`, `bson`)
	if err != nil {
		t.Fatal(err)
	}
	if s != "_id" {
		t.Fatal("expected:", "_id", "got:", s)
	}
}

func TestTag2(t *testing.T) {
	d := data()
	s, err := Tag(d, `Model`, `bson`)
	if err != nil {
		t.Fatal(err)
	}
	if s != ",inline" {
		t.Fatal("expected:", ",inline", "got:", s)
	}
}

func TestTag3(t *testing.T) {
	d := data()
	_, err := Tag(d, `Model`, `ison`)
	if err != ErrNotFound {
		t.Fatal(err)
	}
}

func TestTag4(t *testing.T) {
	d := data()
	_, err := Tag(d, `ModelBase`, `json`)
	if err != ErrNotFound {
		t.Fatal(err)
	}
}

func TestCAS1(t *testing.T) {
	d := data()
	newValue := `new test`
	oldValue := name
	err := CAS(d, `Name`, oldValue, newValue)
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != newValue {
		t.Fatal()
	}
}

func TestCAS2(t *testing.T) {
	d := data()
	newValue := t3.Add(time.Hour * 9)
	oldValue := &t3
	err := CAS(d, `DeletedAt`, oldValue, &newValue)
	if err != nil {
		t.Fatal(err)
	}
	if *d.DeletedAt != newValue {
		t.Fatal()
	}
	if d.DeletedAt != &newValue {
		t.Fatal()
	}
}

func TestCAS3(t *testing.T) {
	d := data()
	newValue := t1.Add(time.Hour * 9)
	oldValue := t1
	err := CAS(d, `CreatedAt`, oldValue, newValue)
	if err != nil {
		t.Fatal(err)
	}
	if d.CreatedAt != newValue {
		t.Fatal()
	}
}

func TestCAS4(t *testing.T) {
	d := data()
	newValue := t1.Add(time.Hour * 9)
	oldValue := t1.Add(time.Minute)
	err := CAS(d, `CreatedAt`, oldValue, newValue)
	if err != ErrNotMatched {
		t.Fatal(err)
	}
}
