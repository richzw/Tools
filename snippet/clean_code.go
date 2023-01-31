package snippet

// https://mp.weixin.qq.com/s/rhe1hFgIfFDHZldGhnzfJQ
func getPercentageRoundsSubstring(percentage float64) string {
	symbols := "★★★★★★★★★★☆☆☆☆☆☆☆☆☆☆"
	offset := 10 - int(percentage*10.0)
	return symbols[offset*3 : (offset+10)*3]
}
