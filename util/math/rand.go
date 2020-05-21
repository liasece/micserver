package math

import (
	crand "crypto/rand"
	"encoding/hex"
	"math/rand"
	"time"
)

// RandBetween rand a num between min to max
func RandBetween(min, max int) int {
	if min >= max || max == 0 {
		return max
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	random := r.Intn(max-min) + min
	return random
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// GenerateRandomBytes case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := crand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// GenerateRandomString case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return hex.EncodeToString(b), err
}

// SliceOutOfOrder 乱序切片
func SliceOutOfOrder(in []string) []string {
	rr := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(in)
	for i := l - 1; i > 0; i-- {
		r := rr.Intn(i)
		in[r], in[i] = in[i], in[r]
	}
	return in
}

// SliceOutOfOrderByInt 乱序切片
func SliceOutOfOrderByInt(in []uint64) []uint64 {
	rr := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(in)
	for i := l - 1; i > 0; i-- {
		r := rr.Intn(i)
		in[r], in[i] = in[i], in[r]
	}
	return in
}

type valueWeightItem struct {
	weight uint32
	value  uint64
}

// ValueWeightPair 权值对，根据权重随机一个值出来
type ValueWeightPair struct {
	allweight uint32
	valuelist []valueWeightItem
}

// NewValueWeightPair new value weight pair
func NewValueWeightPair() *ValueWeightPair {
	vwp := new(ValueWeightPair)
	vwp.valuelist = make([]valueWeightItem, 0)
	return vwp
}

// Add add weight
func (v *ValueWeightPair) Add(weight uint32, value uint64) {
	valueinfo := valueWeightItem{}
	valueinfo.weight = weight
	valueinfo.value = value
	v.valuelist = append(v.valuelist, valueinfo)
	v.allweight += weight
}

// Random random a num
func (v *ValueWeightPair) Random() uint64 {
	if v.allweight > 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randvalue := uint32(r.Intn(int(v.allweight)))
		addweight := uint32(0)
		for i := 0; i < len(v.valuelist); i++ {
			addweight += v.valuelist[i].weight
			if randvalue <= addweight {
				return v.valuelist[i].value
			}
		}
	}
	return 0
}

// RandWeight 根据权重列表随机出一个结果，返回命中下标
func RandWeight(weight []uint32) uint32 {
	total := uint32(0)
	for _, v := range weight {
		total += v
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randvalue := uint32(r.Intn(int(total)))
	tmp := uint32(0)
	for i, v := range weight {
		tmp += v
		if tmp >= randvalue {
			return uint32(i)
		}
	}
	return 0
}

// UtilWeightInterface infterface
type UtilWeightInterface interface {
	GetWeight() uint32
}

// RandWeightStruct 根据权重列表随机出一个结果，返回命中下标
func RandWeightStruct(weight []UtilWeightInterface) interface{} {
	if len(weight) == 0 {
		return nil
	}
	total := uint32(0)
	for _, v := range weight {
		total += v.GetWeight()
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randvalue := uint32(r.Intn(int(total)))
	tmp := uint32(0)
	for i, v := range weight {
		tmp += v.GetWeight()
		if tmp <= randvalue {
			return weight[i]
		}
	}
	return weight[0]
}
