package snippet

import "sort"

func subSets(arr []int) [][]int {
	res, track := [][]int{}, []int{}
	backtrack(arr, 0, track, &res)
	return res
}

func backtrack(arr []int, start int, track []int, res *[][]int) {
	if len(track) > 0 {
		tmp := make([]int, len(track))
		copy(tmp, track)
		*res = append(*res, tmp)
	}

	for i := start; i < len(arr); i++ {
		track = append(track, arr[i])

		backtrack(arr, i+1, track, res)

		track = track[:(len(track) - 1)]
	}
}

func subSetsDup(arr []int) [][]int {
	res, track := [][]int{}, []int{}
	sort.Ints(arr)
	backtrackDup(arr, 0, track, &res)
	return res
}

func backtrackDup(arr []int, start int, track []int, res *[][]int) {
	if len(track) > 0 {
		tmp := make([]int, len(track))
		copy(tmp, track)
		*res = append(*res, tmp)
	}

	for i := start; i < len(arr); i++ {
		if i > start && arr[i] == arr[i-1] {
			continue
		}
		track = append(track, arr[i])

		backtrackDup(arr, i+1, track, res)

		track = track[:(len(track) - 1)]
	}
}

func combine(arr []int, n int) [][]int {
	res, track := [][]int{}, []int{}
	backtrackLimit(arr, 0, n, track, &res)
	return res
}

func backtrackLimit(arr []int, start, limit int, track []int, res *[][]int) {
	if len(track) == limit {
		tmp := make([]int, len(track))
		copy(tmp, track)
		*res = append(*res, tmp)
		return
	}

	for i := start; i < len(arr); i++ {
		track = append(track, arr[i])

		backtrackLimit(arr, i+1, limit, track, res)

		track = track[:(len(track) - 1)]
	}
}

func permute(arr []int) [][]int {
	res, track := [][]int{}, []int{}
	visited := make([]int, len(arr))
	backtrackPermute(arr, track, visited, &res)
	return res
}

func backtrackPermute(arr, track, visited []int, res *[][]int) {
	if len(track) == len(arr) {
		tmp := make([]int, len(track))
		copy(tmp, track)
		*res = append(*res, tmp)
		return
	}

	for i := 0; i < len(arr); i++ {
		if visited[i] == 1 {
			continue
		}

		track = append(track, arr[i])
		visited[i] = 1

		backtrackPermute(arr, track, visited, res)

		track = track[:len(track)-1]
		visited[i] = 0
	}
}
