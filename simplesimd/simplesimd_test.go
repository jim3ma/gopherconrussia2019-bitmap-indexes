package simplesimd

import (
	"math/rand"
	"testing"
	"time"
)

const restaurants = 65536
const bitmapLength = restaurants / 8 // 8192 bytes

func TestSIMDSimpleAnd(t *testing.T) {
	var (
		bitmapA   = make([]byte, bitmapLength)
		bitmapB   = make([]byte, bitmapLength)
		bitmapRes = make([]byte, bitmapLength)
	)

	for i := 0; i < len(bitmapRes); i++ {
		bitmapA[i] = 3 << 1
		bitmapB[i] = 3
	}

	andSIMD(bitmapA, bitmapB, bitmapRes)

	for i := 0; i < len(bitmapRes); i++ {
		if bitmapRes[i] != 2 {
			t.Fatalf("byte %d of result is %d (expected 2)", i, bitmapRes[i])
		}
	}
}

func TestSIMDSimpleOr(t *testing.T) {
	var (
		bitmapA   = make([]byte, bitmapLength)
		bitmapB   = make([]byte, bitmapLength)
		bitmapRes = make([]byte, bitmapLength)
	)

	for i := 0; i < len(bitmapRes); i++ {
		bitmapA[i] = 3 << 1
		bitmapB[i] = 3
	}

	orSIMD(bitmapA, bitmapB, bitmapRes)

	for i := 0; i < len(bitmapRes); i++ {
		if bitmapRes[i] != 7 {
			t.Fatalf("byte %d of result is %d (expected 7)", i, bitmapRes[i])
		}
	}
}

func TestSIMDSimpleAndNot(t *testing.T) {
	var (
		bitmapA   = make([]byte, bitmapLength)
		bitmapB   = make([]byte, bitmapLength)
		bitmapRes = make([]byte, bitmapLength)
	)

	for i := 0; i < len(bitmapRes); i++ {
		bitmapA[i] = 255
		bitmapB[i] = ^(byte(1) << 5)
	}

	andnotSIMD(bitmapA, bitmapB, bitmapRes)

	var expected byte = 1 << 5
	for i := 0; i < len(bitmapRes); i++ {
		if bitmapRes[i] != byte(expected) {
			t.Fatalf("byte %d of result is %d (expected %d)", i, bitmapRes[i], expected)
		}
	}
}

func TestSimpleBitmapIndex(t *testing.T) {
	var (
		nearMetro      = make([]byte, bitmapLength)
		privateParking = make([]byte, bitmapLength)
		terrace        = make([]byte, bitmapLength)
		reservations   = make([]byte, bitmapLength)
		veganFriendly  = make([]byte, bitmapLength)
		expensive      = make([]byte, bitmapLength)
	)

	rand.Seed(time.Now().UnixNano())

	fill(nearMetro, 0.1)
	fill(privateParking, 0.01)
	fill(terrace, 0.05)
	fill(reservations, 0.95)
	fill(veganFriendly, 0.2)
	fill(expensive, 0.1)

	resBitmap := make([]byte, bitmapLength)

	not(expensive, resBitmap)
	and(terrace, resBitmap, resBitmap)
	and(reservations, resBitmap, resBitmap)

	resRestaurants := indexes(resBitmap)

	t.Log(len(resRestaurants))
}

func TestSIMDBitmapIndex(t *testing.T) {
	var (
		nearMetro      = make([]byte, bitmapLength)
		privateParking = make([]byte, bitmapLength)
		terrace        = make([]byte, bitmapLength)
		reservations   = make([]byte, bitmapLength)
		veganFriendly  = make([]byte, bitmapLength)
		expensive      = make([]byte, bitmapLength)
	)

	rand.Seed(time.Now().UnixNano())

	fill(nearMetro, 0.1)
	fill(privateParking, 0.01)
	fill(terrace, 0.05)
	fill(reservations, 0.95)
	fill(veganFriendly, 0.2)
	fill(expensive, 0.1)

	resBitmap := make([]byte, bitmapLength)

	andnotSIMD(terrace, expensive, resBitmap)
	andSIMD(reservations, resBitmap, resBitmap)

	resRestaurants := indexes(resBitmap)

	t.Log(len(resRestaurants))
}

func BenchmarkSimpleBitmapIndex(b *testing.B) {
	var (
		nearMetro      = make([]byte, bitmapLength)
		privateParking = make([]byte, bitmapLength)
		terrace        = make([]byte, bitmapLength)
		reservations   = make([]byte, bitmapLength)
		veganFriendly  = make([]byte, bitmapLength)
		expensive      = make([]byte, bitmapLength)
	)

	rand.Seed(time.Now().UnixNano())

	fill(nearMetro, 0.1)
	fill(privateParking, 0.01)
	fill(terrace, 0.05)
	fill(reservations, 0.95)
	fill(veganFriendly, 0.2)
	fill(expensive, 0.1)

	resBitmap := make([]byte, bitmapLength)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		not(expensive, resBitmap)
		and(terrace, resBitmap, resBitmap)
		and(reservations, resBitmap, resBitmap)
	}
}

func BenchmarkSimpleSIMDBitmapIndex(b *testing.B) {
	var (
		nearMetro      = make([]byte, bitmapLength)
		privateParking = make([]byte, bitmapLength)
		terrace        = make([]byte, bitmapLength)
		reservations   = make([]byte, bitmapLength)
		veganFriendly  = make([]byte, bitmapLength)
		expensive      = make([]byte, bitmapLength)
	)

	rand.Seed(time.Now().UnixNano())

	fill(nearMetro, 0.1)
	fill(privateParking, 0.01)
	fill(terrace, 0.05)
	fill(reservations, 0.95)
	fill(veganFriendly, 0.2)
	fill(expensive, 0.1)

	resBitmap := make([]byte, bitmapLength)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		andnotSIMD(terrace, expensive, resBitmap)
		andSIMD(reservations, resBitmap, resBitmap)
	}
}
