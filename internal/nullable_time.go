package internal

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullableTime struct {
	sql.NullTime
}

func NewNullableTime(t time.Time) NullableTime {
	return NullableTime{
		NullTime: sql.NullTime{
			Time:  t,
			Valid: !t.IsZero(),
		},
	}
}

// MarshalJSON converts NullableTime to JSON
func (ns NullableTime) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Time)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON converts JSON to NullableTime
func (ns *NullableTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		ns.Time = time.Time{}
		return nil
	}
	if err := json.Unmarshal(data, &ns.Time); err != nil {
		return err
	}
	ns.Valid = true
	return nil
}
