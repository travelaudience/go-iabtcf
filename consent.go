package iabtcf

import "time"

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
	SpecialFeatureOptIns   map[int]bool
	PurposesConsent        map[int]bool
	PurposesLITransparency map[int]bool
	PurposeOneTreatment    bool
	PublisherCC            string
	MaxVendorID            int
	IsRangeEncoding        bool
	ConsentedVendors       map[int]bool
	NumEntries             int
	RangeEntries           []*RangeEntry
}

// RangeEntry defines a range groups of Vendor IDs who have been disclosed to a user
type RangeEntry struct {
	StartOrOnlyVendorId int
	EndVendorID         int
}

// EveryPurposeAllowed returns true if every purpose number is allowed in
//// the ParsedConsent, otherwise false
func (p *Consent) EveryPurposeAllowed(ps []int) bool {
	for _, rp := range ps {
		if !p.PurposesConsent[rp] {
			return false
		}
	}
	return true
}

// EverySpecialFeatureAllowed returns true if every special feature number is allowed in
// the ParsedConsent, otherwise false
func (p *Consent) EverySpecialFeatureAllowed(ps []int) bool {
	for _, rp := range ps {
		if !p.SpecialFeatureOptIns[rp] {
			return false
		}
	}
	return true
}

// VendorAllowed checks if vendor is in the list of vendors user has given his consent to
func (p *Consent) VendorAllowed(v int) bool {
	if p.IsRangeEncoding {
		for _, e := range p.RangeEntries {
			if e.StartOrOnlyVendorId <= v && v <= e.EndVendorID {
				return true
			}
		}
		return false
	}

	return p.ConsentedVendors[v]
}
