package config

import "testing"
import "github.com/stretchr/testify/require"

type A struct {
	C string
}

type Config struct {
	A1 A
	A2 A
	A3 []string
	B  int
}

func TestConfig(t *testing.T) {
	cfg := []byte(`
a1:
 c: test1
a2:
 c: test2
a3:
 - s1
 - s2
b: 123`)

	p := &Config{}
	err := LoadFromBytes(cfg, &p)
	require.NoError(t, err)

	require.Equal(t, &Config{
		A1: A{C: "test1"},
		A2: A{C: "test2"},
		A3: []string{"s1", "s2"},
		B:  123,
	}, p)

}
