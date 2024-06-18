package iabtcf

import (
	"encoding/base64"
	"strings"
)

// ParseCoreString parses a core string and returns a Consent object
//
// note: the consent string is base64 decoded.
// Then each field is parsed and stored in a Consent object.
// This parser is optimized for checking multiple vendors + most of the fields.
func ParseCoreString(c string) (*Consent, error) {
	if c == "" {
		return nil, ErrEmptyString
	}
	// extract core string
	cs, _, _ := strings.Cut(c, ".")

	var b, err = base64.RawURLEncoding.DecodeString(cs)
	if err != nil {
		return nil, ErrDecodeFailed(err)
	}

	r := NewReader(b)
	p := &Consent{}
	p.Version, err = r.ReadInt(6)
	if err != nil {
		return nil, ErrVersionFailed(err)
	}
	p.Created, err = r.ReadTime()
	if err != nil {
		return nil, ErrCreatedFailed(err)
	}
	p.LastUpdated, err = r.ReadTime()
	if err != nil {
		return nil, ErrLastUpdatedFailed(err)
	}
	p.CMPID, err = r.ReadInt(12)
	if err != nil {
		return nil, ErrCMPIDFailed(err)
	}
	p.CMPVersion, err = r.ReadInt(12)
	if err != nil {
		return nil, ErrCMPVersionFailed(err)
	}
	p.ConsentScreen, err = r.ReadInt(6)
	if err != nil {
		return nil, ErrConsentScreenFailed(err)
	}
	p.ConsentLanguage, err = r.ReadString(12)
	if err != nil {
		return nil, ErrConsentLanguageFailed(err)
	}
	p.VendorListVersion, err = r.ReadInt(12)
	if err != nil {
		return nil, ErrVendorListVersionFailed(err)
	}
	p.TcfPolicyVersion, err = r.ReadInt(6)
	if err != nil {
		return nil, ErrTcfPolicyVersionFailed(err)
	}
	p.IsServiceSpecific, err = r.ReadBool()
	if err != nil {
		return nil, ErrIsServiceSpecificFailed(err)
	}
	p.UseNonStandardStacks, err = r.ReadBool()
	if err != nil {
		return nil, ErrUseNonStandardStacksFailed(err)
	}
	p.SpecialFeatureOptIns, err = r.ReadBitField(12)
	if err != nil {
		return nil, ErrSpecialFeatureOptInsFailed(err)
	}
	p.PurposesConsent, err = r.ReadBitField(24)
	if err != nil {
		return nil, ErrPurposesConsentFailed(err)
	}
	p.PurposesLITransparency, err = r.ReadBitField(24)
	if err != nil {
		return nil, ErrPurposesLITransparencyFailed(err)
	}
	p.PurposeOneTreatment, err = r.ReadBool()
	if err != nil {
		return nil, ErrPurposeOneTreatmentFailed(err)
	}
	p.PublisherCC, err = r.ReadString(12)
	if err != nil {
		return nil, ErrPublisherCCFailed(err)
	}
	p.MaxVendorID, err = r.ReadInt(16)
	if err != nil {
		return nil, ErrMaxVendorIDFailed(err)
	}
	p.IsRangeEncoding, err = r.ReadBool()
	if err != nil {
		return nil, ErrIsRangeEncodingFailed(err)
	}

	if p.IsRangeEncoding {
		p.NumEntries, err = r.ReadInt(12)
		if err != nil {
			return nil, ErrNumEntriesFailed(err)
		}
		p.RangeEntries, err = r.ReadRangeEntries(p.NumEntries)
		if err != nil {
			return nil, ErrRangeEntriesFailed(err)
		}
	} else {
		p.ConsentedVendors, err = r.ReadBitField(p.MaxVendorID)
		if err != nil {
			return nil, ErrConsentedVendorsFailed(err)
		}
	}

	return p, nil
}
