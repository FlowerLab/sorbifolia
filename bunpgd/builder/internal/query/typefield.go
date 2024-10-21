package query

import (
	"cmp"
	"reflect"
	"slices"
	"strings"
	"sync"
)

type StructFields struct {
	List      []Field
	NameIndex map[string]int
}

var fieldCache sync.Map // map[reflect.Type]StructFields

// CachedTypeFields is like TypeFields but uses a cache to avoid repeated work.
func CachedTypeFields(t reflect.Type) StructFields {
	val, ok := fieldCache.Load(t)
	if !ok {
		val, _ = fieldCache.LoadOrStore(t, TypeFields(t))
	}
	return val.(StructFields)
}

// TypeFields returns a List of fields that JSON should recognize for the given type.
// The algorithm is breadth-first search over the set of structs to include - the top struct
// and then any reachable anonymous structs.
func TypeFields(t reflect.Type) StructFields {
	var (
		current          []Field // Anonymous fields to explore at the current level and the next.
		next             = []Field{{Typ: t}}
		count, nextCount map[reflect.Type]int      // Count of queued names for current level and the next.
		visited          = map[reflect.Type]bool{} // Types already visited at an earlier level.
		fields           []Field                   // Fields found.
	)

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.Typ] {
				continue
			}
			visited[f.Typ] = true

			// Scan f.Typ for fields to include.
			for i := 0; i < f.Typ.NumField(); i++ {
				sf := f.Typ.Field(i)
				if sf.Anonymous {
					if t = sf.Type; t.Kind() == reflect.Pointer {
						t = t.Elem()
					}
					if !sf.IsExported() && t.Kind() != reflect.Struct {
						// Ignore embedded fields of unexported non-struct types.
						continue
					}
					// Do not ignore embedded fields of unexported struct types
					// since they may have exported fields.
				} else if !sf.IsExported() {
					// Ignore unexported non-embedded fields.
					continue
				}

				field := ParseField(sf)
				if field.Name == "" {
					continue
				}

				field.Index = make([]int, len(f.Index)+1)
				copy(field.Index, f.Index)
				field.Index[len(f.Index)] = i

				if field.Flag.And(reflect.Pointer) {
					field.Typ = field.Typ.Elem()
				}

				// Record found Field and Index sequence.
				if field.Name != "" || !sf.Anonymous || field.Typ.Kind() != reflect.Struct {
					fields = append(fields, field)
					if count[f.Typ] > 1 {
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 or 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				// Record new anonymous struct to explore in next round.
				nextCount[field.Typ]++
				if nextCount[field.Typ] == 1 {
					next = append(next, field)
				}
			}
		}
	}

	slices.SortFunc(fields, func(a, b Field) int {
		if c := strings.Compare(a.Name, b.Name); c != 0 {
			return c
		}
		if c := cmp.Compare(len(a.Index), len(b.Index)); c != 0 {
			return c
		}
		return slices.Compare(a.Index, b.Index)
	})

	// Delete all fields that are hidden by the Go rules for embedded fields,
	// except that fields with JSON tags are promoted.

	// The fields are sorted in primary order of Name, secondary order
	// of Field Index length. Loop over names; for each Name, delete
	// hidden fields by choosing the one dominant Field that survives.
	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance {
		// One iteration per Name.
		// Find the sequence of fields with the Name of this first Field.
		fi := fields[i]
		name := fi.Name
		for advance = 1; i+advance < len(fields); advance++ {
			fj := fields[i+advance]
			if fj.Name != name {
				break
			}
		}
		if advance == 1 { // Only one Field with this Name
			out = append(out, fi)
			continue
		}
		dominant, ok := dominantField(fields[i : i+advance])
		if ok {
			out = append(out, dominant)
		}
	}

	fields = out
	slices.SortFunc(fields, func(i, j Field) int {
		return slices.Compare(i.Index, j.Index)
	})

	nameIndex := make(map[string]int, len(fields))
	for i, v := range fields {
		nameIndex[v.Name] = i
	}
	return StructFields{List: fields, NameIndex: nameIndex}
}

func dominantField(fields []Field) (Field, bool) {
	if len(fields) > 1 && len(fields[0].Index) == len(fields[1].Index) {
		return Field{}, false
	}
	return fields[0], true
}
