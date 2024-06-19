package iabtcf

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// ParseCoreString parses a core string and returns a Consent object
//
// note: the consent string is base64 decoded.
// Then each field is parsed and stored in a Consent object.
// This parser is optimized for checking multiple vendors + most of the fields.
func ParseCoreString(c string) (*Consent, error) {
	if c == "" {
		return nil, fmt.Errorf("consent string is empty")
	}
	// extract core string
	cs, _, _ := strings.Cut(c, ".")

	var b, err = base64.RawURLEncoding.DecodeString(cs)
	if err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	r := NewReader(b)
	p := &Consent{}
	p.Version, err = r.ReadInt(6)
	if err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}
	p.Created, err = r.ReadTime()
	if err != nil {
		return nil, fmt.Errorf("created parse failed: %w", err)
	}
	p.LastUpdated, err = r.ReadTime()
	if err != nil {
		return nil, fmt.Errorf("last updated parse failed: %w", err)
	}
	p.CMPID, err = r.ReadInt(12)
	if err != nil {
		return nil, fmt.Errorf("cmp id parse failed: %w", err)
	}
	p.CMPVersion, err = r.ReadInt(12)
	if err != nil {
		return nil, fmt.Errorf("cmp version parse failed: %w", err)
	}
	p.ConsentScreen, err = r.ReadInt(6)
	if err != nil {
		return nil, fmt.Errorf("consent screen parse failed: %w", err)
	}
	p.ConsentLanguage, err = r.ReadString(12)
	if err != nil {
		return nil, fmt.Errorf("consent language parse failed: %w", err)
	}
	p.VendorListVersion, err = r.ReadInt(12)
	if err != nil {
		return nil, fmt.Errorf("vendor list version parse failed: %w", err)
	}
	p.TcfPolicyVersion, err = r.ReadInt(6)
	if err != nil {
		return nil, fmt.Errorf("tcf policy version parse failed: %w", err)
	}
	p.IsServiceSpecific, err = r.ReadBool()
	if err != nil {
		return nil, fmt.Errorf("is service specific parse failed: %w", err)
	}
	p.UseNonStandardStacks, err = r.ReadBool()
	if err != nil {
		return nil, fmt.Errorf("use non standard stacks parse failed: %w", err)
	}
	p.SpecialFeatureOptIns, err = r.ReadBitField(12)
	if err != nil {
		return nil, fmt.Errorf("special feature opt-ins parse failed: %w", err)
	}
	p.PurposesConsent, err = r.ReadBitField(24)
	if err != nil {
		return nil, fmt.Errorf("purposes consent parse failed: %w", err)
	}
	p.PurposesLITransparency, err = r.ReadBitField(24)
	if err != nil {
		return nil, fmt.Errorf("purposes li transparency parse failed: %w", err)
	}
	p.PurposeOneTreatment, err = r.ReadBool()
	if err != nil {
		return nil, fmt.Errorf("purpose one treatment parse failed: %w", err)
	}
	p.PublisherCC, err = r.ReadString(12)
	if err != nil {
		return nil, fmt.Errorf("publisher country code parse failed: %w", err)
	}
	p.MaxVendorID, err = r.ReadInt(16)
	if err != nil {
		return nil, fmt.Errorf("max vendor id parse failed: %w", err)
	}
	p.IsRangeEncoding, err = r.ReadBool()
	if err != nil {
		return nil, fmt.Errorf("is range encoding parse failed: %w", err)
	}

	if p.IsRangeEncoding {
		p.NumEntries, err = r.ReadInt(12)
		if err != nil {
			return nil, fmt.Errorf("num range entries parse failed: %w", err)
		}
		p.RangeEntries, err = r.ReadRangeEntries(p.NumEntries)
		if err != nil {
			return nil, fmt.Errorf("range entries parse failed: %w", err)
		}
	} else {
		p.ConsentedVendors, err = r.ReadBitField(p.MaxVendorID)
		if err != nil {
			return nil, fmt.Errorf("consented vendors parse failed: %w", err)
		}
	}

	return p, nil
}
