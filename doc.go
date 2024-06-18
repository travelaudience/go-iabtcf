// Package iabtcf implements a parser of the IAB Transparency and Consent String (TC String) v2.0
//
// This package provides 2 modes:
// - a normal parser that will parse all fields of the TC String and will store them in a Consent object
// - a lazy parser that will just decode the base64 consent string and will store the bytes in a LazyConsent object
//
// The normal parser is useful when you need to access most of the fields and especially if you need to check multiple vendors.
//
// The lazy parser is useful when you need to access only a few fields and only one vendor.
// The parsing is done only when the field is accessed.
// The lazy parser is not optimized for checking multiple vendors.
// Another drawback of the lazy parser is that the client will have to handle the errors when accessing the fields.
package iabtcf
