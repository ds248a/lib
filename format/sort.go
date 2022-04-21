package format

import (
	"reflect"

	"sort"
)

const brokenNaNs = false

// SortedMap организует отсортированный хеш Key=>Value для структур reflect.Value.
// Ключи и значения выровнены в порядке индекса: Value[i] соответствует Key[i].
type SortedMap struct {
	Key   []reflect.Value
	Value []reflect.Value
}

func (o *SortedMap) Len() int           { return len(o.Key) }
func (o *SortedMap) Less(i, j int) bool { return compare(o.Key[i], o.Key[j]) < 0 }
func (o *SortedMap) Swap(i, j int) {
	o.Key[i], o.Key[j] = o.Key[j], o.Key[i]
	o.Value[i], o.Value[j] = o.Value[j], o.Value[i]
}

// Sort принимает reflect.Map и возвращмет структуру SortedMap, с сортировкой по ключу.
// Это позволяте устранить проблемы когда ключ имеет значение NaN.
func Sort(mapValue reflect.Value) *SortedMap {
	if mapValue.Type().Kind() != reflect.Map {
		return nil
	}
	key, value := mapElems(mapValue)
	sorted := &SortedMap{
		Key:   key,
		Value: value,
	}
	sort.Stable(sorted)
	return sorted
}

func mapElems(mapValue reflect.Value) ([]reflect.Value, []reflect.Value) {
	n := mapValue.Len()
	key := make([]reflect.Value, 0, n)
	value := make([]reflect.Value, 0, n)
	iter := mapValue.MapRange()

	for iter.Next() {
		key = append(key, iter.Key())
		value = append(value, iter.Value())
	}
	return key, value
}

// compare сравнивает два значения одного типа.
// Возвращает -1, 0 или 1
// К примеру: a > b (1), a == b (0), a < b (-1).
// Если типы отличаются, возвращает (-1).
func compare(aVal, bVal reflect.Value) int {
	aType, bType := aVal.Type(), bVal.Type()
	if aType != bType {
		return -1 // No good answer possible, but don't return 0: they're not equal.
	}
	switch aVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		a, b := aVal.Int(), bVal.Int()
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		a, b := aVal.Uint(), bVal.Uint()
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	case reflect.String:
		a, b := aVal.String(), bVal.String()
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	case reflect.Float32, reflect.Float64:
		return floatCompare(aVal.Float(), bVal.Float())
	case reflect.Complex64, reflect.Complex128:
		a, b := aVal.Complex(), bVal.Complex()
		if c := floatCompare(real(a), real(b)); c != 0 {
			return c
		}
		return floatCompare(imag(a), imag(b))
	case reflect.Bool:
		a, b := aVal.Bool(), bVal.Bool()
		switch {
		case a == b:
			return 0
		case a:
			return 1
		default:
			return -1
		}
	case reflect.Ptr:
		a, b := aVal.Pointer(), bVal.Pointer()
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	case reflect.Chan:
		if c, ok := nilCompare(aVal, bVal); ok {
			return c
		}
		ap, bp := aVal.Pointer(), bVal.Pointer()
		switch {
		case ap < bp:
			return -1
		case ap > bp:
			return 1
		default:
			return 0
		}
	case reflect.Struct:
		for i := 0; i < aVal.NumField(); i++ {
			if c := compare(aVal.Field(i), bVal.Field(i)); c != 0 {
				return c
			}
		}
		return 0
	case reflect.Array:
		for i := 0; i < aVal.Len(); i++ {
			if c := compare(aVal.Index(i), bVal.Index(i)); c != 0 {
				return c
			}
		}
		return 0
	case reflect.Interface:
		if c, ok := nilCompare(aVal, bVal); ok {
			return c
		}
		c := compare(reflect.ValueOf(aVal.Elem().Type()), reflect.ValueOf(bVal.Elem().Type()))
		if c != 0 {
			return c
		}
		return compare(aVal.Elem(), bVal.Elem())
	default:
		// Некоторые типы не могут отображаться в качестве ключей: (maps, funcs, slices)
		panic("bad type in compare: " + aType.String())
	}
}

// nilCompare проверяет, является ли какое-либо значение нулевым (nil).
// Если любое переданное значение равно nil, bool принимает true, а int - соответствующее значение.
// В качестве аргумента могут быть преданы: chan, func, interface, map, pointer, or slice.
func nilCompare(aVal, bVal reflect.Value) (int, bool) {
	if aVal.IsNil() {
		if bVal.IsNil() {
			return 0, true
		}
		return -1, true
	}
	if bVal.IsNil() {
		return 1, true
	}
	return 0, false
}

// floatCompare сопоставляет значения типа float.
func floatCompare(a, b float64) int {
	switch {
	case isNaN(a):
		return -1
	case isNaN(b):
		return 1
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func isNaN(a float64) bool {
	return a != a
}
