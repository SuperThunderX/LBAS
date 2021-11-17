package main

import (
	"math"
	"sort"

	"github.com/digisan/gotk/slice/tf64"
	"github.com/digisan/gotk/slice/ti"
	"github.com/digisan/gotk/slice/tu8i"
)

func abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

func smooth(pts []int) (ret []int) {
	ret = make([]int, len(pts))
	copy(ret, pts[:4])
	copy(ret[len(ret)-4:], pts[len(pts)-4:])
	for i := 4; i < len(pts)-4; i++ {
		ret[i] = ((-21)*pts[i-4] + 14*pts[i-3] + 39*pts[i-2] + 54*pts[i-1] + 59*pts[i] + 54*pts[i+1] + 39*pts[i+2] + 14*pts[i+3] - 21*pts[i+4]) / 231
	}
	for i := 0; i < len(ret); i++ {
		if ret[i] < 0 {
			ret[i] = 0
		}
	}
	return
}

func histogram(bytes []byte) (m map[byte]int, maxIdx byte, maxCnt int) {
	m = make(map[byte]int)
	for i := 0; i < 256; i++ {
		m[byte(i)] = 0
	}
	for i := 0; i < len(bytes); i++ {
		v := bytes[i]
		m[v]++
		if m[v] > maxCnt {
			maxCnt = m[v]
			maxIdx = v
		}
	}
	return
}

func derivative1(data []int) (ret []float64) {
	temp := make([]float64, len(data))
	for i := 0; i < len(data); i++ {
		temp[i] = float64(data[i])
	}
	ret = make([]float64, len(temp))
	for i := 4; i < len(temp)-4; i++ {
		ret[i] = (86*temp[i-4] - 142*temp[i-3] - 193*temp[i-2] - 126*temp[i-1] + 126*temp[i+1] + 193*temp[i+2] + 142*temp[i+3] - 86*temp[i+4]) / 1188
	}
	return
}

func maxmin(data ...int) (max, min, maxabs int) {
	max = ti.Max(data...)
	min = ti.Min(data...)
	abs1, abs2 := max, min
	if abs1 < 0 {
		abs1 = -abs1
	}
	if abs2 < 0 {
		abs2 = -abs2
	}
	maxabs = ti.Max(abs1, abs2)
	return
}

func minabs(data ...float64) float64 {
	temp := make([]float64, len(data))
	copy(temp, data)
	temp = tf64.FM(temp, nil, func(i int, e float64) float64 {
		if e < 0 {
			return -e
		}
		return e
	})
	return tf64.Min(temp...)
}

func FacIsPeak(data []int, halfstep int) func(i int) bool {

	var (
		deri  = derivative1(data)
		E     = 0.1
		iPeak = []int{}
	)

	return func(i int) bool {
		if i >= halfstep && i < len(deri)-halfstep {

			nUp1, nEven1, nDown1 := 0, 0, 0
			for j := i - halfstep; j < i; j++ {
				switch {
				case deri[j] > E: // up
					nUp1++
				case deri[j] >= -E && deri[j] <= E: // top
					nEven1++
				default: // down
					nDown1++
				}
			}
			if nDown1 > 0 || (nUp1 <= nEven1) {
				return false
			}

			nUp2, nEven2, nDown2 := 0, 0, 0
			for j := i + halfstep; j > i; j-- {
				switch {
				case deri[j] < -E: // down
					nDown2++
				case deri[j] >= -E && deri[j] <= E: // top
					nEven2++
				default: // up
					nUp2++
				}
			}
			if nUp2 > 0 || (nDown2 <= nEven2) {
				return false
			}

			e := minabs(deri[i-halfstep : i+halfstep+1]...)
			if deri[i] == e || deri[i] == -e {
				if len(iPeak) > 0 && i <= iPeak[len(iPeak)-1]+halfstep {
					return false
				}
				iPeak = append(iPeak, i)
				return true
			}
		}
		return false
	}
}

func FacIsBottom(data []int, halfstep int) func(i int) bool {

	var (
		deri    = derivative1(data)
		E       = 0.1
		iBottom = []int{}
	)

	return func(i int) bool {
		if i >= halfstep && i < len(deri)-halfstep {

			nUp1, nEven1, nDown1 := 0, 0, 0
			for j := i - halfstep; j < i; j++ {
				switch {
				case deri[j] < -E: // down
					nDown1++
				case deri[j] >= -E && deri[j] <= E: // bottom
					nEven1++
				default: // up
					nUp1++
				}
			}
			if nUp1 > 0 || (nDown1 <= nEven1) {
				return false
			}

			nUp2, nEven2, nDown2 := 0, 0, 0
			for j := i + halfstep; j > i; j-- {
				switch {
				case deri[j] > E: // up
					nUp2++
				case deri[j] >= -E && deri[j] <= E: // bottom
					nEven2++
				default: // down
					nDown2++
				}
			}
			if nDown2 > 0 || (nUp2 <= nEven2) {
				return false
			}

			e := minabs(deri[i-halfstep : i+halfstep+1]...)
			if deri[i] == e || deri[i] == -e {
				if len(iBottom) > 0 && i <= iBottom[len(iBottom)-1]+halfstep {
					return false
				}
				iBottom = append(iBottom, i)
				return true
			}
		}
		return false
	}
}

func Peaks(data map[byte]int, halfstep, nSmooth, n int) map[byte]int {
	m := make(map[byte]int)
	ks, vs := tu8i.Map2KVs(data, func(i, j byte) bool { return i < j }, nil)
	for i := 0; i < nSmooth; i++ {
		vs = smooth(vs)
	}
	isPeak := FacIsPeak(vs, halfstep)
	for i := 0; i < len(vs); i++ {
		if isPeak(i) {
			m[ks[i]] = vs[i]
		}
	}

	// adjust to max value
	mp := make(map[byte]int)
	for k, v := range m {
		if max, n := ti.MaxIdx(vs[k-1 : k+2]...); max > v {
			mp[k-1+byte(n)] = max
		} else {
			mp[k] = v
		}
	}

	if n < 0 {
		return mp
	}

	ks, vs = tu8i.Map2KVs(mp, nil, func(i, j int) bool { return j < i })
	mp = make(map[byte]int)
	for i := 0; i < n && i < len(ks); i++ {
		mp[ks[i]] = vs[i]
	}
	return mp
}

func Bottoms(data map[byte]int, halfstep, nSmooth, n int) map[byte]int {
	m := make(map[byte]int)
	ks, vs := tu8i.Map2KVs(data, func(i, j byte) bool { return i < j }, nil)
	for i := 0; i < nSmooth; i++ {
		vs = smooth(vs)
	}
	isBottom := FacIsBottom(vs, halfstep)
	for i := 0; i < len(vs); i++ {
		if isBottom(i) {
			m[ks[i]] = vs[i]
		}
	}

	// adjust to min value
	mp := make(map[byte]int)
	for k, v := range m {
		if min, n := ti.MinIdx(vs[k-1 : k+2]...); min < v {
			mp[k-1+byte(n)] = min
		} else {
			mp[k] = v
		}
	}

	if n < 0 {
		return mp
	}

	ks, vs = tu8i.Map2KVs(mp, nil, func(i, j int) bool { return i < j })
	mp = make(map[byte]int)
	for i := 0; i < n && i < len(ks); i++ {
		mp[ks[i]] = vs[i]
	}
	return mp
}

// return positions of max-up to max-down
func slope(data []int, step, nSmooth int) (sp []int) {

	for i := 0; i < nSmooth; i++ {
		data = smooth(data)
	}

	slope := []struct {
		ix int
		dy int
	}{}

	for i := step - 1; i < len(data); i++ {
		a := data[i-(step-1)]
		b := data[i]
		slope = append(slope, struct {
			ix int
			dy int
		}{
			ix: i - step/2,
			dy: b - a,
		})
	}

	sort.SliceStable(slope, func(i, j int) bool {
		return slope[i].dy > slope[j].dy
	})

	for _, s := range slope {
		sp = append(sp, s.ix)
	}
	return
}

func AveByte(data []byte) float64 {
	sum := 0.0
	for _, v := range data {
		sum += float64(v)
	}
	return sum / float64(len(data))
}

func StdDevByte(data []byte) float64 {
	sum2 := 0.0
	ave := AveByte(data)
	for _, v := range data {
		d2 := (float64(v) - ave) * (float64(v) - ave)
		sum2 += d2
	}
	return math.Sqrt(sum2 / float64(len(data)))
}

func AveFloat(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func StdDevFloat(data []float64) float64 {
	sum2 := 0.0
	ave := AveFloat(data)
	for _, v := range data {
		d2 := (v - ave) * (v - ave)
		sum2 += d2
	}
	return math.Sqrt(sum2 / float64(len(data)))
}
