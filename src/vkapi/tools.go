package vkapi

import (
	"math/rand"
	"strconv"
	"time"
)

//GetRandomID returns string with random_id for random_id field in some vkAPI calls
func GetRandomID() string {
	rand.Seed(time.Now().UnixNano())
	if time.Now().UnixNano()%2 == 0 {
		return strconv.Itoa(int(rand.Int31()))
	}
	return strconv.Itoa(int(rand.Int31() * -1))
}
