/*
Package populate implements convertation from interface{} to map
*/
package populate

import (
	"fmt"
	"math"
	"strconv"

	"github.com/alrusov/misc"
)

//----------------------------------------------------------------------------------------------------------------------------//

// Fields --
type Fields struct {
	list misc.StringMap
}

//----------------------------------------------------------------------------------------------------------------------------//

func populate(name []byte, data interface{}, fields *Fields) {
	switch fVal := data.(type) {
	case map[string]interface{}:
		if len(name) > 0 {
			name = append(name, '.')
		}
		for fName, field := range fVal {
			populate(append(name, []byte(fName)...), field, fields)
		}

	case []interface{}:
		for i, field := range fVal {
			fName :=
				append(
					append(
						append(name, '['),
						[]byte(strconv.Itoa(i))...),
					']')
			populate(fName, field, fields)
		}

	default:
		v := ""
		switch data.(type) {
		case bool:
			if data.(bool) {
				v = "true"
			} else {
				v = "false"
			}
		case string:
			v = data.(string)
		case int64:
			v = fmt.Sprintf("%d", data.(int64))
		case float64:
			v = misc.TrimStringAsFloat(fmt.Sprintf("%f", data.(float64)))
		}
		fields.list[string(name)] = v
	}

}

// Do --
func Do(data interface{}) *Fields {
	fields := &Fields{
		list: make(misc.StringMap),
	}
	populate(make([]byte, 0), data, fields)
	return fields
}

//----------------------------------------------------------------------------------------------------------------------------//

// String --
func (f Fields) String(name string) (val string, exists bool, err error) {
	val, exists = f.list[name]
	if !exists {
		return "", false, nil
	}

	return val, true, nil
}

//----------------------------------------------------------------------------------------------------------------------------//

// Bool --
func (f Fields) Bool(name string) (val bool, exists bool, err error) {
	v, exists := f.list[name]
	if !exists {
		return false, false, nil
	}

	switch v {
	case "true":
		return true, true, nil
	case "false":
		return false, true, nil
	}

	return false, true, fmt.Errorf(`Bad boolean "%s"`, v)
}

//----------------------------------------------------------------------------------------------------------------------------//

// Int64 --
func (f Fields) Int64(name string) (val int64, exists bool, err error) {
	v, exists := f.list[name]
	if !exists {
		return 0, false, nil
	}

	val, err = strconv.ParseInt(v, 10, 64)

	return val, true, err
}

//----------------------------------------------------------------------------------------------------------------------------//

// Int32 --
func (f Fields) Int32(name string) (val int32, exists bool, err error) {
	v, exists, err := f.Int64(name)
	if !exists || err != nil {
		return 0, exists, err
	}

	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0, true, fmt.Errorf("int32 overflow (%d)", v)
	}

	val = int32(v)
	return val, true, err
}

//----------------------------------------------------------------------------------------------------------------------------//

// Float64 --
func (f Fields) Float64(name string) (val float64, exists bool, err error) {
	v, exists := f.list[name]
	if !exists {
		return 0., false, nil
	}

	val, err = strconv.ParseFloat(v, 64)

	return val, true, err
}

//----------------------------------------------------------------------------------------------------------------------------//
