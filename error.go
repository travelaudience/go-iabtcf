package iabtcf

import (
	"fmt"
)

var (
	ErrEmptyString      = fmt.Errorf("consent string is empty")
	ErrNilBits          = fmt.Errorf("bits is nil")
	ErrNegativeBitIndex = func(index int) error {
		return fmt.Errorf("negative bit index: %d", index)
	}
	ErrBitIndexHigherThanUpperBound = func(index, upperBound int) error {
		return fmt.Errorf("bit index is higher than upper bound %d: %d", upperBound, index)
	}
	ErrBitLowerThanLowerBound = func(number, lowerBound int) error {
		return fmt.Errorf("bit #%d is lower than lower bound %d", number, lowerBound)
	}
	ErrBitHigherThanUpperBound = func(number, upperBound int) error {
		return fmt.Errorf("bit #%d is higher than upper bound %d", number, upperBound)
	}
	ErrInvalidNbBitsMultiple = func(nbBits, multiple int) error {
		return fmt.Errorf("number of bits is not multiple of %d: %d", multiple, nbBits)
	}

	ErrDecodeFailed                 = func(cause error) error { return fmt.Errorf("decode failed: %w", cause) }
	ErrVersionFailed                = func(cause error) error { return fmt.Errorf("version parse failed: %w", cause) }
	ErrCreatedFailed                = func(cause error) error { return fmt.Errorf("created parse failed: %w", cause) }
	ErrLastUpdatedFailed            = func(cause error) error { return fmt.Errorf("last updated parse failed: %w", cause) }
	ErrCMPIDFailed                  = func(cause error) error { return fmt.Errorf("cmp id parse failed: %w", cause) }
	ErrCMPVersionFailed             = func(cause error) error { return fmt.Errorf("cmp version parse failed: %w", cause) }
	ErrConsentScreenFailed          = func(cause error) error { return fmt.Errorf("consent screen parse failed: %w", cause) }
	ErrConsentLanguageFailed        = func(cause error) error { return fmt.Errorf("consent language parse failed: %w", cause) }
	ErrVendorListVersionFailed      = func(cause error) error { return fmt.Errorf("vendor list version parse failed: %w", cause) }
	ErrTcfPolicyVersionFailed       = func(cause error) error { return fmt.Errorf("tcf policy version parse failed: %w", cause) }
	ErrIsServiceSpecificFailed      = func(cause error) error { return fmt.Errorf("is service specific parse failed: %w", cause) }
	ErrUseNonStandardStacksFailed   = func(cause error) error { return fmt.Errorf("use non standard stacks parse failed: %w", cause) }
	ErrSpecialFeatureOptInsFailed   = func(cause error) error { return fmt.Errorf("special feature opt-ins parse failed: %w", cause) }
	ErrPurposesConsentFailed        = func(cause error) error { return fmt.Errorf("purposes consent parse failed: %w", cause) }
	ErrPurposesLITransparencyFailed = func(cause error) error { return fmt.Errorf("purposes li transparency parse failed: %w", cause) }
	ErrPurposeOneTreatmentFailed    = func(cause error) error { return fmt.Errorf("purpose one treatment parse failed: %w", cause) }
	ErrPublisherCCFailed            = func(cause error) error { return fmt.Errorf("publisher country code parse failed: %w", cause) }
	ErrMaxVendorIDFailed            = func(cause error) error { return fmt.Errorf("max vendor id parse failed: %w", cause) }
	ErrIsRangeEncodingFailed        = func(cause error) error { return fmt.Errorf("is range encoding parse failed: %w", cause) }
	ErrNumEntriesFailed             = func(cause error) error { return fmt.Errorf("num range entries parse failed: %w", cause) }
	ErrIsRangeFailed                = func(cause error) error { return fmt.Errorf("is range parse failed: %w", cause) }
	ErrStartRangeFailed             = func(cause error) error { return fmt.Errorf("start range parse failed: %w", cause) }
	ErrEndRangeFailed               = func(cause error) error { return fmt.Errorf("end range parse failed: %w", cause) }
	ErrRangeEntriesFailed           = func(cause error) error { return fmt.Errorf("range entries parse failed: %w", cause) }
	ErrConsentedVendorsFailed       = func(cause error) error { return fmt.Errorf("consented vendors parse failed: %w", cause) }
)
