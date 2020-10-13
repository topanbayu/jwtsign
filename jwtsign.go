// Copyright 2015 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package body allows for the replacement of message body on responses.
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

// Modifier substitutes the body on an HTTP response.
type Modifier struct {
	secret []byte
}

type modifierJSON struct {
	Secret []byte               `json:"secret"`
	Scope  []parse.ModifierType `json:"scope"`
}

// NewModifier constructs and returns a body.Modifier.
func NewModifier(b []byte) *Modifier {
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

// ModifyResponse sets the Content-Type header and overrides the response body.
func (m *Modifier) ModifyResponse(res *http.Response) error {
	// Replace the existing body, close it first.
	res.Body.Close()
	res.Body = ioutil.NopCloser(bytes.NewReader(m.secret))

	return nil
}
