package pow

import "math/bits"

func leadingZerosCount(data []byte) byte {
	count := 0
	for _, v := range data {
		if v == 0 {
			count += 8
		} else {
			count += bits.LeadingZeros8(v)
			break
		}
	}
	if count > 255 {
		return 255
	}
	return byte(count)
}
