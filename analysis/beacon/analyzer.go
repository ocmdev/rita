package beacon

import (
	"math"
	"sort"

	dataBeacon "github.com/activecm/rita/datatypes/beacon"

	"github.com/activecm/rita/util"
)

// analyze src, dst pairs with their connection data
func (t *Beacon) analyze() {
	for data := range t.analysisChannel {
		//sort the size and timestamps since they may have arrived out of order
		sort.Sort(util.SortableInt64(data.ts))
		sort.Sort(util.SortableInt64(data.origIPBytes))

		//remove subsecond communications
		//these will appear as beacons if we do not remove them
		//subsecond beacon finding *may* be implemented later on...
		data.ts = util.RemoveSortedDuplicates(data.ts)

		//If removing duplicates lowered the conn count under the threshold,
		//remove this data from the analysis
		if len(data.ts) < t.res.Config.S.Beacon.DefaultConnectionThresh {
			continue
		}

		//store the diff slice length since we use it a lot
		//for timestamps this is one less then the data slice length
		//since we are calculating the times in between readings
		tsLength := len(data.ts) - 1
		dsLength := len(data.origIPBytes)

		//find the duration of this connection
		//perfect beacons should fill the observation period
		duration := float64(data.ts[tsLength]-data.ts[0]) /
			float64(t.maxTime-t.minTime)

		//find the delta times between the timestamps
		diff := make([]int64, tsLength)
		for i := 0; i < tsLength; i++ {
			diff[i] = data.ts[i+1] - data.ts[i]
		}

		//perfect beacons should have symmetric delta time and size distributions
		//Bowley's measure of skew is used to check symmetry
		sort.Sort(util.SortableInt64(diff))
		tsSkew := float64(0)
		dsSkew := float64(0)

		//tsLength -1 is used since diff is a zero based slice
		tsLow := diff[util.Round(.25*float64(tsLength-1))]
		tsMid := diff[util.Round(.5*float64(tsLength-1))]
		tsHigh := diff[util.Round(.75*float64(tsLength-1))]
		tsBowleyNum := tsLow + tsHigh - 2*tsMid
		tsBowleyDen := tsHigh - tsLow

		//we do the same for datasizes
		dsLow := data.origIPBytes[util.Round(.25*float64(dsLength-1))]
		dsMid := data.origIPBytes[util.Round(.5*float64(dsLength-1))]
		dsHigh := data.origIPBytes[util.Round(.75*float64(dsLength-1))]
		dsBowleyNum := dsLow + dsHigh - 2*dsMid
		dsBowleyDen := dsHigh - dsLow

		//tsSkew should equal zero if the denominator equals zero
		//bowley skew is unreliable if Q2 = Q1 or Q2 = Q3
		if tsBowleyDen != 0 && tsMid != tsLow && tsMid != tsHigh {
			tsSkew = float64(tsBowleyNum) / float64(tsBowleyDen)
		}

		if dsBowleyDen != 0 && dsMid != dsLow && dsMid != dsHigh {
			dsSkew = float64(dsBowleyNum) / float64(dsBowleyDen)
		}

		//perfect beacons should have very low dispersion around the
		//median of their delta times
		//Median Absolute Deviation About the Median
		//is used to check dispersion
		devs := make([]int64, tsLength)
		for i := 0; i < tsLength; i++ {
			devs[i] = util.Abs(diff[i] - tsMid)
		}

		dsDevs := make([]int64, dsLength)
		for i := 0; i < dsLength; i++ {
			dsDevs[i] = util.Abs(data.origIPBytes[i] - dsMid)
		}

		sort.Sort(util.SortableInt64(devs))
		sort.Sort(util.SortableInt64(dsDevs))

		tsMadm := devs[util.Round(.5*float64(tsLength-1))]
		dsMadm := dsDevs[util.Round(.5*float64(dsLength-1))]

		//Store the range for human analysis
		tsIntervalRange := diff[tsLength-1] - diff[0]
		dsRange := data.origIPBytes[dsLength-1] - data.origIPBytes[0]

		//get a list of the intervals found in the data,
		//the number of times the interval was found,
		//and the most occurring interval
		intervals, intervalCounts, tsMode, tsModeCount := createCountMap(diff)
		dsSizes, dsCounts, dsMode, dsModeCount := createCountMap(data.origIPBytes)

		//more skewed distributions recieve a lower score
		//less skewed distributions recieve a higher score
		tsSkewScore := 1.0 - math.Abs(tsSkew) //smush tsSkew
		dsSkewScore := 1.0 - math.Abs(dsSkew) //smush dsSkew

		//lower dispersion is better, cutoff dispersion scores at 30 seconds
		tsMadmScore := 1.0 - float64(tsMadm)/30.0
		if tsMadmScore < 0 {
			tsMadmScore = 0
		}

		//lower dispersion is better, cutoff dispersion scores at 32 bytes
		dsMadmScore := 1.0 - float64(dsMadm)/32.0
		if dsMadmScore < 0 {
			dsMadmScore = 0
		}

		tsDurationScore := duration

		//smaller data sizes receive a higher score
		dsSmallnessScore := 1.0 - (float64(dsMode) / 65535.0)
		if dsSmallnessScore < 0 {
			dsSmallnessScore = 0
		}

		output := dataBeacon.BeaconAnalysisOutput{
			UconnID:           data.uconnID,
			TS_iSkew:          tsSkew,
			TS_iDispersion:    tsMadm,
			TS_duration:       duration,
			TS_iRange:         tsIntervalRange,
			TS_iMode:          tsMode,
			TS_iModeCount:     tsModeCount,
			TS_intervals:      intervals,
			TS_intervalCounts: intervalCounts,
			DS_skew:           dsSkew,
			DS_dispersion:     dsMadm,
			DS_range:          dsRange,
			DS_sizes:          dsSizes,
			DS_sizeCounts:     dsCounts,
			DS_mode:           dsMode,
			DS_modeCount:      dsModeCount,
		}

		//score numerators
		tsSum := (tsSkewScore + tsMadmScore + tsDurationScore)
		dsSum := (dsSkewScore + dsMadmScore + dsSmallnessScore)

		//score averages
		output.TS_score = tsSum / 3.0
		output.DS_score = dsSum / 3.0
		output.Score = (tsSum + dsSum) / 6.0

		t.writeChannel <- &output
	}
	t.analysisWg.Done()
}

// createCountMap returns a distinct data array, data count array, the mode,
// and the number of times the mode occured
func createCountMap(data []int64) ([]int64, []int64, int64, int64) {
	//create interval counts for human analysis
	dataMap := make(map[int64]int64)
	for _, d := range data {
		dataMap[d]++
	}

	distinct := make([]int64, len(dataMap))
	counts := make([]int64, len(dataMap))

	i := 0
	for k, v := range dataMap {
		distinct[i] = k
		counts[i] = v
		i++
	}

	mode := distinct[0]
	max := counts[0]
	for idx, count := range counts {
		if count > max {
			max = count
			mode = distinct[idx]
		}
	}
	return distinct, counts, mode, max
}
