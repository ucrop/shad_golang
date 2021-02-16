package main

import "image/color"

// MyTime ...
type MyTime struct {
	H int `json:"h"`
	M int `json:"m"`
	S int `json:"s"`
}

// Req ...
type Req struct {
	Time string `json:"time"`
	K    string `json:"k"`
}

// Resp ...
type Resp struct {
	Time MyTime `json:"time"`
	K    int    `json:"k"`
}

// Picture ...
type Picture struct {
	Matrix [][]color.RGBA
}

// NewPicture ...
func NewPicture(h, w int) *Picture {
	tmpPicture := &Picture{
		Matrix: make([][]color.RGBA, h),
	}
	for i := 0; i < h; i++ {
		tmpPicture.Matrix[i] = make([]color.RGBA, w)
	}
	return tmpPicture
}
