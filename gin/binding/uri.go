// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

type uriBinding struct{}

// Name ...
func (uriBinding) Name() string {
	return "uri"
}

// BindUri ...
func (uriBinding) BindUri(m map[string][]string, obj any) error {
	if err := mapURI(obj, m); err != nil {
		return err
	}
	return validate(obj)
}
