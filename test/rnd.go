package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	var flt int
	//var rndstr string
	var rndbyte [20]byte
	flt = time.Now().Second() + time.Now().Minute()*60 + time.Now().Hour()*3600
	flt *= 10000000000
	flt += time.Now().Nanosecond()
	fmt.Printf("now %d\n", flt)
	rand.Seed(int64(flt))
	for i := 0; i < 20; i++ {
		rndnum := int(rand.Float32() * 1000000)
		rndnum %= 62
		fmt.Printf("rand now %d\n", rndnum)
		if rndnum < 26 {
			rndbyte[i] = byte(int('a') + rndnum)
		} else if rndnum < 52 {
			rndbyte[i] = byte(int('A') + rndnum - 26)
		} else {
			rndbyte[i] = byte(int('0') + rndnum - 52)
		}

	}
	//rndbyte[20] = byte(000)

	rndstr := string(rndbyte[:])
	fmt.Printf("rndbyte %v rndstr %s", rndbyte, rndstr)

}
