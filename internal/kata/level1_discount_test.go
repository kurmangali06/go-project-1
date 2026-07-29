package kata

import "testing"

func TestDiscount(t *testing.T) {
	tests := []struct {
		name    string
		price   uint32
		percent uint8
		want    uint32
	}{
		{name: "без скидки", price: 100, percent: 0, want: 100},
		{name: "скидка 100%", price: 100, percent: 100, want: 0},
		{name: "скидка больше 100%", price: 100, percent: 110, want: 0},
		{name: "округление вниз", price: 999, percent: 10, want: 899},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Discount(tt.price, tt.percent); got != tt.want {
				t.Errorf("Discount(%d, %d) = %d, ожидали %d", tt.price, tt.percent, got, tt.want)
			}
		})
	}
}
