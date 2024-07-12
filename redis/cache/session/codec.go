package session

import (
	"bytes"
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

type MessagePackCodec struct{}

func (MessagePackCodec) Encode(deadline time.Time, values map[string]interface{}) ([]byte, error) {
	s := &struct {
		Deadline int64                  `msgpack:"deadline"`
		Values   map[string]interface{} `msgpack:"values"`
	}{
		Deadline: deadline.UnixNano(),
		Values:   values,
	}

	var b bytes.Buffer
	if err := msgpack.NewEncoder(&b).Encode(s); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (MessagePackCodec) Decode(b []byte) (deadline time.Time, values map[string]interface{}, error error) {
	aux := &struct {
		Deadline int64                  `msgpack:"deadline"`
		Values   map[string]interface{} `msgpack:"values"`
	}{}

	dec := msgpack.NewDecoder(bytes.NewReader(b))
	if err := dec.Decode(aux); err != nil {
		return time.Time{}, nil, err
	}

	return time.Unix(0, aux.Deadline), aux.Values, nil
}
