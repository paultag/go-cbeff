// {{{ Copyright (c) Paul R. Tagliamonte <paultag@gmail.com>, 2019
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE. }}}

// May the FSM have mercy on my poor soul. This package is a shim because
// Go doesn't have any JPEG 2000 parsers, because JPEG 2000 is a dead spec
// that no one uses except for CBEFF, apparently.
//
// As a result, this has been subpackaged, since using this package will
// require cgo linking to imagick. On GNU/Linux this means installing
// the `libmagickwand-dev` package.

// Package jpeg2000 contains a thin wrapper on top of imagick to convert
// JPEG2000 bytes to a golang image.Image.
package jpeg2000

import (
	"bytes"
	"image"
	"image/png"

	"gopkg.in/gographics/imagick.v2/imagick"
)

// Parse a JEPG 2000 into an image.Image.
func Parse(data []byte) (image.Image, error) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.ReadImageBlob(data); err != nil {
		return nil, err
	}

	mw.SetImageFormat("png")
	return png.Decode(bytes.NewReader(mw.GetImageBlob()))
}

// vim: foldmethod=marker
