package inspector

import "testing"

func TestHeatmap(t *testing.T) {
	coordinator := NewCoordinator("angular/angular")

	coordinator.Heatmap()
}
