package iabtcf

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMemory(t *testing.T) {

	nbRoutine := 50
	nbParse := 5000
	forceToFailed := false
	lazy := true

	const TravelAudienceVendorID = 423
	var PurposesList = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var SpecialFeaturesList = []int{1}

	filename := "mem.pprof"
	now := func() string { return time.Now().Format("15:04:05.000") }

	fmt.Printf("[%s] remove %s \n", now(), filename)
	_ = os.Remove(filename)

	fmt.Printf("[%s] create %s \n", now(), filename)
	f, err := os.Create(filename)
	if err != nil {
		t.Errorf("failed to create %s: %v", filename, err)
	}
	defer f.Close()

	fmt.Printf("[%s] gc \n", now())
	runtime.GC()

	// note: this consent string contains more than 4000 vendors
	c := "CP9Qr_AP9Qr_AAfETDFRAwEsAP_gAEPgAAigg1NX_H__bX9v-Xr36ft0eY1f99j77uQxBhfJs-4FzLvW_JwX32EzNE36tqYKmRIEu3bBIQFtHJnUTVihaogVrzHsYkGchTNKJ-BkiHMRe2dYCF5vmYtj-QKZ5_p_d3f52T_9_dv-3dzzz91nv3f9f-f1eLida59tH_v_bRKb-_If9_7-_4v0_t_rk2_eTVv_9evv79-u_t____9_9____4AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAEQamr_j__tr-3_L179P26PMav--x993IYgwvk2fcC5l3rfk4L77CZmib9W1MFTIkCXbtgkIC2jkzqJqxQtUQK15j2MSDOQpmlE_AyRDmIvbOsBC83zMWx_IFM8_0_u7v87J_-_u3_bu555-6z37v-v_P6vFxOtc-2j_3_tolN_fkP-_9_f8X6f2_1ybfvJq3_-vX39-_Xf2____-_-____8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIAAACAA"

	var wg sync.WaitGroup

	fmt.Printf("[%s] start \n", now())
	ch := make(chan struct{})

	wg.Add(nbRoutine)
	for i := 0; i < nbRoutine; i++ {
		go func(i int, c string) {
			defer wg.Done()
			start := time.Now()
			for j := 0; j < nbParse; j++ {
				if lazy {
					parsed, err := LazyParseCoreString(c)
					if err != nil {
						fmt.Printf("[error] unable to parse consent string: %s \n", err.Error())
					} else {
						allowed := parsed.VendorAllowed(TravelAudienceVendorID)
						purposesAllowed := parsed.EveryPurposeAllowed(PurposesList)
						specialFeaturesAllowed := parsed.EveryPurposeAllowed(SpecialFeaturesList)
						if !allowed {
							fmt.Printf("[error] travel audience not consented: purposesAllowed=%t, specialFeaturesAllowed=%t \n", purposesAllowed, specialFeaturesAllowed)
						}
					}
				} else {
					parsed, err := ParseCoreString(c)
					if err != nil {
						fmt.Printf("[error] unable to parse consent string: %s \n", err.Error())
					} else {
						allowed := parsed.VendorAllowed(TravelAudienceVendorID)
						purposesAllowed := parsed.EveryPurposeAllowed(PurposesList)
						specialFeaturesAllowed := parsed.EveryPurposeAllowed(SpecialFeaturesList)
						if !allowed {
							fmt.Printf("[error] travel audience not consented: purposesAllowed=%t, specialFeaturesAllowed=%t \n", purposesAllowed, specialFeaturesAllowed)
						}
					}
				}
				ch <- struct{}{}
			}
			elapsed := time.Since(start)
			fmt.Printf("[%s] routine #%3d > elapsed: %v, avg: %v \n", now(), i, elapsed, elapsed/time.Duration(nbParse))
		}(i, c)
	}

	count := 0
	stop := nbRoutine * nbParse
	write := stop * 8 / 10
	fmt.Printf("[%s] wait for 80%% : write=%d, stop=%d \n", now(), write, stop)
	for {
		<-ch
		count++

		if count == write {
			fmt.Printf("[%s] count %d >>> WRITE %s \n", now(), count, filename)
			err = pprof.WriteHeapProfile(f)
			if err != nil {
				t.Errorf("failed to write %s: %v", filename, err)
			}
		}
		if count == stop {
			fmt.Printf("[%s] count %d >>> STOP \n", now(), count)
			break
		}
	}

	fmt.Printf("[%s] closing channel... \n", now())
	close(ch)

	fmt.Printf("[%s] wait routines... \n", now())
	wg.Wait()
	fmt.Printf("[%s] end \n", now())

	require.False(t, forceToFailed, "force to failed")
}
