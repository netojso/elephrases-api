package internal

import (
	"database/sql"
	"encoding/json"
	"regexp"
)

type NullableString struct {
	sql.NullString
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

// IsValidEmail checks if the NullableString contains a valid email address
func (ns NullableString) IsValidUuid() bool {
	if !ns.Valid {
		return false
	}
	return regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`).MatchString(ns.String)
}
