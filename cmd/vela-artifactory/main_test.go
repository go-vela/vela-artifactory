// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"testing"

	. "github.com/franela/goblin"
)

func TestMain(t *testing.T) {
	g := Goblin(t)

	g.Describe("unmarshalActions", func() {
		g.It("should get actions from JSON string", func() {
			rawJSON := `[{"action":"delete","path":"some/path"},{"action":"upload","sources":["some/sources/*"],"path":"some/other/path","explode":true}]`

			actions, err := unmarshalActions(rawJSON)

			g.Assert(err == nil).IsTrue(fmt.Sprintf("Failed to unmarshal JSON into []Action: %s", err))
			g.Assert(len(actions)).Equal(2)
		})

		g.It("should get arguments from JSON string", func() {
			rawJSON := `[{"action":"delete","path":"some/path"},{"action":"upload","sources":["some/sources/*"],"path":"some/other/path","explode":true}]`

			actions, err := unmarshalActions(rawJSON)

			g.Assert(err == nil).IsTrue(fmt.Sprintf("Failed to unmarshal JSON into []Action: %s", err))
			g.Assert(string(actions[0].RawArguments)).Equal(`{"action":"delete","path":"some/path"}`)
		})
	})
}
