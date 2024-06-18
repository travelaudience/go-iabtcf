package iabtcf

import (
	"encoding/base64"
	"strings"
	"time"
)

// LazyParseCoreString parses a TCF consent string into a LazyConsent
//
// note: the parser is only decoding the base64 string and storing the bytes.
// the parsing is done only when the field is accessed.
// the client will have to handle the errors when accessing the fields
//
// note: the lazy parser is optimized for checking only one vendor + few fields
func LazyParseCoreString(c string) (*LazyConsent, error) {
	if c == "" {
		return nil, ErrEmptyString
	}
	// extract core string
	cs, _, _ := strings.Cut(c, ".")

	var bytes, err = base64.RawURLEncoding.DecodeString(cs)
	if err != nil {
		return nil, ErrDecodeFailed(err)
	}

	return NewLazyConsent(bytes), nil
}

// //////////////////////////////////////////////////
// lazy consent

// LazyConsent provides methods to extract data in a lazy mode ( to avoid parsing all fields for nothing )
type LazyConsent struct {
	Bits
}

func NewLazyConsent(bytes []byte) *LazyConsent {
	// note: at this point, we do not know the real size of the bitmap
	// assumption: the length should be at least lower than 8 times the number of bytes
	length := 8 * len(bytes)
	return &LazyConsent{
		Bits: Bits{
			Length: length,
			Bytes:  bytes,
		},
	}
}

// Version returns the version of the consent string
func (c *LazyConsent) Version() (int, error) {
	value, err := c.ReadIntField(VersionField.Offset, VersionField.NbBits)
	if err != nil {
		return 0, ErrVersionFailed(err)
	}
	return value, nil
}

// Created returns the creation date of the consent string
func (c *LazyConsent) Created() (time.Time, error) {
	value, err := c.ReadTimeField(CreatedField.Offset)
	if err != nil {
		return time.Time{}, ErrCreatedFailed(err)
	}
	return value, nil
}

// LastUpdated returns the last update date of the consent string
func (c *LazyConsent) LastUpdated() (time.Time, error) {
	value, err := c.ReadTimeField(LastUpdatedField.Offset)
	if err != nil {
		return time.Time{}, ErrLastUpdatedFailed(err)
	}
	return value, nil
}

// CMPID returns the Consent Management Platform ID
func (c *LazyConsent) CMPID() (int, error) {
	value, err := c.ReadIntField(CMPIDField.Offset, CMPIDField.NbBits)
	if err != nil {
		return 0, ErrCMPIDFailed(err)
	}
	return value, nil
}

// CMPVersion returns the Consent Management Platform version
func (c *LazyConsent) CMPVersion() (int, error) {
	value, err := c.ReadIntField(CMPVersionField.Offset, CMPVersionField.NbBits)
	if err != nil {
		return 0, ErrCMPVersionFailed(err)
	}
	return value, nil
}

// ConsentScreen returns the consent screen number
func (c *LazyConsent) ConsentScreen() (int, error) {
	value, err := c.ReadIntField(ConsentScreenField.Offset, ConsentScreenField.NbBits)
	if err != nil {
		return 0, ErrConsentScreenFailed(err)
	}
	return value, nil
}

// ConsentLanguage returns the consent language
func (c *LazyConsent) ConsentLanguage() (string, error) {
	value, err := c.ReadStringField(ConsentLanguageField.Offset, ConsentLanguageField.NbBits)
	if err != nil {
		return "", ErrConsentLanguageFailed(err)
	}
	return value, nil
}

// VendorListVersion returns the vendor list version
func (c *LazyConsent) VendorListVersion() (int, error) {
	value, err := c.ReadIntField(VendorListVersionField.Offset, VendorListVersionField.NbBits)
	if err != nil {
		return 0, ErrVendorListVersionFailed(err)
	}
	return value, nil
}

// TcfPolicyVersion returns the TCF policy version
func (c *LazyConsent) TcfPolicyVersion() (int, error) {
	value, err := c.ReadIntField(TcfPolicyVersionField.Offset, TcfPolicyVersionField.NbBits)
	if err != nil {
		return 0, ErrTcfPolicyVersionFailed(err)
	}
	return value, nil
}

// IsServiceSpecific checks if the consent is service specific
func (c *LazyConsent) IsServiceSpecific() (bool, error) {
	value, err := c.ReadBoolField(IsServiceSpecificField.Offset)
	if err != nil {
		return false, ErrIsServiceSpecificFailed(err)
	}
	return value, nil
}

// UseNonStandardStacks checks if the consent uses non standard stacks
func (c *LazyConsent) UseNonStandardStacks() (bool, error) {
	value, err := c.ReadBoolField(UseNonStandardStacksField.Offset)
	if err != nil {
		return false, ErrUseNonStandardStacksFailed(err)
	}
	return value, nil
}

// PurposeOneTreatment checks if the consent is for one treatment
func (c *LazyConsent) PurposeOneTreatment() (bool, error) {
	value, err := c.ReadBoolField(PurposeOneTreatmentField.Offset)
	if err != nil {
		return false, ErrPurposeOneTreatmentFailed(err)
	}
	return value, nil
}

// PublisherCC returns the publisher country code
func (c *LazyConsent) PublisherCC() (string, error) {
	value, err := c.ReadStringField(PublisherCCField.Offset, PublisherCCField.NbBits)
	if err != nil {
		return "", ErrPublisherCCFailed(err)
	}
	return value, nil
}

// IsRangeEncoding checks if the consent is using range encoding
func (c *LazyConsent) IsRangeEncoding() (bool, error) {
	value, err := c.ReadBoolField(IsRangeEncodingField.Offset)
	if err != nil {
		return false, ErrIsRangeEncodingFailed(err)
	}
	return value, nil
}

// NumRangeEntries returns the number of range entries
func (c *LazyConsent) NumRangeEntries() (int, error) {
	value, err := c.ReadIntField(NumRangeEntriesField.Offset, NumRangeEntriesField.NbBits)
	if err != nil {
		return 0, ErrNumEntriesFailed(err)
	}
	return value, nil
}

// MaxVendorID returns the maximum vendor ID
func (c *LazyConsent) MaxVendorID() (int, error) {
	value, err := c.ReadIntField(MaxVendorIDField.Offset, MaxVendorIDField.NbBits)
	if err != nil {
		return 0, ErrMaxVendorIDFailed(err)
	}
	return value, nil
}

// EveryPurposeAllowed checks if every purpose number is allowed
func (c *LazyConsent) EveryPurposeAllowed(numbers []int) (bool, error) {
	for _, number := range numbers {
		allowed, err := c.PurposeAllowed(number)
		if err != nil {
			return false, err
		}
		if !allowed {
			return false, nil
		}
	}
	return true, nil
}

// PurposeAllowed checks if purpose is allowed
func (c *LazyConsent) PurposeAllowed(number int) (bool, error) {
	value, err := c.readBitNumber(number, PurposesConsentField.Offset, PurposesConsentField.NbBits)
	if err != nil {
		return false, ErrPurposesConsentFailed(err)
	}
	return value, nil
}

// PurposeLITransparencyAllowed checks if purposeLITransparency is allowed
func (c *LazyConsent) PurposeLITransparencyAllowed(number int) (bool, error) {
	value, err := c.readBitNumber(number, PurposesLITransparencyField.Offset, PurposesLITransparencyField.NbBits)

	if err != nil {
		return false, ErrPurposesLITransparencyFailed(err)
	}
	return value, nil
}

// EverySpecialFeatureAllowed checks every special feature number is allowed
func (c *LazyConsent) EverySpecialFeatureAllowed(numbers []int) (bool, error) {
	for _, number := range numbers {
		allowed, err := c.SpecialFeatureAllowed(number)
		if err != nil {
			return false, err
		}
		if !allowed {
			return false, nil
		}
	}
	return true, nil
}

// SpecialFeatureAllowed checks if special feature is allowed
func (c *LazyConsent) SpecialFeatureAllowed(number int) (bool, error) {
	value, err := c.readBitNumber(number, SpecialFeatureOptInsField.Offset, SpecialFeatureOptInsField.NbBits)
	if err != nil {
		return false, ErrSpecialFeatureOptInsFailed(err)
	}
	return value, nil
}

// VendorAllowed checks if vendor is in the list of vendors user has given his consent to
func (c *LazyConsent) VendorAllowed(number int) (bool, error) {

	isRangeEncoding, err := c.IsRangeEncoding()
	if err != nil {
		return false, err
	}

	if isRangeEncoding {
		numEntries, err := c.NumRangeEntries()
		if err != nil {
			return false, err
		}

		offset := NumRangeEntriesField.NextOffset()
		for i := 0; i < int(numEntries); i++ {

			var isRange bool
			if isRange, err = c.ReadBoolField(offset); err != nil {
				return false, ErrIsRangeFailed(err)
			}
			offset += 1

			var start, end int
			if start, err = c.ReadIntField(offset, 16); err != nil {
				return false, ErrStartRangeFailed(err)
			}
			offset += 16

			if isRange {
				if end, err = c.ReadIntField(offset, 16); err != nil {
					return false, ErrEndRangeFailed(err)
				}
				offset += 16
			} else {
				end = start
			}

			if start <= number && number <= end {
				return true, nil
			}
		}

		return false, nil
	}

	maxVendorId, err := c.MaxVendorID()
	if err != nil {
		return false, err
	}

	consentedVendor, err := c.readBitNumber(number, ConsentedVendorsOffset, int(maxVendorId))
	if err != nil {
		return false, ErrConsentedVendorsFailed(err)
	}
	return consentedVendor, nil
}

// readBitNumber reads bit number as bool and checks boundaries
func (c *LazyConsent) readBitNumber(number, offset, maxNbbBits int) (bool, error) {
	if c == nil {
		return false, ErrNilBits
	}
	if number < 1 {
		return false, ErrBitLowerThanLowerBound(number, 1)
	}
	if number > maxNbbBits {
		return false, ErrBitHigherThanUpperBound(number, maxNbbBits)
	}
	value, err := c.ReadBoolField(offset + number - 1)
	if err != nil {
		return false, err
	}
	return value, nil
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
	return NewConsentField(36)
}

func NewConsentBoolField() *ConsentField {
	return NewConsentField(1)
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
