package randomdata

import (
	"math/rand"
	"strconv"
	"time"
)

const BufferSize = 100000

var (
	randPool chan string
)

func init() {
	randPool = make(chan string, BufferSize)
	go run()
}

func run() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r = rand.New(rand.NewSource(time.Now().UnixNano()))
		default:
			if len(randPool) < BufferSize {
				randPool <- strconv.Itoa(r.Intn(900000) + 100000)
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func GenRandomSixNumberString() string {
	select {
	case randNumStr := <-randPool:
		return randNumStr
	case <-time.After(10 * time.Millisecond):
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		return strconv.Itoa(r.Intn(900000) + 100000)
	}
}
