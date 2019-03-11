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

// CBEFF (Common Biometric Exchange Formats Framework) is a set of ISO
// standards defining an approach to facilitate serialisation and sharing of
// biometric data in an implementation agnostic manner. This is achieved
// through use of a data structure which both describes, and contains,
// biometric data.
//
// This format is most notibly used as part of the US Government's FIPS 201
// PIV II smartcard.
//
// Currently, only Facial support has been partially implemented, due to
// the paywalled documetation on large swaths of this encoding.
// This package has been implemented only based on freely avalible US Government
// provided examples of PIV II data and documentation. As such this may be
// incomplete or slightly wrong in some aspects. Please submit pull requests
// as issues are triaged.
package cbeff // import "pault.ag/go/cbeff"

// vim: foldmethod=marker
