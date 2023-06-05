package snippet

//For min(ctz(x), ctz(y)), we can use ctz(x | y) to gain better performance. But what about max(ctz(x), ctz(y))?
/*
x + y = min(x, y) + max(x, y)
and thus
max(ctz(x), ctz(y)) = ctz(x) + ctz(y) - min(ctz(x), ctz(y))

max(ctz(a),ctz(b))
ctz((a|-a)&(b|-b))
ctz(a)+ctz(b)-ctz(a&b)
*/

import "math/bits"

func maxr_zero(x, y uint64) int32 {
	loxs := ^x & (x - 1) // low zeros of x
	loys := ^y & (y - 1) // low zeros of y
	return int32(bits.TrailingZeros64((loxs | loys) + 1))
}
