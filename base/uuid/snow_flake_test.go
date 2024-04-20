package uuid

import "testing"

func TestUUID(t *testing.T) {
	t.Log(UUID.UUID())

}

func TestGenerateDup(t *testing.T) {
	node := UUID

	var x, y int64
	for i := 0; i < 1000000; i++ {
		y = node.UUID()
		if x == y {
			t.Errorf("x(%d) & y(%d) are the same", x, y)
		}
		x = y
	}
}

// I feel like there's probably a better way
func TestRace(t *testing.T) {
	node := UUID

	go func() {
		for i := 0; i < 1000000000; i++ {
			NewSnowFlake()
		}
	}()

	for i := 0; i < 4000; i++ {
		node.UUID()
	}
}
