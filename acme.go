// Copyright (c) Alisdair MacLeod <copying@alisdairmacleod.co.uk>
//
// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
// REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
// INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
// LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
// OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
// PERFORMANCE OF THIS SOFTWARE.

// Package acme provides useful extensions to the 9fans.net/go/acme library for writing acme utilities.
package acme

import (
	"fmt"

	"9fans.net/go/acme"
)

// ReplaceSelection replaces the currently selected text with the result of replaceFunc and reselects the text.
func ReplaceSelection(winid int, replaceFunc func(selection string) (string, error)) error {
	win, err := acme.Open(winid, nil)
	if err != nil {
		return fmt.Errorf("opening window: %w", err)
	}
	// acme zeroes addr the first time it is read.
	if _, _, addrErr := win.ReadAddr(); addrErr != nil {
		return fmt.Errorf("reading addr: %w", addrErr)
	}
	selection := win.Selection()
	if selection == "" {
		return fmt.Errorf("no selection")
	}
	replacement, err := replaceFunc(selection)
	if err != nil {
		return fmt.Errorf("running replacement function: %w", err)
	}
	q0, _, err := win.ReadAddr()
	if err != nil {
		return fmt.Errorf("reading addr: %w", err)
	}
	// after reading the selection q0 will be set to the end of the selection.
	startAddr := q0 - len(selection)
	endAddr := startAddr + len(replacement)
	if err := win.Ctl("addr=dot"); err != nil {
		return fmt.Errorf("setting addr=dot: %w", err)
	}
	if _, err := win.Write("data", []byte(replacement)); err != nil {
		return fmt.Errorf("writing data: %w", err)
	}
	if err := win.Addr("#%d,#%d", startAddr, endAddr); err != nil {
		return fmt.Errorf("setting addr: %w", err)
	}
	if err := win.Ctl("dot=addr"); err != nil {
		return fmt.Errorf("setting dot=addr: %w", err)
	}
	return nil
}
