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

			gotVersion, err := gotParsed.Version()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.Version, gotVersion)

			gotCreated, err := gotParsed.Created()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.Created, gotCreated)

			gotLastUpdated, err := gotParsed.LastUpdated()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.LastUpdated, gotLastUpdated)

			gotCMPID, err := gotParsed.CMPID()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.CMPID, gotCMPID)

			gotCMPVersion, err := gotParsed.CMPVersion()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.CMPVersion, gotCMPVersion)

			gotConsentScreen, err := gotParsed.ConsentScreen()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.ConsentScreen, gotConsentScreen)

			gotConsentLanguage, err := gotParsed.ConsentLanguage()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.ConsentLanguage, gotConsentLanguage)

			gotVendorListVersion, err := gotParsed.VendorListVersion()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.VendorListVersion, gotVendorListVersion)

			gotTcfPolicyVersion, err := gotParsed.TcfPolicyVersion()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.TcfPolicyVersion, gotTcfPolicyVersion)

			gotIsServiceSpecific, err := gotParsed.IsServiceSpecific()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.IsServiceSpecific, gotIsServiceSpecific)

			gotUseNonStandardStacks, err := gotParsed.UseNonStandardStacks()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.UseNonStandardStacks, gotUseNonStandardStacks)

			gotPurposeOneTreatment, err := gotParsed.PurposeOneTreatment()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.PurposeOneTreatment, gotPurposeOneTreatment)

			gotPublisherCC, err := gotParsed.PublisherCC()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.PublisherCC, gotPublisherCC)

			gotIsRangeEncoding, err := gotParsed.IsRangeEncoding()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.IsRangeEncoding, gotIsRangeEncoding)

			for number := 1; number <= wantParsed.SpecialFeatureOptIns.Length; number++ {
				gotAllowed, err := gotParsed.SpecialFeatureAllowed(number)
				require.NoError(t, err, "unexpected error")
				require.Equal(t, wantParsed.SpecialFeatureAllowed(number), gotAllowed)
			}

			for number := 1; number <= wantParsed.PurposesConsent.Length; number++ {
				gotAllowed, err := gotParsed.PurposeAllowed(number)
				require.NoError(t, err, "unexpected error")
				require.Equal(t, wantParsed.PurposeAllowed(number), gotAllowed)
			}

			for number := 1; number <= wantParsed.PurposesLITransparency.Length; number++ {
				gotAllowed, err := gotParsed.PurposeLITransparencyAllowed(number)
				require.NoError(t, err, "unexpected error")
				require.Equal(t, wantParsed.PurposeLITransparencyAllowed(number), gotAllowed)
			}

			if wantParsed.IsRangeEncoding {
				gotNumRangeEntries, err := gotParsed.NumRangeEntries()
				require.NoError(t, err, "unexpected error")
				require.Equal(t, wantParsed.NumEntries, gotNumRangeEntries)
			}

			gotMaxVendorID, err := gotParsed.MaxVendorID()
			require.NoError(t, err, "unexpected error")
			require.Equal(t, wantParsed.MaxVendorID, gotMaxVendorID)

			for number := 1; number <= wantParsed.MaxVendorID; number++ {
				gotAllowed, err := gotParsed.VendorAllowed(number)
				require.NoError(t, err, "unexpected error")
				require.Equal(t, wantParsed.VendorAllowed(number), gotAllowed)
			}
		})
	}

}
