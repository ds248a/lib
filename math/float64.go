package math

import (
	"bytes"
	"math"
	"strconv"

	"github.com/ds248a/lib/bpool"
)

type Float64 struct{}

// NewFloat64 для операций без потери точности в десятичных дробях.
func NewFloat64() *Float64 {
	return &Float64{}
}

// DecimalPlaces возвращает количество десятичных знаков.
func (fc *Float64) DecimalPlaces(x float64) int {
	buf := bpool.Get()
	defer bpool.Put(buf)

	buf.B = strconv.AppendFloat(buf.B, x, 'f', -1, 64)
	i := bytes.IndexByte(buf.B, '.')

	if i > -1 {
		return len(buf.B) - i - 1
	}

	return 0
}

// Max возвращает большее из двух значений.
func (fc *Float64) Max(x, y int) int {
	if x > y {
		return x
	}

	return y
}

// Add возвращает a + b.
func (fc *Float64) Add(a, b float64) float64 {
	exp := math.Pow10(fc.Max(fc.DecimalPlaces(a), fc.DecimalPlaces(b)))

	intA := math.Round(a * exp)
	intB := math.Round(b * exp)

	return (intA + intB) / exp
}

// Sub возвращает a - b.
func (fc *Float64) Sub(a, b float64) float64 {
	exp := math.Pow10(fc.Max(fc.DecimalPlaces(a), fc.DecimalPlaces(b)))

	intA := math.Round(a * exp)
	intB := math.Round(b * exp)

	return (intA - intB) / exp
}

// Mul возвращает a * b.
func (fc *Float64) Mul(a, b float64) float64 {
	placesA := fc.DecimalPlaces(a)
	placesB := fc.DecimalPlaces(b)

	expA := math.Pow10(placesA)
	expB := math.Pow10(placesB)

	intA := math.Round(a * expA)
	intB := math.Round(b * expB)

	exp := math.Pow10(placesA + placesB)

	return (intA * intB) / exp
}

// Div возвращает a / b.
func (fc *Float64) Div(a, b float64) float64 {
	placesA := fc.DecimalPlaces(a)
	placesB := fc.DecimalPlaces(b)

	expA := math.Pow10(placesA)
	expB := math.Pow10(placesB)

	intA := math.Round(a * expA)
	intB := math.Round(b * expB)

	exp := math.Pow10(placesA - placesB)

	return (intA / intB) / exp
}

// Mod возвращает a % b.
func (fc *Float64) Mod(a, b float64) float64 {
	quo := math.Round(fc.Div(a, b))

	return fc.Sub(a, fc.Mul(b, quo))
}
