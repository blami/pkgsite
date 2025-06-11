// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fetch

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/pkgsite/internal/testing/testhelper"
)

func TestMatchingFiles(t *testing.T) {
	plainGoBody := `
		package plain
		type Value int`
	jsGoBody := `
		// +build js,wasm

		// Package js only works with wasm.
		package js
		type Value int`
	synctestBody := `
		//go:build goexperiment.synctest

		package synctest
		var X int
	`
	jsonv2Body := `
		//go:build goexperiment.jsonv2

		package json
		var X int
	`

	plainContents := map[string]string{
		"README.md":  "THIS IS A README",
		"LICENSE.md": testhelper.MITLicense,
		"plain.go":   plainGoBody,
	}

	jsContents := map[string]string{
		"README.md":  "THIS IS A README",
		"LICENSE.md": testhelper.MITLicense,
		"js.go":      jsGoBody,
	}

	synctestContents := map[string]string{
		"plain.go": plainGoBody,
		"st.go":    synctestBody,
	}

	jsonv2Contents := map[string]string{
		"plain.go": plainGoBody,
		"json.go":  jsonv2Body,
	}

	for _, test := range []struct {
		name         string
		goos, goarch string
		importPath   string
		contents     map[string]string
		want         map[string][]byte
	}{
		{
			name:     "plain-linux",
			goos:     "linux",
			goarch:   "amd64",
			contents: plainContents,
			want: map[string][]byte{
				"plain.go": []byte(plainGoBody),
			},
		},
		{
			name:     "plain-js",
			goos:     "js",
			goarch:   "wasm",
			contents: plainContents,
			want: map[string][]byte{
				"plain.go": []byte(plainGoBody),
			},
		},
		{
			name:     "wasm-linux",
			goos:     "linux",
			goarch:   "amd64",
			contents: jsContents,
			want:     map[string][]byte{},
		},
		{
			name:     "wasm-js",
			goos:     "js",
			goarch:   "wasm",
			contents: jsContents,
			want: map[string][]byte{
				"js.go": []byte(jsGoBody),
			},
		},
		{
			name:     "synctest",
			goos:     "linux",
			goarch:   "amd64",
			contents: synctestContents,
			want: map[string][]byte{
				"plain.go": []byte(plainGoBody),
				"st.go":    []byte(synctestBody),
			},
		},
		{
			name:       "json",
			goos:       "linux",
			goarch:     "amd64",
			importPath: "encoding/json",
			contents:   jsonv2Contents,
			want: map[string][]byte{
				"plain.go": []byte(plainGoBody),
			},
		},
		{
			name:       "jsonv2",
			goos:       "linux",
			goarch:     "amd64",
			importPath: "encoding/json/v2",
			contents:   jsonv2Contents,
			want: map[string][]byte{
				"plain.go": []byte(plainGoBody),
				"json.go":  []byte(jsonv2Body),
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			files := map[string][]byte{}
			for n, c := range test.contents {
				files[n] = []byte(c)
			}
			got, err := matchingFiles(test.goos, test.goarch, test.importPath, files)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
