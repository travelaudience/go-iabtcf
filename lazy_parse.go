package iabtcf

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// LazyParseCoreString parses a TCF consent string into a LazyConsent
//
// note: The parser is only decoding the base64 string,
// checking if consent string has at least the minimum length to decode most of the fields,
// and then storing the bytes.
// The field parsing is done only when the field is accessed.
// Since the minimum length check is done, all fields ( except vendor fields ) can be accessed without error.
// For the vendor part, if the consent string is too short or invalid, the vendor will be considered as not allowed.
//
// note: the lazy parser is optimized for checking only one vendor + few fields
func LazyParseCoreString(c string) (*LazyConsent, error) {
	if c == "" {
		return nil, fmt.Errorf("consent string is empty")
	}
	everything := strings.Split(c, ".")

	// extract core string
	var bytes, err = base64.RawURLEncoding.DecodeString(everything[0])
	if err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	// Extract disclosed vendors and publisher TC blocks.  There are an arbitrary number of these blocks in any order,
	// and each block needs to be inspected to see what it is.
	consent := NewLazyConsent(bytes)
	for _, extra := range everything[1:] {
		if bytes, err := base64.RawURLEncoding.DecodeString(extra); err == nil {
			consent.Extras = append(consent.Extras, Bits(bytes))
		}
	}

	// note: is_range_encoding is the last fixed field.
	// after this bit, we will have either range entries ( up to num_entries ) or vendor bitset ( up to max_vendor_id ).
	// we are just checking here that we are able to read at minimum the fixed fields.
	// if after this bit, the consent string is too short or invalid, we will just return that the vendor is not allowed
	if consent.Core.Length() < IsRangeEncodingField.NextOffset() {
		return nil, fmt.Errorf("consent string is too short")
	}

	return consent, nil
}

// //////////////////////////////////////////////////
// lazy consent

// LazyConsent provides methods to extract data in a lazy mode ( to avoid parsing all fields for nothing )
type LazyConsent struct {
	Core   Bits
	Extras []Bits
}

func NewLazyConsent(bytes []byte) *LazyConsent {
	return &LazyConsent{
		Core: Bits(bytes),
	}
}

// Version returns the version of the consent string
func (c *LazyConsent) Version() int {
	return c.Core.ReadIntField(VersionField.Offset, VersionField.NbBits)
}

// Created returns the creation date of the consent string
func (c *LazyConsent) Created() time.Time {
	return c.Core.ReadTimeField(CreatedField.Offset)
}

// LastUpdated returns the last update date of the consent string
func (c *LazyConsent) LastUpdated() time.Time {
	return c.Core.ReadTimeField(LastUpdatedField.Offset)
}

// CMPID returns the Consent Management Platform ID
func (c *LazyConsent) CMPID() int {
	return c.Core.ReadIntField(CMPIDField.Offset, CMPIDField.NbBits)
}

// CMPVersion returns the Consent Management Platform version
func (c *LazyConsent) CMPVersion() int {
	return c.Core.ReadIntField(CMPVersionField.Offset, CMPVersionField.NbBits)
}

// ConsentScreen returns the consent screen number
func (c *LazyConsent) ConsentScreen() int {
	return c.Core.ReadIntField(ConsentScreenField.Offset, ConsentScreenField.NbBits)
}

// ConsentLanguage returns the consent language
func (c *LazyConsent) ConsentLanguage() string {
	return c.Core.ReadStringField(ConsentLanguageField.Offset, ConsentLanguageField.NbBits)
}

// VendorListVersion returns the vendor list version
func (c *LazyConsent) VendorListVersion() int {
	return c.Core.ReadIntField(VendorListVersionField.Offset, VendorListVersionField.NbBits)
}

// TcfPolicyVersion returns the TCF policy version
func (c *LazyConsent) TcfPolicyVersion() int {
	return c.Core.ReadIntField(TcfPolicyVersionField.Offset, TcfPolicyVersionField.NbBits)
}

// IsServiceSpecific checks if the consent is service specific
func (c *LazyConsent) IsServiceSpecific() bool {
	return c.Core.ReadBoolField(IsServiceSpecificField.Offset)
}

// UseNonStandardStacks checks if the consent uses non standard stacks
func (c *LazyConsent) UseNonStandardStacks() bool {
	return c.Core.ReadBoolField(UseNonStandardStacksField.Offset)
}

// PurposeOneTreatment checks if the consent is for one treatment
func (c *LazyConsent) PurposeOneTreatment() bool {
	return c.Core.ReadBoolField(PurposeOneTreatmentField.Offset)
}

// PublisherCC returns the publisher country code
func (c *LazyConsent) PublisherCC() string {
	return c.Core.ReadStringField(PublisherCCField.Offset, PublisherCCField.NbBits)
}

// IsRangeEncoding checks if the consent is using range encoding
func (c *LazyConsent) IsRangeEncoding() bool {
	return c.Core.ReadBoolField(IsRangeEncodingField.Offset)
}

// NumRangeEntries returns the number of range entries
func (c *LazyConsent) NumRangeEntries() int {
	return c.Core.ReadIntField(NumRangeEntriesField.Offset, NumRangeEntriesField.NbBits)
}

// MaxVendorID returns the maximum vendor ID
func (c *LazyConsent) MaxVendorID() int {
	return c.Core.ReadIntField(MaxVendorIDField.Offset, MaxVendorIDField.NbBits)
}

// EveryPurposeAllowed checks if every purpose number is allowed
func (c *LazyConsent) EveryPurposeAllowed(numbers []int) bool {
	for _, number := range numbers {
		if allowed := c.PurposeAllowed(number); !allowed {
			return false
		}
	}
	return true
}

// PurposeAllowed checks if purpose is allowed
func (c *LazyConsent) PurposeAllowed(number int) bool {
	return c.Core.ReadBitNumber(number, PurposesConsentField.Offset, PurposesConsentField.NbBits)
}

// EveryPurposeLITransparencyAllowed checks if every purposeLITransparency number is allowed
func (c *LazyConsent) EveryPurposeLITransparencyAllowed(numbers []int) bool {
	for _, number := range numbers {
		if allowed := c.PurposeLITransparencyAllowed(number); !allowed {
			return false
		}
	}
	return true
}

// PurposeLITransparencyAllowed checks if purposeLITransparency is allowed
func (c *LazyConsent) PurposeLITransparencyAllowed(number int) bool {
	return c.Core.ReadBitNumber(number, PurposesLITransparencyField.Offset, PurposesLITransparencyField.NbBits)
}

// EverySpecialFeatureAllowed checks every special feature number is allowed
func (c *LazyConsent) EverySpecialFeatureAllowed(numbers []int) bool {
	for _, number := range numbers {
		if allowed := c.SpecialFeatureAllowed(number); !allowed {
			return false
		}
	}
	return true
}

// SpecialFeatureAllowed checks if special feature is allowed
func (c *LazyConsent) SpecialFeatureAllowed(number int) bool {
	return c.Core.ReadBitNumber(number, SpecialFeatureOptInsField.Offset, SpecialFeatureOptInsField.NbBits)
}

// VendorAllowed checks if vendor is in the list of vendors user has given his consent to
func (c *LazyConsent) VendorAllowed(number int) bool {

	if c.IsRangeEncoding() {

		numEntries := c.NumRangeEntries()
		if numEntries == 0 {
			return false
		}

		offset := NumRangeEntriesField.NextOffset()
		for i := 0; i < int(numEntries); i++ {

			isRange := c.Core.ReadBoolField(offset)
			offset += 1

			start := c.Core.ReadIntField(offset, 16)
			offset += 16

			end := start
			if isRange {
				end = c.Core.ReadIntField(offset, 16)
				offset += 16
			}

			if start <= number && number <= end {
				return true
			}
		}

		return false
	}

	maxVendorId := c.MaxVendorID()
	if maxVendorId == 0 {
		return false
	}

	return c.Core.ReadBitNumber(number, ConsentedVendorsOffset, maxVendorId)
}

// //////////////////////////////////////////////////
// consent field helpers

var (
	VersionField                = NewConsentIntField(6)
	CreatedField                = NewConsentTimeField()
	LastUpdatedField            = NewConsentTimeField()
	CMPIDField                  = NewConsentIntField(12)
	CMPVersionField             = NewConsentIntField(12)
	ConsentScreenField          = NewConsentIntField(6)
	ConsentLanguageField        = NewConsentStringField(12)
	VendorListVersionField      = NewConsentIntField(12)
	TcfPolicyVersionField       = NewConsentIntField(6)
	IsServiceSpecificField      = NewConsentBoolField()
	UseNonStandardStacksField   = NewConsentBoolField()
	SpecialFeatureOptInsField   = NewConsentBitsField(12)
	PurposesConsentField        = NewConsentBitsField(24)
	PurposesLITransparencyField = NewConsentBitsField(24)
	PurposeOneTreatmentField    = NewConsentBoolField()
	PublisherCCField            = NewConsentStringField(12)
	MaxVendorIDField            = NewConsentIntField(16)
	IsRangeEncodingField        = NewConsentBoolField()

	// if range encoding, number of range entries, then each range entries
	NumRangeEntriesField = NewConsentIntField(12)

	// if not range encoding, one bit for each vendor up to the max vendor id
	ConsentedVendorsOffset = IsRangeEncodingField.NextOffset()
)

type ConsentField struct {
	Offset int
	NbBits int
}

var previousField *ConsentField

func NewConsentIntField(nbBits int) *ConsentField {
	return NewConsentField(nbBits)
}

func NewConsentTimeField() *ConsentField {
	return NewConsentField(timeNbBits)
}

func NewConsentBoolField() *ConsentField {
	return NewConsentField(boolNbBits)
}

func NewConsentBitsField(nbBits int) *ConsentField {
	return NewConsentField(nbBits)
}

func NewConsentStringField(nbBits int) *ConsentField {
	return NewConsentField(nbBits)
}

func NewConsentField(nbBits int) *ConsentField {
	field := &ConsentField{Offset: previousField.NextOffset(), NbBits: nbBits}
	previousField = field
	return field
}

func NewConsentFieldFromOffset(offset, nbBits int) *ConsentField {
	field := &ConsentField{Offset: offset, NbBits: nbBits}
	previousField = field
	return field
}

func (f *ConsentField) NextOffset() int {
	if f == nil {
		return 0
	}
	return f.Offset + f.NbBits
}
