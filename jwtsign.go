package jwt

import (
	"encoding/json"
	"net/http"

	"github.com/google/martian/log"
	"github.com/google/martian/parse"
)

func init() {
	parse.Register("jwt.Sign", modifierFromJSON)
}

type Modifier struct {
	secret []byte
}

type modifierJSON struct {
	Secret []byte               `json:"secret"`
	Scope  []parse.ModifierType `json:"scope"`
}

// NewModifier constructs and returns a body.Modifier.
func NewModifier(b []byte) *Modifier {
	log.Debugf("len(b): %d, content: %s", len(b), string(b))
	return &Modifier{
		secret: b,
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
	log.Debugf("Modifier: %v", m)

	res.Header.Add("Secret", string(m.secret))
	// res.Body.Close()
	// res.Body = ioutil.NopCloser(bytes.NewReader(m.secret))

	return nil
}
