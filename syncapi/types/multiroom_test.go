package types

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
)

func TestMarshallMultiRoom(t *testing.T) {
	is := is.New(t)
	m, err := json.Marshal(MultiRoom{"@3:example.com": map[string]MultiRoomData{"location": MultiRoomData(`{"foo":"bar"}`)}})
	is.NoErr(err)
	is.Equal(m, []byte(`{"@3:example.com":{"location":{"foo":"bar"}}}`))
}
