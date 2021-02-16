// +build !solution

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func check(r *MyTime) bool {
	if r.H < 0 || r.H > 23 {
		fmt.Println("Bad Check H")
		return false
	}

	if r.M < 0 || r.M > 59 {
		fmt.Println("Bad Check M")
		return false
	}

	if r.S < 0 || r.S > 59 {
		fmt.Println("Bad Check S")
		return false
	}
	return true
}

func ParseTime(tm string) (*MyTime, bool) {
	res := &MyTime{}

	if tm == "" {
		now := time.Now()
		res.H = now.Hour()
		res.M = now.Minute()
		res.S = now.Second()
		return res, true
	}

	digits := strings.Split(tm, ":")
	if len(digits) != 3 {
		return nil, false
	}

	h, ok := strconv.Atoi(digits[0])
	if ok != nil || len(digits[0]) != 2 {
		return nil, false
	}
	m, ok := strconv.Atoi(digits[1])
	if ok != nil || len(digits[1]) != 2 {
		return nil, false
	}
	s, ok := strconv.Atoi(digits[2])
	if ok != nil || len(digits[2]) != 2 {
		return nil, false
	}

	res = &MyTime{
		H: h,
		M: m,
		S: s,
	}

	if check(res) {
		return res, true
	} else {
		return nil, false
	}
}

func ParseK(k string) (int, bool) {
	if k == "" {
		return 1, true
	}
	res, err := strconv.Atoi(k)
	if err != nil || res < 1 || res > 30 {
		return res, false
	}
	return res, true
}

func GetNum(x int) string {
	if x == 0 {
		return Zero
	}
	if x == 1 {
		return One
	}
	if x == 2 {
		return Two
	}
	if x == 3 {
		return Three
	}
	if x == 4 {
		return Four
	}
	if x == 5 {
		return Five
	}
	if x == 6 {
		return Six
	}
	if x == 7 {
		return Seven
	}
	if x == 8 {
		return Eight
	}
	if x == 9 {
		return Nine
	}
	return "SUKAAAA"
}

func DrawNumber(p *Picture, num int, lastInd *int) {
	first := num / 10
	second := num % 10
	oneDigit := GetNum(first)
	twoDigit := GetNum(second)

	i := 0
	j := 0
	for _, v := range oneDigit {
		if v == '\n' {
			i++
			j = 0
			continue
		}

		if v == '.' {
			p.Matrix[i][*lastInd+j] = White
		} else {
			p.Matrix[i][*lastInd+j] = Cyan
		}
		j++
	}
	*lastInd += 8

	i, j = 0, 0
	for _, v := range twoDigit {
		if v == '\n' {
			i++
			j = 0
			continue
		}

		if v == '.' {
			p.Matrix[i][*lastInd+j] = White
		} else {
			p.Matrix[i][*lastInd+j] = Cyan
		}
		j++
	}
	*lastInd += 8
}

func DrawColon(p *Picture, lastInd *int) {
	i, j := 0, 0
	for _, v := range Colon {
		if v == '\n' {
			i++
			j = 0
			continue
		}

		if v == '.' {
			p.Matrix[i][*lastInd+j] = White
		} else {
			p.Matrix[i][*lastInd+j] = Cyan
		}
		j++
	}
	*lastInd += 4
}

func GetPicture(time *MyTime, k int) *Picture {
	n, m := (8*2 + 4 + 8*2 + 4 + 8*2), 12
	res := NewPicture(m, n)

	i := 0
	DrawNumber(res, time.H, &i)
	DrawColon(res, &i)
	DrawNumber(res, time.M, &i)
	DrawColon(res, &i)
	DrawNumber(res, time.S, &i)

	return res
}

func handleDigitalClock(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	structBody := &Req{
		Time: "r.FormValue()",
		K:    "r.FormValue()",
	}

	arTime, ok := r.Form["time"]
	if !ok || len(arTime) == 0 {
		structBody.Time = ""
	} else {
		if arTime[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		structBody.Time = arTime[0]
	}
	arK, ok := r.Form["k"]
	if !ok || len(arK) == 0 {
		structBody.K = ""
	} else {
		if arK[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		structBody.K = arK[0]
	}

	time, ok := ParseTime(structBody.Time)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	k, ok := ParseK(structBody.K)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tmpImg := GetPicture(time, k)
	finalN := k * (8*2 + 4 + 8*2 + 4 + 8*2)
	finalM := k * 12
	img := image.NewRGBA(image.Rect(0, 0, finalN, finalM))
	for i := 0; i < finalN; i++ {
		for j := 0; j < finalM; j++ {
			img.Set(i, j, tmpImg.Matrix[j/k][i/k])
		}
	}

	f, _ := os.Create("/tmp/img.png")
	png.Encode(f, img)

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func main() {
	if len(os.Args) != 3 {
		err := fmt.Errorf("Usage: ./m --port. Need two args you send: %d", len(os.Args))
		if err != nil {
			panic(err)
		}
		return
	}

	if os.Args[1] != "-port" {
		err := fmt.Errorf("Usage: ./m -port. Need port arg you send: -->  %s", os.Args[1])
		if err != nil {
			panic(err)
		}
	}
	http.HandleFunc("/", handleDigitalClock)
	if err := http.ListenAndServe(":"+os.Args[2], nil); err != nil {
		panic(err)
	}
}
