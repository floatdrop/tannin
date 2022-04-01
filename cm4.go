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

// cm4 is a small conservative-update count-min sketch implementation with 4-bit counters
type cm4 struct {
	s    [depth]nvec
	mask uint32
}

const depth = 4

func newCM4(w int) *cm4 {
	if w < 1 {
		panic("cm4: bad width")
	}

	w32 := nextPowerOfTwo(uint32(w))
	c := cm4{
		mask: w32 - 1,
	}

	for i := 0; i < depth; i++ {
		c.s[i] = newNvec(int(w32))
	}

	return &c
}

func (c *cm4) add(keyh uintptr) {
	h1, h2 := uint32(keyh), uint32(keyh>>32)

	for i := range c.s {
		pos := (h1 + uint32(i)*h2) & c.mask
		c.s[i].inc(pos)
	}
}

func (c *cm4) estimate(keyh uintptr) byte {
	h1, h2 := uint32(keyh), uint32(keyh>>32)

	var min byte = 255
	for i := 0; i < depth; i++ {
		pos := (h1 + uint32(i)*h2) & c.mask
		v := c.s[i].get(pos)
		if v < min {
			min = v
		}
	}
	return min
}

func (c *cm4) reset() {
	for _, n := range c.s {
		n.reset()
	}
}

// nybble vector
type nvec []byte

func newNvec(w int) nvec {
	return make(nvec, w/2)
}

func (n nvec) get(i uint32) byte {
	// Ugly, but as a single expression so the compiler will inline it :/
	return byte(n[i/2]>>((i&1)*4)) & 0x0f
}

func (n nvec) inc(i uint32) {
	idx := i / 2
	shift := (i & 1) * 4
	v := (n[idx] >> shift) & 0x0f
	if v < 15 {
		n[idx] += 1 << shift
	}
}

func (n nvec) reset() {
	for i := range n {
		n[i] = (n[i] >> 1) & 0x77
	}
}

// return the integer >= i which is a power of two
func nextPowerOfTwo(i uint32) uint32 {
	n := i - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}
