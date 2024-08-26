package json

import (
	"context"
	"core"
	"core/generic"
	"core/types"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"
)

// Map is standard map, that implements GetterSetter
type Map map[string]interface{}

func (that Map) Contains(name string) bool {
	_, err := that.GetProperty(context.Background(), name)
	return err == nil
}

func (that Map) Hash() (string, error) {
	data, err := MarshalIndent(that)
	if err != nil {
		return "", fmt.Errorf("MarshalSortedIndent: %w", err)
	}

	return digestString(string(data)), nil
}

func (that Map) SaveToFile(filename string) error {
	data, err := MarshalIndent(that)
	if err != nil {
		return fmt.Errorf("MarshalIndent: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("WriteFile: %w", err)
	}

	return nil
}

func (that Map) Marshal() (RawMessage, error) {
	return Marshal(that)
}

func (that Map) MarshalIndent() (RawMessage, error) {
	return MarshalIndent(that)
}

func (that Map) Unmarshal(data []byte) error {
	return Unmarshal(data, &that)
}

func (that Map) Json() (RawMessage, error) {
	res, err := Marshal(that)
	if err != nil {
		return nil, fmt.Errorf("Marshal: %w", err)
	}
	return res, nil
}

func (that Map) String() string {
	return fmt.Sprintf("%#v", that)
}

func (that Map) GetProperty(
	ctx context.Context,
	name string,
) (interface{}, error) {
	if strings.HasPrefix(name, "$.") {
		var path = strings.Split(name[2:], ".")

		if len(path) == 1 {
			if v, ok := that[name]; ok {
				return v, nil
			}
			return nil, fmt.Errorf("key %q not found %w", name, core.ErrNoMatch)
		}

		var val interface{} = that
		for _, key := range path {
			if vv, ok := val.(Map); ok {
				if v, ok := vv[key]; ok {
					val = v
				} else {
					return nil, fmt.Errorf("key %q not found %w", name, core.ErrNoMatch)
				}
			} else {
				return nil, fmt.Errorf("key %q not found %w", name, core.ErrNoMatch)
			}
		}

		return val, nil
	}

	if v, ok := that[name]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("key %q not found %w", name, core.ErrNoMatch)
}

func (that Map) SetProperty(
	ctx context.Context,
	name string,
	value interface{},
) error {
	if strings.HasPrefix(name, "$.") {
		path := strings.Split(name[2:], ".")
		if len(path) == 1 {
			that[name] = value
			return nil
		}

		var val interface{} = that
		for i := 0; i < len(path)-1; i++ {
			key := path[i]
			if vv, ok := val.(Map); ok {
				if v, ok := vv[key]; ok {
					val = v
				} else {
					vv[key] = make(Map)
					val = vv[key]
				}
			} else {
				return fmt.Errorf("key %q not found %w", name, core.ErrNoMatch)
			}
		}

		if vv, ok := val.(Map); ok {
			vv[path[len(path)-1]] = value
			return nil
		}

		return fmt.Errorf("key %q not found %w", name, core.ErrNoMatch)
	}

	that[name] = value
	return nil
}

// ExpandBy is routine, that allow expand map by additional values from another map (recursive).
func (that Map) ExpandBy(right Map) {
	for key, rightVal := range right {
		if leftVal, present := that[key]; present {
			if left, ok := leftVal.(Map); ok {
				if right, ok := rightVal.(Map); ok {
					left.ExpandBy(right)
				}
			}
		} else {
			// key not in left so we can just shove it in
			if right, ok := rightVal.(Map); ok {
				mm := make(Map)
				that[key] = mm
				mm.ExpandBy(right)
			} else {
				that[key] = generic.CloneValue(rightVal)
			}
		}
	}
}

func (that Map) ToBoolean(
	ctx context.Context,
	name string,
	defVal bool,
) bool {
	val, err := that.GetProperty(ctx, name)
	if err != nil || val == nil {
		return defVal
	}
	return types.Type.Boolean.Cast(val, defVal)
}

func (that Map) ToString(
	ctx context.Context,
	name string,
	defVal string,
) string {
	val, err := that.GetProperty(ctx, name)
	if err != nil || val == nil {
		return defVal
	}
	return types.Type.String.Cast(val, defVal)
}

func (that Map) ToInteger(
	ctx context.Context,
	name string,
	defVal int64,
) int64 {
	val, err := that.GetProperty(ctx, name)
	if err != nil || val == nil {
		return defVal
	}
	return types.Type.Integer.Cast(val, defVal)
}

func (that Map) ToFloat(
	ctx context.Context,
	name string,
	defVal float64,
) float64 {
	val, err := that.GetProperty(ctx, name)
	if err != nil || val == nil {
		return defVal
	}
	return types.Type.Float.Cast(val, defVal)
}

func (that Map) ToDuration(
	ctx context.Context,
	name string,
	defVal time.Duration,
) time.Duration {
	val, err := that.GetProperty(ctx, name)
	if err != nil || val == nil {
		return defVal
	}
	return types.Type.Duration.Cast(val, defVal)
}

func (that Map) ToJson(
	ctx context.Context,
	name string,
	defVal RawMessage,
) RawMessage {
	val, err := that.GetProperty(ctx, name)
	if err != nil || val == nil {
		return defVal
	}
	return types.Type.Json.Cast(val, defVal)
}

func (that Map) ToMap(
	ctx context.Context,
	name string,
) Map {
	if mm, ok := that[name]; ok {
		if mmm, ok := mm.(Map); ok {
			return mmm
		}
		if mmm, ok := mm.(map[string]interface{}); ok {
			return mmm
		}
	}
	return nil
}

func (that Map) ToMaps(
	ctx context.Context,
	name string,
) []Map {
	if mm, ok := that[name]; ok {
		if mmm, ok := mm.([]Map); ok {
			return mmm
		}
		if mmm, ok := mm.([]interface{}); ok {
			m1 := make([]Map, 0, len(mmm))
			for _, m2 := range mmm {
				if m3, ok := m2.(map[string]interface{}); ok {
					m1 = append(m1, m3)
				}
			}
			return m1
		}
	}
	return nil
}

func (that Map) ToSlice(
	ctx context.Context,
	name string,
) []interface{} {
	if mm, ok := that[name]; ok {
		switch v := mm.(type) {
		case []interface{}:
			return v
		default:
			val := reflect.ValueOf(v)
			if val.Kind() == reflect.Slice {
				list := make([]interface{}, val.Len())
				for i := 0; i < val.Len(); i++ {
					list[i] = val.Index(i).Interface()
				}
				return list
			}
		}
	}
	return nil
}

func (that Map) GetBoolean(
	ctx context.Context,
	name string,
	defVal bool,
) (res bool, err error) {
	return types.Type.Boolean.Get(ctx, that, name, defVal)
}

func (that Map) GetString(
	ctx context.Context,
	name string,
	defVal string,
) (res string, err error) {
	return types.Type.String.Get(ctx, that, name, defVal)
}

func (that Map) GetInteger(
	ctx context.Context,
	name string,
	defVal int64,
) (res int64, err error) {
	return types.Type.Integer.Get(ctx, that, name, defVal)
}

func (that Map) GetFloat(
	ctx context.Context,
	name string,
	defVal float64,
) (res float64, err error) {
	return types.Type.Float.Get(ctx, that, name, defVal)
}

func (that Map) GetDuration(
	ctx context.Context,
	name string,
	defVal time.Duration,
) (res time.Duration, err error) {
	return types.Type.Duration.Get(ctx, that, name, defVal)
}

func (that Map) GetJson(
	ctx context.Context,
	name string,
	defVal RawMessage,
) (res RawMessage, err error) {
	return types.Type.Json.Get(ctx, that, name, defVal)
}

func (that Map) SetBoolean(
	ctx context.Context,
	name string,
	value bool,
) error {
	return that.SetProperty(ctx, name, value)
}

func (that Map) SetString(
	ctx context.Context,
	name string,
	value string,
) error {
	return that.SetProperty(ctx, name, value)
}

func (that Map) SetInteger(
	ctx context.Context,
	name string,
	value int64,
) error {
	return that.SetProperty(ctx, name, value)
}

func (that Map) SetFloat(
	ctx context.Context,
	name string,
	value float64,
) error {
	return that.SetProperty(ctx, name, value)
}

func (that Map) SetDuration(
	ctx context.Context,
	name string,
	value time.Duration,
) error {
	return that.SetProperty(ctx, name, value)
}

func (that Map) SetJson(
	ctx context.Context,
	name string,
	value RawMessage,
) error {
	return that.SetProperty(ctx, name, value)
}

// Scope is routine, that allow access to the branch of base Map as sub Map.
func (that Map) Scope(name string) Map {
	if mm, ok := that[name]; ok {
		if mmm, ok := mm.(Map); ok {
			return mmm
		}
	}
	return nil
}

func (that Map) NewScope(name string) Map {
	if mm, ok := that[name]; ok {
		if mmm, ok := mm.(Map); ok {
			return mmm
		}
	}
	mmm := make(Map)
	that[name] = mmm
	return mmm
}

func (that Map) Clone() Map {
	cp := make(Map)
	for k, v := range that {
		vm, ok := v.(Map)
		if ok {
			cp[k] = vm.Clone()
		} else {
			cp[k] = v
		}
	}

	return cp
}

// NewMap is constructor for creating Map from JSON source.
func NewMap(source RawMessage) (Map, error) {
	var res map[string]interface{}
	err := Unmarshal(source, &res)
	//	err := JsonUnmarshal([]byte(source), &res)
	if err != nil {
		return nil, fmt.Errorf("create map: %w", err)
	}
	return NewMapFromStruct(res), nil
}

// NewMapFromFiles returns map as result join from files.
func NewMapFromFiles(files ...string) (m Map, err error) {
	m = make(Map)
	for _, file := range files {
		mm, err := NewMapFromFile(file)
		if err != nil {
			return nil, fmt.Errorf("NewMapFromFile: %w", err)
		}
		m.ExpandBy(mm)
	}
	return
}

// NewMapFromFile returns map, that loaded from file
func NewMapFromFile(filename string) (Map, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: %w", err)
	}

	m, err := NewMap(data)
	if err != nil {
		return nil, fmt.Errorf("NewMapFromJson: %w", err)
	}

	return m, nil
}

// NewMaps is constructor for creating Maps from JSON source.
func NewMaps(source RawMessage) ([]Map, error) {
	var res []RawMessage
	err := Unmarshal(source, &res)
	if err != nil {
		return nil, fmt.Errorf("create slice: %w", err)
	}
	result := make([]Map, len(res))
	for i, raw := range res {
		m, err := NewMap(raw)
		if err != nil {
			return nil, fmt.Errorf("NewMapFromJson: %w", err)
		}
		result[i] = m
	}
	return result, nil
}

// NewMapFromStruct is constructor for creating Map from standard map.
func NewMapFromStruct(
	m map[string]interface{},
) Map {
	cp := make(Map)
	for k, v := range m {
		cp[k] = newMapValue(v)
	}
	return cp
}

// NewMapFromNativeStruct is constructor for creating Map from native struct.
func NewMapFromNativeStruct(val interface{}) Map {
	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Struct {
		return nil
	}
	typeOfS := v.Type()

	values := make(Map)
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i)
		key := field.Name
		if isPublic(key) {
			val := newMapValue(v.Field(i).Interface())
			values[key] = val
		}
	}
	return values
}

func newMapValue(
	value interface{},
) interface{} {
	switch v := value.(type) {
	case Map:
		m := make(Map, len(v))
		for k, vv := range v {
			m[k] = newMapValue(vv)
		}
		return m
	case map[string]interface{}:
		m := make(Map, len(v))
		for k, vv := range v {
			m[k] = newMapValue(vv)
		}
		return m
	case []interface{}:
		list := make([]interface{}, len(v))
		for i, v := range v {
			list[i] = newMapValue(v)
		}
		if !isAllMaps(list) {
			return list
		}
		lst := make([]Map, len(list))
		for i, v := range list {
			lst[i] = v.(Map)
		}
		return lst
	case []Map:
		list := make([]interface{}, len(v))
		for i, v := range v {
			list[i] = newMapValue(v)
		}
		return list
	case []map[string]interface{}:
		list := make([]Map, len(v))
		for i, v := range v {
			list[i] = NewMapFromStruct(v)
		}
		return list
	default:
		return v
	}
}

func isPublic(name string) bool {
	return len(name) != 0 && name[0] >= 'A' && name[0] <= 'Z'
}

func isAllMaps(list []interface{}) bool {
	for _, v := range list {
		if _, ok := v.(Map); !ok {
			return false
		}
	}
	return true
}

func digestString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
