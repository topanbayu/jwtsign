package jwt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/martian/parse"
)

func init() {
	parse.Register("jwt.Sign", modifierFromJSON)
}

type Modifier struct {
	secret string
}

type modifierJSON struct {
	Secret string               `json:"secret"`
	Scope  []parse.ModifierType `json:"scope"`
}

func NewModifier(secret string) *Modifier {
	return &Modifier{
		secret: secret,
	}
}

func modifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &modifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	mod := NewModifier(msg.Secret)
	return parse.NewResult(mod, msg.Scope)
}

func (m *Modifier) ModifyResponse(res *http.Response) error {
	nBody := &bytes.Buffer{}
	nBody.Write([]byte("{secret: }" + m.secret))

	newBody := ioutil.NopCloser(nBody)

	newRes := &http.Response{
		StatusCode: res.StatusCode,
		Status:     res.Status,
		Proto:      res.Proto,
		ProtoMajor: res.ProtoMajor,
		ProtoMinor: res.ProtoMinor,
		Header:     res.Header,
		Body:       newBody,
		Request:    res.Request,
	}
	*res = *newRes

	return nil
}
