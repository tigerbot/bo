package util_test

import (
	"reflect"
	"testing"

	. "util"
)

type testStruct struct {
	Str        string
	Boolean    bool
	Integer    int
	SmallInt   int32
	BigInt     int64
	SmallFloat float32
	BigFloat   float64

	StrMap     map[string]string
	IntSlice   []int
	FloatArray [3]float64

	StructPtr *testStruct
	FloatPtr  *float64
}

func TestPtrCopy(t *testing.T) {
	var floatVal float64 = 123909834.3134
	val := testStruct{
		Str:        "hello",
		Boolean:    true,
		Integer:    9843,
		SmallInt:   98,
		BigInt:     5646573214365,
		SmallFloat: 654.1658,
		BigFloat:   23454773215.21345841,

		StrMap: map[string]string{
			"teal":      "bluish",
			"primaries": "red, green, blue",
			"purple":    "red+blue",
		},
		IntSlice:   []int{2, 5, 4864, 24657, 14},
		FloatArray: [3]float64{89746.4657, 65777, 32146694238.21657},

		StructPtr: &testStruct{
			Str:      "ptr struct",
			BigInt:   987564987654,
			IntSlice: []int{287, 75, 84, 2657, 1874},
		},
		FloatPtr: &floatVal,
	}

	cp, err := Copy(&val)
	if err != nil {
		t.Fatalf("failed to copy test struct: %v", err)
	}

	typed, ok := cp.(*testStruct)
	if !ok {
		t.Fatalf("copy result isn't expected type: %T", cp)
	}

	if !reflect.DeepEqual(&val, cp) {
		t.Fatalf("copy doesn't match the original value:\n%+v\n%+v", val, cp)
	}

	typed.Str = "this is now my own"
	if reflect.DeepEqual(&val, cp) {
		t.Error("changing the copy's string didn't break equality")
	}
	typed.Str = val.Str
	if !reflect.DeepEqual(&val, cp) {
		t.Fatal("Couldn't restored equality")
	}

	typed.BigFloat = 1290.33
	if reflect.DeepEqual(&val, cp) {
		t.Error("changing the copy's float64 didn't break equality")
	}
	typed.BigFloat = val.BigFloat
	if !reflect.DeepEqual(&val, cp) {
		t.Fatal("Couldn't restored equality")
	}

	typed.SmallInt = 547
	if reflect.DeepEqual(&val, cp) {
		t.Error("changing the copy's int32 didn't break equality")
	}
	typed.SmallInt = val.SmallInt
	if !reflect.DeepEqual(&val, cp) {
		t.Fatal("Couldn't restored equality")
	}

	typed.StrMap["different"] = "my reference has been manipulated"
	if reflect.DeepEqual(&val, cp) {
		t.Error("adding to the copy's map didn't break equality")
	}
	delete(typed.StrMap, "different")
	if !reflect.DeepEqual(&val, cp) {
		t.Fatal("Couldn't restored equality")
	}

	typed.StrMap["teal"] = "greenish"
	if reflect.DeepEqual(&val, cp) {
		t.Error("changing the copy's map didn't break equality")
	}
	typed.StrMap["teal"] = val.StrMap["teal"]
	if !reflect.DeepEqual(&val, cp) {
		t.Fatal("Couldn't restored equality")
	}

	typed.IntSlice[3] = 3
	if reflect.DeepEqual(&val, cp) {
		t.Error("changing the copy's slice didn't break equality")
	}
	typed.IntSlice[3] = val.IntSlice[3]
	if !reflect.DeepEqual(&val, cp) {
		t.Fatal("Couldn't restored equality")
	}

	typed.FloatArray[2] = 789465.4657
	if reflect.DeepEqual(&val, cp) {
		t.Error("changing the copy's slice didn't break equality")
	}
	typed.FloatArray[2] = val.FloatArray[2]
	if !reflect.DeepEqual(&val, cp) {
		t.Fatal("Couldn't restored equality")
	}

	if val.StructPtr != nil {
		typed.StructPtr.Str = "other test value"
		if reflect.DeepEqual(&val, cp) {
			t.Error("changing the struct pointers string value didn't break equality")
		}
		typed.StructPtr.Str = val.StructPtr.Str
		if !reflect.DeepEqual(&val, cp) {
			t.Fatal("Couldn't restored equality")
		}

		typed.StructPtr.IntSlice[2] = 287234
		if reflect.DeepEqual(&val, cp) {
			t.Error("changing the copy's slice didn't break equality")
		}
		typed.StructPtr.IntSlice[2] = val.StructPtr.IntSlice[2]
		if !reflect.DeepEqual(&val, cp) {
			t.Fatal("Couldn't restored equality")
		}
	}

	if val.FloatPtr != nil {
		*typed.FloatPtr = 9856545.4534
		if reflect.DeepEqual(&val, cp) {
			t.Error("changing the float pointer's value didn't break equality")
		}
		*typed.FloatPtr = *val.FloatPtr
		if !reflect.DeepEqual(&val, cp) {
			t.Fatal("Couldn't restored equality")
		}
	}
}
