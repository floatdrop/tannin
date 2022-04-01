// MIT License

// Copyright (c) 2021 Damian Gryski

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package tannin

import (
	"testing"
)

func TestNvec(t *testing.T) {

	n := newNvec(8)

	n.inc(0)
	if n[0] != 0x01 {
		t.Errorf("n[0]=0x%02x, want 0x01: (n=% 02x)", n[0], n)
	}
	if w := n.get(0); w != 1 {
		t.Errorf("n.get(0)=%d, want 1", w)
	}
	if w := n.get(1); w != 0 {
		t.Errorf("n.get(1)=%d, want 0", w)
	}

	n.inc(1)
	if n[0] != 0x11 {
		t.Errorf("n[0]=0x%02x, want 0x11: (n=% 02x)", n[0], n)
	}
	if w := n.get(0); w != 1 {
		t.Errorf("n.get(0)=%d, want 1", w)
	}
	if w := n.get(1); w != 1 {
		t.Errorf("n.get(1)=%d, want 1", w)
	}

	for i := 0; i < 14; i++ {
		n.inc(1)
	}
	if n[0] != 0xf1 {
		t.Errorf("n[1]=0x%02x, want 0xf1: (n=% 02x)", n[0], n)
	}
	if w := n.get(1); w != 15 {
		t.Errorf("n.get(1)=%d, want 15", w)
	}
	if w := n.get(0); w != 1 {
		t.Errorf("n.get(0)=%d, want 1", w)
	}

	// ensure clamped
	for i := 0; i < 3; i++ {
		n.inc(1)
		if n[0] != 0xf1 {
			t.Errorf("n[0]=0x%02x, want 0xf1: (n=% 02x)", n[0], n)
		}
	}

	n.reset()

	if n[0] != 0x70 {
		t.Errorf("n[0]=0x%02x, want 0x70 (n=% 02x)", n[0], n)
	}
}

func TestCM4(t *testing.T) {

	cm := newCM4(32)

	hash := uintptr(0x0ddc0ffeebadf00d)

	cm.add(hash)
	cm.add(hash)

	if got := cm.estimate(hash); got != 2 {
		t.Errorf("cm.estimate(%x)=%d, want 2\n", hash, got)
	}
}