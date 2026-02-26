package iabtcf

// HasDisclosedVendorsBlock returns true if there is at least one disclosedVendors block in the consent string.
func (c *LazyConsent) HasDisclosedVendorsBlock() bool {

	for _, block := range c.Extras {
		blockType := block.ReadIntField(0, 3)
		if blockType == 1 {
			return true
		}
	}
	return false

}

// IsVendorDisclosed examines all of the disclosedVendor blocks and returns true if the given vendor ID is found
// in any of them.  Returns false if there are no disclosed vendor blocks in the consent string.
//
// See also https://github.com/InteractiveAdvertisingBureau/GDPR-Transparency-and-Consent-Framework/blob/master/TCFv2/IAB%20Tech%20Lab%20-%20Consent%20string%20and%20vendor%20list%20formats%20v2.md#disclosed-vendors
func (c *LazyConsent) IsVendorDisclosed(vendorID int) bool {

	if vendorID <= 0 {
		return false
	}

	for _, block := range c.Extras {
		blockType := block.ReadIntField(0, 3)
		if blockType != 1 {
			// not a disclosedVendor block
			continue
		}

		maxVendorID := block.ReadIntField(3, 16)
		if vendorID > maxVendorID {
			continue
		}

		isRangeEncoding := block.ReadBoolField(19)
		if isRangeEncoding {
			// range encoding
			numEntries := block.ReadIntField(20, 12)
			offset := 32
			for range numEntries {
				isaRange := block.ReadBoolField(offset)
				offset++
				startID := block.ReadIntField(offset, 16)
				offset += 16
				if vendorID == startID {
					// found vendor ID
					return true
				}
				if isaRange {
					endID := block.ReadIntField(offset, 16)
					offset += 16
					if vendorID > startID && vendorID <= endID {
						// vendor ID is in a range
						return true
					}
				}
			}
		} else {
			// Bit field encoding: bit fields start at offset 20, and vendor ID starts at 1
			if block.ReadBoolField(19 + vendorID) {
				// found vendor ID
				return true
			}
		}
	}

	return false
}
