package util

import (
	"fmt"
	"reflect"
)

// Recursive function that does all of the heavy lifting for the Copy function
func copyRecursive(in reflect.Value) (out reflect.Value) {
	defer func() {
		// since we started using this in update_record we always need out to
		// be settable, so we make a new value that should always be settable
		if !out.CanSet() {
			realOut := reflect.New(out.Type()).Elem()
			realOut.Set(out)
			out = realOut
		}
	}()

	// check for all nillable values first so we don't have to run a separate
	// check in all of the relevant kind handlers
	switch in.Kind() {
	case reflect.Chan, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if in.IsNil() {
			return in
		}
	}

	switch in.Kind() {
	case reflect.Ptr:
		inElem := in.Elem()
		out = reflect.New(inElem.Type())
		out.Elem().Set(copyRecursive(inElem))

	case reflect.Struct:
		out = reflect.New(in.Type()).Elem()
		for i, num := 0, in.Type().NumField(); i < num; i += 1 {
			if dest := out.Field(i); dest.CanSet() {
				dest.Set(copyRecursive(in.Field(i)))
			}
		}

	case reflect.Slice:
		out = reflect.MakeSlice(in.Type(), 0, in.Len())
		for i, max := 0, in.Len(); i < max; i += 1 {
			out = reflect.Append(out, copyRecursive(in.Index(i)))
		}

	case reflect.Map:
		out = reflect.MakeMap(in.Type())
		for _, key := range in.MapKeys() {
			out.SetMapIndex(key, copyRecursive(in.MapIndex(key)))
		}

	case reflect.Chan:
		out = reflect.MakeChan(in.Type().Elem(), in.Cap())
	case reflect.Interface:
		out = copyRecursive(in.Elem()).Convert(in.Type())
	default:
		// most things aren't stored by reference, so if we haven't explicitly
		// handled the kind it's probably safe to just use the raw value.
		out = in
	}

	return out
}

// Creates a deep copy using reflection
func Copy(orig interface{}) (cp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panicked copying struct: %v -\n\torig:%v\n\tcopy:%v\n", r, orig, cp)
			cp = nil
			err = fmt.Errorf("caught panic: %v", r)
		}
	}()

	return copyRecursive(reflect.ValueOf(orig)).Interface(), nil
}
