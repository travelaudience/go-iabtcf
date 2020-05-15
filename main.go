package iabtcf

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func ParseCoreString(c string) (*Consent, error) {
	if c == "" {
		return nil, errors.New("string is empty")
	}
	// extract core string
	cs := strings.Split(c, ".")[0]

	var b, err = base64.RawURLEncoding.DecodeString(cs)
	if err != nil {
		return nil, fmt.Errorf("DecodeString failed: %s", err.Error())
	}

	r := NewReader(b)
	p := &Consent{}
	p.Version, err = r.ReadInt(6)
	if err != nil {
		return nil, fmt.Errorf("Version parse failed: %s", err.Error())
	}
	p.Created, err = r.ReadTime()
	if err != nil {
		return nil, fmt.Errorf("Created parse failed: %s", err.Error())
	}
	p.LastUpdated, err = r.ReadTime()
	if err != nil {
		return nil, fmt.Errorf("LastUpdated parse failed: %s", err.Error())
	}
	p.CMPID, err = r.ReadInt(12)
	if err != nil {
		return nil, fmt.Errorf("CMPID parse failed: %s", err.Error())
	}
	p.CMPVersion, err = r.ReadInt(12)
	if err != nil {
		return nil, fmt.Errorf("CMPVersion parse failed: %s", err.Error())
	}
	p.ConsentScreen, err = r.ReadInt(6)
	if err != nil {
		return nil, fmt.Errorf("ConsentScreen parse failed: %s", err.Error())
	}
	p.ConsentLanguage, err = r.ReadString(12)
	if err != nil {
		return nil, fmt.Errorf("ConsentLanguage parse failed: %s", err.Error())
	}
	p.VendorListVersion, err = r.ReadInt(12)
	if err != nil {
		return nil, fmt.Errorf("VendorListVersion parse failed: %s", err.Error())
	}
	p.TcfPolicyVersion, err = r.ReadInt(6)
	if err != nil {
		return nil, fmt.Errorf("TcfPolicyVersion parse failed: %s", err.Error())
	}
	p.IsServiceSpecific, err = r.ReadBool()
	if err != nil {
		return nil, fmt.Errorf("IsServiceSpecific parse failed: %s", err.Error())
	}
	p.UseNonStandardStacks, err = r.ReadBool()
	if err != nil {
		return nil, fmt.Errorf("UseNonStandardStacks parse failed: %s", err.Error())
	}
	p.SpecialFeatureOptIns, err = r.ReadBitField(12)
	if err != nil {
		return nil, fmt.Errorf("SpecialFeatureOptIns parse failed: %s", err.Error())
	}
	p.PurposesConsent, err = r.ReadBitField(24)
	if err != nil {
		return nil, fmt.Errorf("PurposesConsent parse failed: %s", err.Error())
	}
	p.PurposesLITransparency, err = r.ReadBitField(24)
	if err != nil {
		return nil, fmt.Errorf("PurposesLITransparency parse failed: %s", err.Error())
	}
	p.PurposeOneTreatment, err = r.ReadBool()
	if err != nil {
		return nil, fmt.Errorf("PurposeOneTreatment parse failed: %s", err.Error())
	}
	p.PublisherCC, err = r.ReadString(12)
	if err != nil {
		return nil, fmt.Errorf("PublisherCC parse failed: %s", err.Error())
	}
	p.MaxVendorID, err = r.ReadInt(16)
	if err != nil {
		return nil, fmt.Errorf("MaxVendorID parse failed: %s", err.Error())
	}
	p.IsRangeEncoding, err = r.ReadBool()
	if err != nil {
		return nil, fmt.Errorf("IsRangeEncoding parse failed: %s", err.Error())
	}

	if p.IsRangeEncoding {
		p.NumEntries, err = r.ReadInt(12)
		if err != nil {
			return nil, fmt.Errorf("NumEntries parse failed: %s", err.Error())
		}
		p.RangeEntries, err = r.ReadRangeEntries(uint(p.NumEntries))
		if err != nil {
			return nil, fmt.Errorf("RangeEntries parse failed: %s", err.Error())
		}
	} else {
		p.ConsentedVendors, err = r.ReadBitField(uint(p.MaxVendorID))
		if err != nil {
			return nil, fmt.Errorf("ConsentedVendors parse failed: %s", err.Error())
		}
	}

	return p, nil
}
