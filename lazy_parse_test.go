package iabtcf

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLazyParse(t *testing.T) {

	type TestCase struct {
		consent string
	}

	testCases := map[string]*TestCase{
		"v2-small": {
			consent: "COzcJxTOzcJxTBcAAAENAiCMAP_AAAAAAAAADTwAQDTgAAAA.IF5EX2S5OI2tho2YdF7BEYYwfJxyigMgShgQIsS8NwIeFbBoGPmAAHBG4JAQAGBAkkACBAQIsHGBcCQABgIgRiRCMQEGMjzNKBJBAggkbI0FACCVmnkHS3ZCY70-6u__bA",
		},
		"range-encoding": {
			consent: "COzcJxTOzcJxTBcAAAENAiCMAP_AAAAAAAAADTwAQDTgAAAA.IF5EX2S5OI2tho2YdF7BEYYwfJxyigMgShgQIsS8NwIeFbBoGPmAAHBG4JAQAGBAkkACBAQIsHGBcCQABgIgRiRCMQEGMjzNKBJBAggkbI0FACCVmnkHS3ZCY70-6u__bA",
		},
		"v2-big": {
			consent: "CP9Qr_AP9Qr_AAfETDFRAwEsAP_gAEPgAAigg1NX_H__bX9v-Xr36ft0eY1f99j77uQxBhfJs-4FzLvW_JwX32EzNE36tqYKmRIEu3bBIQFtHJnUTVihaogVrzHsYkGchTNKJ-BkiHMRe2dYCF5vmYtj-QKZ5_p_d3f52T_9_dv-3dzzz91nv3f9f-f1eLida59tH_v_bRKb-_If9_7-_4v0_t_rk2_eTVv_9evv79-u_t____9_9____4AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAEQamr_j__tr-3_L179P26PMav--x993IYgwvk2fcC5l3rfk4L77CZmib9W1MFTIkCXbtgkIC2jkzqJqxQtUQK15j2MSDOQpmlE_AyRDmIvbOsBC83zMWx_IFM8_0_u7v87J_-_u3_bu555-6z37v-v_P6vFxOtc-2j_3_tolN_fkP-_9_f8X6f2_1ybfvJq3_-vX39-_Xf2____-_-____8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIAAACAA",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Helper()
			fmt.Printf("\n[test] ---------- %s ---------- \n", name)

			// comparing normal parsing with lazy parsing
			// assumption: normal parsing is the expected result
			wantParsed, wantErr := ParseCoreString(tc.consent)
			gotParsed, gotErr := LazyParseCoreString(tc.consent)

			fmt.Printf("\n[test] bits: %s \n", gotParsed.ToBitString())

			require.Equal(t, wantErr, gotErr)

			require.Equal(t, wantParsed.Version, gotParsed.Version())
			require.Equal(t, wantParsed.Created, gotParsed.Created())
			require.Equal(t, wantParsed.LastUpdated, gotParsed.LastUpdated())
			require.Equal(t, wantParsed.CMPID, gotParsed.CMPID())
			require.Equal(t, wantParsed.CMPVersion, gotParsed.CMPVersion())
			require.Equal(t, wantParsed.ConsentScreen, gotParsed.ConsentScreen())
			require.Equal(t, wantParsed.ConsentLanguage, gotParsed.ConsentLanguage())
			require.Equal(t, wantParsed.VendorListVersion, gotParsed.VendorListVersion())
			require.Equal(t, wantParsed.TcfPolicyVersion, gotParsed.TcfPolicyVersion())
			require.Equal(t, wantParsed.IsServiceSpecific, gotParsed.IsServiceSpecific())
			require.Equal(t, wantParsed.UseNonStandardStacks, gotParsed.UseNonStandardStacks())
			require.Equal(t, wantParsed.PurposeOneTreatment, gotParsed.PurposeOneTreatment())
			require.Equal(t, wantParsed.PublisherCC, gotParsed.PublisherCC())
			require.Equal(t, wantParsed.IsRangeEncoding, gotParsed.IsRangeEncoding())

			for number := 1; number <= SpecialFeatureOptInsField.NbBits; number++ {
				require.Equal(t, wantParsed.SpecialFeatureAllowed(number), gotParsed.SpecialFeatureAllowed(number), "special feature %d", number)
			}
			for number := 1; number <= PurposesConsentField.NbBits; number++ {
				require.Equal(t, wantParsed.PurposeAllowed(number), gotParsed.PurposeAllowed(number), "purpose %d", number)
			}
			for number := 1; number <= PurposesLITransparencyField.NbBits; number++ {
				require.Equal(t, wantParsed.PurposeLITransparencyAllowed(number), gotParsed.PurposeLITransparencyAllowed(number), "purpose LIT %d", number)
			}
			require.Equal(t, wantParsed.MaxVendorID, gotParsed.MaxVendorID())
			if wantParsed.IsRangeEncoding {
				require.Equal(t, wantParsed.NumEntries, gotParsed.NumRangeEntries())
			}
			for number := 1; number <= wantParsed.MaxVendorID; number++ {
				require.Equal(t, wantParsed.VendorAllowed(number), gotParsed.VendorAllowed(number), "vendor %d", number)
			}
		})
	}

}
