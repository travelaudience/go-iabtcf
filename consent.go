package iabtcf

import (
	"time"
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
