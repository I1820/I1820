/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 13-11-2017
 * |
 * | File Name:     decoder.go
 * +===============================================
 */

package decoder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Decoder is a way to communicate with user provided decoder
type Decoder struct {
	URL string
}

// New creates new decoder based on given remote address
func New(url string) Decoder {
	return Decoder{
		URL: url,
	}
}

// Decode decodes given data with user provided decoder
func (d Decoder) Decode(payload []byte, id string) (string, error) {
	r, err := http.Post(fmt.Sprintf("%s/api/decode/%s", d.URL, id), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	if r.StatusCode != 200 {
		return "", fmt.Errorf("%s", b)
	}

	return string(b), nil
}
