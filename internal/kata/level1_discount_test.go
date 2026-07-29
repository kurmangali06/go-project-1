package kata

import "testing"

func TestDiscount(t *testing.T) {
	// контракт в комментарии над Discount — тестируй по нему,
	// в тело функции пока не смотри
	got := Discount(100, 0)
	var want uint32 = 100
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	got = Discount(100, 100)
	want = 0
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	got = Discount(100, 110)
	want = 0
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	got = Discount(999, 10)
	want = 899
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
