package json

import (
	"encoding/json"

	"github.com/Kratos40-sba/urlshort/shorturl"
	"github.com/pkg/errors"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*shorturl.Redirect, error) {
	redirect := &shorturl.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}
	return redirect, nil

}
func (r *Redirect) Encode(input *shorturl.Redirect) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}
	return rawMsg, nil
}
