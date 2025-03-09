package jsonObject

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// IMPORTANT: Please don't  try to replace the usage of this type in the DB model fields with something you think is
// 	better (like string, or a type from another package) without actually testing the full set of effects.
// 	IF YOU HAVE AN IDEA TO IMPROVE THE IMPLEMENTATION, PLEASE READ THE COMMENT BELOW FIRST.

// NOTE: You might wonder why this type was created?
// 	The reason was - the "more_data" fields in the database is a JSON/JSONB type. When we get it from the DB,
// 	the driver returns it as a string and if we try to send that value in a response, it goes out as a string
// 	Now, we know what happens when a JSON is sent to clients as a string - all sorts of escape characters
// 	are added to it (the \" and \\\n characters show up). This causes major headache to the clients and is a pain
// 	for the eyes to look at. REST API clients (like Postman) can't parse it easily and tests can fail and so on.
// 	To avoid all those problems, and to send JSON values as they should be sent, this type was created.

// IMPORTANT: THIS TYPE CANNOT HANDLE A TOP LEVEL JSON ARRAY (any valid JSON document starting with '[')
//
//	It can handle arrays nested in a JSON object though.

const topLevelArrayKey = "topLevelArrayKeyb8d8c89aea51f88b1af1144e3e0b8b74ac2a2c257d08cb80ebc99a7262e5dd8c"

type StringAnyMap map[string]any
type Typ struct {
	Valid            bool
	hasTopLevelArray bool
	StringAnyMap
}

var NullJsonObject Typ

// EmptyNotNullJsonObject returns a new blank Typ
// NOTE: We cannot use a var for the EmptyNotNullJsonObject value because when we copy a lot of values around and
//
//	assign the var to multiple values throughout the program in multiple goroutines, we might get the panic message
//	"concurrent map read and map write" indicating that the value is being written and read simultaneously because
//	the same variable is being used at multiple places
func EmptyNotNullJsonObject() Typ {
	return Typ{
		StringAnyMap: StringAnyMap{},
		Valid:        true,
	}
}

func init() {
	NullJsonObject = Typ{
		StringAnyMap: nil,
		Valid:        false,
	}
}

func NewJsonObject(key string, value any) Typ {
	j := Typ{
		StringAnyMap: map[string]any{
			key: value,
		},
		Valid: false,
	}
	return j
}

func (j *Typ) IsEmpty() bool {
	if j.Valid == false || len(j.StringAnyMap) == 0 {
		return true
	}
	return false
}

// ToJsonObject will convert any object type to Typ using json.Marshal and json.Unmarshal
func ToJsonObject(v any) (Typ, error) {
	jsonObj := Typ{
		Valid:        true,
		StringAnyMap: StringAnyMap{},
	}
	var err error

	// If the value is either a byte slice or a string, check if it is already a JSON string or not
	switch v.(type) {
	case string:
		err = jsonObj.Scan(v.(string))
		if err != nil {
			return NullJsonObject, fmt.Errorf("E#1LT8DV - %v", err)
		}
		return jsonObj, nil
	case []byte:
		err = jsonObj.Scan(v.([]byte))
		if err != nil {
			return NullJsonObject, fmt.Errorf("E#1LT8DY - %v", err)
		}
		return jsonObj, nil
	}

	jsonValue, err := json.Marshal(v)
	if err != nil {
		return NullJsonObject, fmt.Errorf("E#1LT8E0 - %v", err)
	}

	err = jsonObj.Scan(jsonValue)
	if err != nil {
		return NullJsonObject, fmt.Errorf("E#1LT8E3 - %v", err)
	}

	return jsonObj, nil
}

func (j *Typ) IsNotEmpty() bool {
	return !j.IsEmpty()
}

func (j *Typ) SetNewTopLevelElement(key string, value any) (replacedExistingKey bool) {
	if j.Valid == false {
		// We are making this object into a valid one
		j.Valid = true
		j.StringAnyMap[key] = value
		return
	}

	replacedExistingKey = false

	if _, ok := j.StringAnyMap[key]; ok {
		// Element already exists
		replacedExistingKey = true
	}

	j.StringAnyMap[key] = value
	return
}

// GetTopLevelElement will return Top-Level element identified by key. If the key does not exist, nil is returned
func (j *Typ) GetTopLevelElement(key string) any {
	if val, ok := j.StringAnyMap[key]; ok {
		return val
	}
	return nil
}

// MARKER: Stringer interface implementation

func (j *Typ) String() string {
	if !j.Valid {
		return ""
	}

	bytes, err := json.Marshal(j.StringAnyMap)
	if err != nil {
		return ""
	}

	return string(bytes)
}

// PrettyString will give the formatted string for this Typ
func (j *Typ) PrettyString() string {
	if !j.Valid {
		return ""
	}

	bytes, err := json.MarshalIndent(j.StringAnyMap, "", "    ")
	if err != nil {
		return ""
	}

	return string(bytes)
}

func (j *Typ) HasTopLevelArray() bool {
	if j.Valid && len(j.StringAnyMap) == 1 && j.hasTopLevelArray {
		if _, ok := j.StringAnyMap[topLevelArrayKey]; ok {
			return true
		}
	}

	return false
}

func (j *Typ) AsByteSlice() []byte {
	return []byte(j.String())
}

// MARKER: DB Interface implementations

// Value implements the driver.Valuer interface. This method returns the JSON-encoded representation of the struct.
func (j *Typ) Value() (driver.Value, error) {
	if j.Valid == false {
		return nil, nil
	}

	if len(j.StringAnyMap) == 0 {
		// Valid empty JSON
		return []byte("{}"), nil
	}

	return j.MarshalJSON()
}

// Scan implements the sql.Scanner interface. This method decodes a JSON-encoded value into the struct fields.
func (j *Typ) Scan(value any) error {
	var arrAnys []any = make([]any, 0)
	switch value.(type) {
	case nil:
		j.Valid = false
		return nil
	case string:
		// Convert to byte slice and try
		err := json.Unmarshal([]byte(value.(string)), &j.StringAnyMap)
		if err != nil {
			err2 := json.Unmarshal(value.([]byte), &arrAnys)
			if err2 != nil {
				return errors.New(fmt.Sprintf("E#1LT8DE - Unmarshalling failed: %v", err))
			}
			j.hasTopLevelArray = true
			j.StringAnyMap = StringAnyMap{
				topLevelArrayKey: arrAnys,
			}
		}
		j.Valid = true
		return nil
	case []byte:
		err := json.Unmarshal(value.([]byte), &j.StringAnyMap)
		if err != nil {
			err2 := json.Unmarshal(value.([]byte), &arrAnys)
			if err2 != nil {
				return errors.New(fmt.Sprintf("E#1LT8DB - Unmarshalling failed: %v", err))
			}
			j.hasTopLevelArray = true
			j.StringAnyMap = StringAnyMap{
				topLevelArrayKey: arrAnys,
			}
		}
		j.Valid = true
		return nil
	default:
		// Attempt to convert
		b, ok := value.([]byte)
		if !ok {
			return errors.New("E#1LT8D7 - Type assertion to []byte failed")
		}

		// return json.Unmarshal(b, &j.StringAnyMap)
		err := json.Unmarshal(b, &j.StringAnyMap)
		if err != nil {
			return errors.New(fmt.Sprintf("E#1LT8D3 - Unmarshalling failed after assertion passed: %v", err))
		}
		j.Valid = true
		return nil
	}
}

// MARKER: Custom implementation of JSON Encoder for this type

// MarshalJSON implements json.Marshaler interface
func (j Typ) MarshalJSON() ([]byte, error) {
	if !j.Valid {
		return []byte("null"), nil
	}
	if j.HasTopLevelArray() {
		return json.Marshal(j.StringAnyMap[topLevelArrayKey])
	}
	return json.Marshal(j.StringAnyMap)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Typ) UnmarshalJSON(dataToUnmarshal []byte) error {
	var err error
	var v any
	if err = json.Unmarshal(dataToUnmarshal, &v); err != nil {
		return err
	}
	switch v.(type) {
	case StringAnyMap:
		err = json.Unmarshal(dataToUnmarshal, &j.StringAnyMap)
	case map[string]any:
		err = json.Unmarshal(dataToUnmarshal, &j.StringAnyMap)
	case []any:
		j.StringAnyMap = StringAnyMap{
			topLevelArrayKey: v.([]any),
		}
		j.hasTopLevelArray = true
	case nil:
		j.Valid = false
		j.StringAnyMap = nil
		return nil
	default:
		err = fmt.Errorf("E#1LT8CZ - Cannot convert object of type %v to Typ", reflect.TypeOf(v).Name())
	}

	j.Valid = true
	if err != nil {
		j.Valid = false
	}

	return err
}
