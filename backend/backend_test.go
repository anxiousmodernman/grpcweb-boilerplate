package backend

import (
	"fmt"
	"testing"
)

func TestBackendData(t *testing.T) {

	assert := func(l []string, r ...string) {
		for i, x := range l {
			fmt.Println(x, r[i])
			if x != r[i] {
				t.Errorf("expected %s got %s", x, r[i])
			}
		}
	}

	a := []string{"10.2.1.20", "10.2.1.21"}
	b := []string{"10.2.1.56", "10.2.1.21"}

	c := combine(a, b)
	fmt.Println(c)

	assert(c, "10.2.1.2x", "10.2.1.21", "10.2.1.56")
}
