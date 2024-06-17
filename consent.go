package iabtcf

import (
	"fmt"
	"time"

	"github.com/rupertchen/go-bits"
)

// Consent represents Core Consent extracted from an IAB Consent String v2.0
// It's implemented according to specification: https://github.com/InteractiveAdvertisingBureau/GDPR-Transparency-and-Consent-Framework/blob/master/TCFv2/IAB%20Tech%20Lab%20-%20Consent%20string%20and%20vendor%20list%20formats%20v2.md
type Consent struct {
	Version                int
	Created                time.Time
	LastUpdated            time.Time
	CMPID                  int
	CMPVersion             int
	ConsentScreen          int
	ConsentLanguage        string
	VendorListVersion      int
	TcfPolicyVersion       int
	IsServiceSpecific      bool
	UseNonStandardStacks   bool
	SpecialFeatureOptIns   *Bits
	PurposesConsent        *Bits
	PurposesLITransparency *Bits
	PurposeOneTreatment    bool
	PublisherCC            string
	MaxVendorID            int
	IsRangeEncoding        bool
	ConsentedVendors       *Bits
	NumEntries             int
	RangeEntries           []*RangeEntry
}

// RangeEntry defines a range groups of Vendor IDs who have been disclosed to a user
type RangeEntry struct {
	StartOrOnlyVendorId int
	EndVendorID         int
}

// BitField defines a bitmap with a specific length
// note: bitmap is not storing the real length of the bitmap, but the number of blocks
// note: bits are right aligned in the bit map > need to compute an offset to get the right bit ( see HasBitNumber )
type BitField struct {
	*bits.Bitmap
	Length int
}

// EveryPurposeAllowed returns true if every purpose number is allowed in
// the ParsedConsent, otherwise false
func (p *Consent) EveryPurposeAllowed(numbers []int) bool {
	for _, number := range numbers {
		if !p.PurposeAllowed(number) {
			return false
		}
	}
	return true
}

// PurposeAllowed checks if purpose is allowed in the ParsedConsent
func (p *Consent) PurposeAllowed(number int) bool {
	return p.PurposesConsent.HasBit(number)
}

// PurposeLITransparencyAllowed checks if purposeLITransparency is allowed in the ParsedConsent
func (p *Consent) PurposeLITransparencyAllowed(number int) bool {
	return p.PurposesLITransparency.HasBit(number)
}

// EverySpecialFeatureAllowed returns true if every special feature number is allowed in
// the ParsedConsent, otherwise false
func (p *Consent) EverySpecialFeatureAllowed(numbers []int) bool {
	for _, number := range numbers {
		if !p.SpecialFeatureAllowed(number) {
			return false
		}
	}
	return true
}

// SpecialFeatureAllowed checks if special feature is allowed in the ParsedConsent
func (p *Consent) SpecialFeatureAllowed(number int) bool {
	return p.SpecialFeatureOptIns.HasBit(number)
}

// VendorAllowed checks if vendor is in the list of vendors user has given his consent to
func (p *Consent) VendorAllowed(number int) bool {

	if p.IsRangeEncoding {
		for _, e := range p.RangeEntries {
			if e.StartOrOnlyVendorId <= number && number <= e.EndVendorID {
				return true
			}
		}
		return false
	}

	return p.ConsentedVendors.HasBit(number)
}

func (b *BitField) HasBitNumber(number int) bool {
	if b == nil || number < 1 || number > b.Length {
		return false
	}
	// note: bits are right aligned
	offset := b.Size() - b.Length
	block, err := b.Get(uint(offset+number-1), 1)
	if err != nil {
		return false
	}
	return block == 1
}

func (b *BitField) ToString() string {
	if b == nil {
		return ""
	}
	result := ""
	for number := 1; number <= b.Length; number++ {
		if number > 1 && number%10 == 1 {
			result += " "
		}
		if b.HasBitNumber(number) {
			result += "1"
		} else {
			result += "0"
		}
	}
	nb := uint(b.Size() / 8)
	result += fmt.Sprintf(" [%d](", nb)
	for i := uint(0); i < nb; i++ {
		block, err := b.Get(i*8, 8)
		if err != nil {
			result += fmt.Sprintf("(%v)", err)
		} else {
			result += fmt.Sprintf("%d > %s", block, byteToBitString(byte(block)))
		}
		if i != 0 {
			result += "|"
		}
	}
	result += ")"
	return result
}

func blockToBitString(block uint64) string {
	result := ""
	tmp := block
	for i := 0; i < 8; i++ {
		if i != 0 {
			result = byteToBitString(uint8(tmp)) + " " + result
		} else {
			result = byteToBitString(uint8(tmp))
		}
		tmp = tmp >> 8
	}
	return fmt.Sprintf("%d=%s", block, result)
}

func byteToBitString(b byte) string {
	result := ""
	for i := 0; i < 8; i++ {
		if b&(1<<i) != 0 {
			result = "1" + result
		} else {
			result = "0" + result
		}
	}
	return result
}
