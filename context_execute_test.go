package template

import (
	"fmt"
	"github.com/ImSingee/dt"
	"github.com/ImSingee/mock/random"
	"github.com/ImSingee/tt"
	"strings"
	"testing"
	"unicode/utf8"
)

func init() {
	random.MustRegister(DefaultFunctionSet.MustRegisterFunction)
	fmt.Println(DefaultFunctionSet)
}

func TestExecute(t *testing.T) {
	mocks := map[string]string{
		"ABCD":         "ABCD",
		"@some":        "@some",
		"@@":           "@",
		"@@number()":   "@number()",
		"@@@number(5)": "@5",
	}

	for a, b := range mocks {
		t.Run("test for "+a, func(t *testing.T) {
			mock, err := Execute(a)
			tt.AssertIsNil(t, err)
			fmt.Println(mock)
			tt.AssertEqual(t, b, mock)
		})
	}

	//fmt.Println(Execute(`@string("aeiou", 3, 5)`))
}

func TestExecuteBoolean(t *testing.T) {
	trueCount := 0
	falseCount := 0
	for i := 0; i < 100; i++ {
		mock, err := Execute("@bool()")
		tt.AssertIsNil(t, err)
		fmt.Println(mock)
		switch mock {
		case "true":
			trueCount += 1
		case "false":
			falseCount += 1
		default:
			t.FailNow()
		}
	}
	tt.AssertNotEqual(t, 0, trueCount)
	tt.AssertNotEqual(t, 0, falseCount)
}

func TestExecuteInteger(t *testing.T) {
	mockName := []string{"int", "integer"}
	for _, name := range mockName {
		t.Run("mock use @"+name, func(t *testing.T) {
			t.Run("no param", func(t *testing.T) {
				for i := 0; i < 100; i++ {
					mock, err := Execute("@" + name + "()")
					tt.AssertIsNil(t, err)
					num, ok := dt.NumberFromString(mock)
					tt.AssertTrue(t, ok)
					tt.AssertTrue(t, num.IsInt64())
				}
			})

			t.Run("one param (min)", func(t *testing.T) {
				for i := 0; i < 100; i++ {
					r := int64(random.Integer(100, 1000))
					mock, err := Execute(fmt.Sprintf("@%s(%d)", name, r))
					tt.AssertIsNil(t, err)
					num, ok := dt.NumberFromString(mock)
					tt.AssertTrue(t, ok)
					tt.AssertTrue(t, num.IsInt64())
					tt.AssertTrue(t, num.Int64() >= r)
				}
			})

			t.Run("two param (min, max)", func(t *testing.T) {
				for i := 0; i < 100; i++ {
					r := int64(random.Integer(100, 1000))
					q := int64(random.Integer(10000, 100000))
					mock, err := Execute(fmt.Sprintf("@%s(%d, %d)", name, r, q))
					tt.AssertIsNil(t, err)
					num, ok := dt.NumberFromString(mock)
					tt.AssertTrue(t, ok)
					tt.AssertTrue(t, num.IsInt64())
					tt.AssertTrue(t, num.Int64() >= r)
					tt.AssertTrue(t, num.Int64() <= q)
				}
			})
		})
	}
}

func TestExecuteFloat(t *testing.T) {
	t.Run("no param", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			mock, err := Execute("@float()")
			tt.AssertIsNil(t, err)
			fmt.Println(mock)
			num, ok := dt.NumberFromString(mock)
			tt.AssertTrue(t, ok)
			tt.AssertTrue(t, num.Float())
		}
	})

	t.Run("one param (min)", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			r := random.Float64(100, 1000)
			mock, err := Execute(fmt.Sprintf("@float(%.3f)", r))
			fmt.Println(mock)
			tt.AssertIsNil(t, err)
			num, ok := dt.NumberFromString(mock)
			tt.AssertTrue(t, ok)
			tt.AssertTrue(t, num.Float())
			tt.AssertTrue(t, num.Float64() >= float64(int64(r)))
		}
	})

	t.Run("two param (min, max)", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			r := random.Float64(100, 1000)
			q := random.Float64(10000, 100000)
			mock, err := Execute(fmt.Sprintf("@float(%.3f, %.3f)", r, q))
			fmt.Println(mock)
			tt.AssertIsNil(t, err)
			num, ok := dt.NumberFromString(mock)
			tt.AssertTrue(t, ok)
			tt.AssertTrue(t, num.Float())
			tt.AssertTrue(t, num.Float64() >= float64(int64(r)))
			tt.AssertTrue(t, num.Float64() <= float64(int64(q)+1))
		}
	})

	t.Run("two param (min, max, d)", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			r := random.Float64(100, 1000)
			q := random.Float64(10000, 100000)
			d := random.Integer(1, 5)
			mock, err := Execute(fmt.Sprintf("@float(%.3f, %.3f, %d)", r, q, d))
			fmt.Println(mock)
			tt.AssertIsNil(t, err)
			num, ok := dt.NumberFromString(mock)
			tt.AssertTrue(t, ok)
			tt.AssertTrue(t, num.Float())
			tt.AssertTrue(t, num.Float64() >= float64(int64(r)))
			tt.AssertTrue(t, num.Float64() <= float64(int64(q)+1))
			tt.AssertTrue(t, strings.Contains(mock, "."))
			tt.AssertEqual(t, d, len(mock)-strings.Index(mock, ".")-1)
		}
	})
}

func TestExecuteChar(t *testing.T) {
	checkIn := func(t *testing.T, c string, pool string) {
		tt.AssertEqual(t, 1, utf8.RuneCountInString(c))
		tt.AssertTrue(t, strings.ContainsAny(c, pool))
	}

	t.Run("no param", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			mock, err := Execute("@char()")
			tt.AssertIsNil(t, err)
			fmt.Println(mock)
			checkIn(t, mock, string(random.MapCharacterPool("numletter")))
		}
	})

	t.Run("one param (pool)", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			mock, err := Execute(`@char("aeiou")`)
			tt.AssertIsNil(t, err)
			fmt.Println(mock)
			checkIn(t, mock, "aeiou")
		}
	})
}

func TestExecuteString(t *testing.T) {
	checkIn := func(t *testing.T, s string, pool string) {
		for _, c := range s {
			tt.AssertTrue(t, strings.ContainsAny(string(c), pool))
		}
	}

	t.Run("@string()", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			mock, err := Execute(`@string()`)
			tt.AssertIsNil(t, err)
			//fmt.Println(mock)
			checkIn(t, mock, string(random.MapCharacterPool("numletter")))
		}
	})

	t.Run("@string(length)", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			length := random.Natural(10, 100)
			mock, err := Execute(fmt.Sprintf(`@string(%d)`, length))
			tt.AssertIsNil(t, err)
			//fmt.Println(mock)
			checkIn(t, mock, string(random.MapCharacterPool("numletter")))
			tt.AssertEqual(t, length, len(mock))
		}
	})

	t.Run("@string(min, max)", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			min := random.Natural(10, 50)
			max := random.Natural(50, 100)

			mock, err := Execute(fmt.Sprintf(`@string(%d, %d)`, min, max))
			tt.AssertIsNil(t, err)
			//fmt.Println(mock)
			checkIn(t, mock, string(random.MapCharacterPool("numletter")))
			tt.AssertTrue(t, len(mock) >= min)
			tt.AssertTrue(t, len(mock) <= max)
		}
	})

	t.Run("@string(pool)", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			mock, err := Execute(fmt.Sprintf(`@string("aeiou")`))
			tt.AssertIsNil(t, err)
			//fmt.Println(mock)
			checkIn(t, mock, "aeiou")
		}
	})

	t.Run("@string(pool, length)", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			length := random.Natural(10, 100)

			mock, err := Execute(fmt.Sprintf(`@string("aeiou", %d)`, length))
			tt.AssertIsNil(t, err)
			//fmt.Println(mock)
			checkIn(t, mock, "aeiou")
			tt.AssertEqual(t, length, len(mock))
		}
	})

	t.Run("@string(pool, min, max)", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			min := random.Natural(10, 50)
			max := random.Natural(50, 100)

			mock, err := Execute(fmt.Sprintf(`@string("aeiou", %d, %d)`, min, max))
			tt.AssertIsNil(t, err)
			//fmt.Println(mock)
			checkIn(t, mock, "aeiou")
			tt.AssertTrue(t, len(mock) >= min)
			tt.AssertTrue(t, len(mock) <= max)
		}
	})
}

func TestExecuteSentence(t *testing.T) {
	t.Run("@sentence()", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			mock, err := Execute(`@sentence()`)
			tt.AssertIsNil(t, err)
			fmt.Println(mock)
			tt.AssertEqual(t, strings.TrimSpace(mock), mock)
			tt.AssertTrue(t,
				strings.HasSuffix(mock, ".") ||
					strings.HasSuffix(mock, "?") ||
					strings.HasSuffix(mock, "!"),
			)
		}
	})
}

func TestExecuteIncrement(t *testing.T) {
	random.IncrementReset()
	t.Run("@increment()", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			mock, err := Execute(`@increment()`)
			tt.AssertIsNil(t, err)
			//fmt.Println(mock)
			tt.AssertEqual(t, dt.ToString(i+1), mock)
		}
	})

	random.IncrementReset()
	t.Run("@increment(step)", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			mock, err := Execute(`@increment(3)`)
			tt.AssertIsNil(t, err)
			//fmt.Println(mock)
			tt.AssertEqual(t, dt.ToString(3*(i+1)), mock)
		}
	})

	random.IncrementReset()
	t.Run("@increment(step, delta)", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			mock, err := Execute(`@increment(3, 1024)`)
			tt.AssertIsNil(t, err)
			//fmt.Println(mock)
			tt.AssertEqual(t, dt.ToString(3*(i+1)+1024), mock)
		}
	})
}
