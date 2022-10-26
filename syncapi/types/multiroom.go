package types

type MultiRoom map[string]map[string]MultiRoomData

type MultiRoomData []byte

func (d MultiRoomData) MarshalJSON() ([]byte, error) {
	return d, nil
}

type MultiRoomDataRow struct {
	Data   []byte
	Type   string
	UserId string
}
