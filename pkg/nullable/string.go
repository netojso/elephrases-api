package nullable

import (
	"database/sql"
	"encoding/json"
)

type NullableString struct {
	sql.NullString
}

// N creates a new NullableString
func NewNullableString(s string) NullableString {
	return NullableString{
		NullString: sql.NullString{
			String: s,
			Valid:  s != "",
		},
	}
}

// MarshalJSON converts NullableString to JSON
func (ns NullableString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON converts JSON to NullableString
func (ns *NullableString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		ns.String = ""
		return nil
	}
	if err := json.Unmarshal(data, &ns.String); err != nil {
		return err
	}
	ns.Valid = true
	return nil
}
