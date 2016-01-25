package inspector
import (
	"testing"
)

func TestMinMax(t *testing.T) {

	tests := []struct {
		start, end, intersect int
	}{
		{
			start: 21,
			end: 34,
			intersect: 0,
		},
		{
			start: 2,
			end: 9,
			intersect: 0,
		},
		{
			start: 10,
			end: 20,
			intersect: 10,
		},
		{
			start: 0,
			end: 5,
			intersect: 0,
		},
		{
			start: 25,
			end: 35,
			intersect: 0,
		},
		{
			start: 12,
			end: 13,
			intersect: 2,
		},
		{
			start: 5,
			end: 25,
			intersect: 10,
		},
	}

	unit := &Unit{
		LineStart: 10,
		LineEnd: 20,
	}

	for _, test := range tests {
		got := unit.numberOfLinesIntersected(test.start, test.end)
		if got != test.intersect {
			t.Errorf("Wanted %d, got %d", test.intersect, got)
		}
	}

}

func TestHarvestINFILTRATOR(t *testing.T) {
	Harvest("docker/docker")
}