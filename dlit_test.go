package dlit

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

func TestNew(t *testing.T) {
	anError := errors.New("this is an error")
	cases := []struct {
		in   interface{}
		want *Literal
	}{
		{6, MustNew(6)},
		{6.0, MustNew(6.0)},
		{6.6, MustNew(6.6)},
		{float32(6.6), MustNew(float32(6.6))},
		{int64(922336854775807), MustNew(922336854775807)},
		{int64(9223372036854775807), MustNew(9223372036854775807)},
		{"98292223372036854775807", MustNew("98292223372036854775807")},
		{"6", MustNew("6")},
		{"6.6", MustNew("6.6")},
		{"abc", MustNew("abc")},
		{true, MustNew(true)},
		{false, MustNew(false)},
		{anError, MustNew(anError)},
	}

	for _, c := range cases {
		got, err := New(c.in)
		if err != nil {
			t.Errorf("New(%v) err: %v", c.in, err)
		}

		if err := checkLitsMatch(got, c.want); err != nil {
			t.Errorf("New(%v): %s", c.in, err)
		}
	}
}

func TestNew_errors(t *testing.T) {
	cases := []struct {
		in        interface{}
		want      *Literal
		wantError error
	}{
		{complex64(1), MustNew(InvalidKindError("complex64")),
			InvalidKindError("complex64")},
		{complex128(1), MustNew(InvalidKindError("complex128")),
			InvalidKindError("complex128")},
	}

	for _, c := range cases {
		got, err := New(c.in)
		if !errorMatch(err, c.wantError) {
			t.Errorf("New(%v) - err == %v, wantError == %v", c.in, err, c.wantError)
		}

		if err := checkLitsMatch(got, c.want); err != nil {
			t.Errorf("New(%v): %s", c.in, err)
		}
	}
}

func TestNewString(t *testing.T) {
	cases := []struct {
		in   string
		want *Literal
	}{
		{"", MustNew("")},
		{"6", MustNew("6")},
		{"6.27", MustNew("6.27")},
		{"Hello how are you today", MustNew("Hello how are you today")},
	}

	for _, c := range cases {
		got := NewString(c.in)
		if err := checkLitsMatch(got, c.want); err != nil {
			t.Errorf("NewString(%v): %s", c.in, err)
		}
	}
}

func TestMustNew(t *testing.T) {
	anError := errors.New("this is an error")
	cases := []struct {
		in   interface{}
		want *Literal
	}{
		{6, MustNew(6)},
		{6.0, MustNew(6.0)},
		{6.6, MustNew(6.6)},
		{float32(6.6), MustNew(float32(6.6))},
		{int64(922336854775807), MustNew(922336854775807)},
		{int64(9223372036854775807), MustNew(9223372036854775807)},
		{"98292223372036854775807", MustNew("98292223372036854775807")},
		{"6", MustNew("6")},
		{"6.6", MustNew("6.6")},
		{"abc", MustNew("abc")},
		{true, MustNew(true)},
		{false, MustNew(false)},
		{anError, MustNew(anError)},
	}

	for _, c := range cases {
		got := MustNew(c.in)
		if err := checkLitsMatch(got, c.want); err != nil {
			t.Errorf("MustNew(%v): %s", c.in, err)
		}
	}
}

func TestMustNew_panic(t *testing.T) {
	cases := []struct {
		in        interface{}
		wantPanic string
	}{
		{6, ""},
		{complex64(1), InvalidKindError("complex64").Error()},
	}

	for _, c := range cases {
		paniced := false
		defer func() {
			if r := recover(); r != nil {
				if r.(string) == c.wantPanic {
					paniced = true
				} else {
					t.Errorf("MustNew(%v) - got panic: %s, wanted: %s",
						c.in, r, c.wantPanic)
				}
			}
		}()
		MustNew(c.in)
		if c.wantPanic != "" && !paniced {
			t.Errorf("MustNew(%v) - failed to panic with: %s", c.in, c.wantPanic)
		}
	}
}

func TestInt(t *testing.T) {
	cases := []struct {
		in        *Literal
		want      int64
		wantIsInt bool
	}{
		{MustNew(6), 6, true},
		{MustNew(6.0), 6, true},
		{MustNew(float32(6.0)), 6, true},
		{MustNew("6"), 6, true},
		{MustNew("6.0"), 6, true},
		{MustNew("6."), 6, true},
		{MustNew("6.0000"), 6, true},
		{MustNew("-6"), -6, true},
		{MustNew("-6.0"), -6, true},
		{MustNew("-6."), -6, true},
		{MustNew(fmt.Sprintf("%d", int64(math.MinInt64))),
			int64(math.MinInt64), true},
		{MustNew(fmt.Sprintf("%d", int64(math.MaxInt64))),
			int64(math.MaxInt64), true},
		{MustNew(fmt.Sprintf("-1%d", int64(math.MinInt64))), 0, false},
		{MustNew(fmt.Sprintf("1%d", int64(math.MaxInt64))), 0, false},
		{MustNew("-9223372036854775809"), 0, false},
		{MustNew("9223372036854775808"), 0, false},
		{MustNew(6.6), 0, false},
		{MustNew("6.6"), 0, false},
		{MustNew("6.06"), 0, false},
		{MustNew("abc"), 0, false},
		{MustNew(true), 0, false},
		{MustNew(false), 0, false},
		{MustNew(".0"), 0, false},
		{MustNew(".23"), 0, false},
		{MustNew(errors.New("This is an error")), 0, false},
	}

	for _, c := range cases {
		got, gotIsInt := c.in.Int()
		if got != c.want || gotIsInt != c.wantIsInt {
			t.Errorf("Int() with Literal: %s - got: %d, %t - want: %d, %t",
				c.in, got, gotIsInt, c.want, c.wantIsInt)
		}
	}
}

func TestFloat(t *testing.T) {
	cases := []struct {
		in          *Literal
		want        float64
		wantIsFloat bool
	}{
		{MustNew(6), 6.0, true},
		{MustNew(int64(922336854775807)), 922336854775807.0, true},
		{MustNew(fmt.Sprintf("%G", float64(math.SmallestNonzeroFloat64))),
			math.SmallestNonzeroFloat64, true},
		{MustNew(fmt.Sprintf("%f", float64(math.MaxFloat64))),
			float64(math.MaxFloat64), true},
		{MustNew(6.0), 6.0, true},
		{MustNew("6"), 6.0, true},
		{MustNew(6.678934), 6.678934, true},
		{MustNew("6.678394"), 6.678394, true},
		{MustNew("abc"), 0, false},
		{MustNew(true), 0, false},
		{MustNew(false), 0, false},
		{MustNew(errors.New("This is an error")), 0, false},
	}

	for _, c := range cases {
		got, gotIsFloat := c.in.Float()
		if got != c.want || gotIsFloat != c.wantIsFloat {
			t.Errorf("Float() with Literal: %s - got: %f, %t - want: %f, %t",
				c.in, got, gotIsFloat, c.want, c.wantIsFloat)
		}
	}
}

func TestBool(t *testing.T) {
	cases := []struct {
		in         *Literal
		want       bool
		wantIsBool bool
	}{
		{MustNew(1), true, true},
		{MustNew(2), false, false},
		{MustNew(0), false, true},
		{MustNew(1.0), true, true},
		{MustNew(2.0), false, false},
		{MustNew(2.25), false, false},
		{MustNew(0.0), false, true},
		{MustNew(true), true, true},
		{MustNew(false), false, true},
		{MustNew("true"), true, true},
		{MustNew("false"), false, true},
		{MustNew("True"), true, true},
		{MustNew("False"), false, true},
		{MustNew("TRUE"), true, true},
		{MustNew("FALSE"), false, true},
		{MustNew("t"), true, true},
		{MustNew("f"), false, true},
		{MustNew("T"), true, true},
		{MustNew("F"), false, true},
		{MustNew("1"), true, true},
		{MustNew("0"), false, true},
		{MustNew("bob"), false, false},
		{MustNew(errors.New("This is an error")), false, false},
	}

	for _, c := range cases {
		got, gotIsBool := c.in.Bool()
		if got != c.want || gotIsBool != c.wantIsBool {
			t.Errorf("Bool() with Literal: %s - got: %t, %t - want: %t, %t",
				c.in, got, gotIsBool, c.want, c.wantIsBool)
		}
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		in   *Literal
		want string
	}{
		{MustNew(124), "124"},
		{MustNew(int64(922336854775807)), "922336854775807"},
		{MustNew(int64(9223372036854775807)), "9223372036854775807"},
		{MustNew("98292223372036854775807"), "98292223372036854775807"},
		{MustNew("Hello my name is fred"), "Hello my name is fred"},
		{MustNew(124.0), "124"},
		{MustNew(124.56728482274629), "124.56728482274629"},
		{MustNew(true), "true"},
		{MustNew(false), "false"},
		{MustNew(errors.New("This is an error")), "This is an error"},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("String() with Literal: %v - got: %v, want: %v",
				c.in, got, c.want)
		}
	}
}

func TestErr(t *testing.T) {
	cases := []struct {
		in   *Literal
		want error
	}{
		{MustNew(1), nil},
		{MustNew(2), nil},
		{MustNew("true"), nil},
		{MustNew(2.25), nil},
		{MustNew("hello"), nil},
		{MustNew(errors.New("This is an error")), errors.New("This is an error")},
	}

	for _, c := range cases {
		got := c.in.Err()
		if !errorMatch(c.want, got) {
			t.Errorf("Err() with Literal: %s - got: %s, want: %s", c.in, got, c.want)
		}
	}
}

func checkLitsMatch(got, want *Literal) error {
	if got.String() != want.String() {
		return fmt.Errorf("got.String(): %s, want.String(): %s", got, want)
	}

	if got.Err() != want.Err() {
		return fmt.Errorf("got.Err(): %s, want.Err(): %s", got.Err(), want.Err())
	}

	canBeIntGot, intGot := got.Int()
	canBeIntWant, intWant := want.Int()
	if canBeIntGot != canBeIntWant || intGot != intWant {
		return fmt.Errorf("Int statuses do not match")
	}

	canBeFloatGot, intGot := got.Float()
	canBeFloatWant, intWant := want.Float()
	if canBeFloatGot != canBeFloatWant || intGot != intWant {
		return fmt.Errorf("Float statuses do not match")
	}

	canBeBoolGot, intGot := got.Bool()
	canBeBoolWant, intWant := want.Bool()
	if canBeBoolGot != canBeBoolWant || intGot != intWant {
		return fmt.Errorf("Bool statuses do not match")
	}
	return nil
}

/*************************
       Benchmarks
*************************/
func BenchmarkInt_unknown(b *testing.B) {
	b.StopTimer()
	var sum int64
	for n := 0; n < b.N; n++ {
		l := MustNew("7.0")
		b.StartTimer()
		v, ok := l.Int()
		b.StopTimer()
		if !ok {
			b.Errorf("Int - ok: %t, want: %t", ok, true)
		}
		sum += v
	}
	if sum != int64(7*b.N) {
		b.Errorf("sum: %d, want: %d", sum, 7*b.N)
	}
}

func BenchmarkFloat_unknown(b *testing.B) {
	b.StopTimer()
	var sum float64
	for n := 0; n < b.N; n++ {
		l := MustNew("7.0")
		b.StartTimer()
		v, ok := l.Float()
		b.StopTimer()
		if !ok {
			b.Errorf("Float - ok: %t, want: %t", ok, true)
		}
		sum += v
	}
	if sum != float64(7.0*b.N) {
		b.Errorf("sum: %f, want: %f", sum, float64(7.0*b.N))
	}
}
func BenchmarkBool_unknown(b *testing.B) {
	b.StopTimer()
	var countTrue int
	for n := 0; n < b.N; n++ {
		l := MustNew("true")
		b.StartTimer()
		v, ok := l.Bool()
		b.StopTimer()
		if !ok {
			b.Errorf("Bool - ok: %t, want: %t", ok, true)
		}
		if v {
			countTrue++
		}
	}
	if countTrue != b.N {
		b.Errorf("countTrue: %d, want: %d", countTrue, b.N)
	}
}

func BenchmarkBool_multiple(b *testing.B) {
	b.StopTimer()
	var countTrue int
	for n := 0; n < b.N; n++ {
		lits := []*Literal{MustNew(true), MustNew(1), MustNew(1.0), MustNew("1")}
		for _, l := range lits {
			b.StartTimer()
			v, ok := l.Bool()
			b.StopTimer()
			if !ok {
				b.Errorf("Bool - ok: %t, want: %t", ok, true)
			}
			if v {
				countTrue++
			}
		}
	}
	if countTrue != 4*b.N {
		b.Errorf("countTrue: %d, want: %d", countTrue, 4*b.N)
	}
}

func BenchmarkNewInt(b *testing.B) {
	b.StopTimer()
	var sum int64
	for n := 0; n < b.N; n++ {
		b.StartTimer()
		l, _ := New(7)
		b.StopTimer()
		v, _ := l.Int()
		sum += v
	}
	if sum != int64(7*b.N) {
		b.Errorf("sum: %d, want: %d", sum, 7*b.N)
	}
}

func BenchmarkNewString(b *testing.B) {
	b.StopTimer()
	var sum int64
	for n := 0; n < b.N; n++ {
		b.StartTimer()
		l, _ := New("7.0")
		b.StopTimer()
		v, _ := l.Int()
		sum += v
	}
	if sum != int64(7*b.N) {
		b.Errorf("sum: %d, want: %d", sum, 7*b.N)
	}
}

/***********************
   Helper functions
************************/
func errorMatch(e1 error, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}
	if e1 == nil || e2 == nil {
		return false
	}
	if e1.Error() == e2.Error() {
		return true
	}
	return false
}
