/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 23-02-2018
 * |
 * | File Name:     encoder/encoder.go
 * +===============================================
 */

package encoder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Encoder is a way to communicate with user provided encoder
type Encoder struct {
	URL string
}

// New creates new encoder based on given remote address
func New(url string) Encoder {
	return Encoder{
		URL: url,
	}
}

// Encode encodes given data with user provided encoder
func (e Encoder) Encode(payload string, id string) ([]byte, error) {
	r, err := http.Post(fmt.Sprintf("%s/api/encode/%s", e.URL, id), "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("%s", b)
	}

	return b, nil
}
