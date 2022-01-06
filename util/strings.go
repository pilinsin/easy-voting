package util

import ()

func strPadding(s, pad string, nPad int) string {
	for i := 0; i < nPad; i++ {
		s += pad
	}
	return s
}

func strSplitEvenly(s string, size int) []string {
	n := 1 + len(s)/size
	nPad := n*size - len(s)
	s = strPadding(s, " ", nPad)

	strs := make([]string, n)
	for idx, _ := range strs {
		strs[idx] = s[idx*size : (idx+1)*size]
	}
	return strs
}
