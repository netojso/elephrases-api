package pkg

import (
	"github.com/google/uuid"
)

type UUID struct {
	value uuid.UUID
}

func (u UUID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.value.String() + `"`), nil
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	parsedUUID, err := uuid.Parse(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	u.value = parsedUUID
	return nil
}

func NewUUID() UUID {
	return UUID{value: uuid.New()}
}

func (u UUID) String() string {
	return u.value.String()
}

func (u UUID) Value() uuid.UUID {
	return u.value
}

func ParseUUID(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, err
	}
	return UUID{value: id}, nil
}
