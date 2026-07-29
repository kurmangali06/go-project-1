package kata

import (
	"errors"
	"testing"
)

func TestParseSKU(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    int64
		wantErr error // nil = ошибки быть не должно
	}{
		{name: "корректный sku", in: "2008", want: 2008, wantErr: nil},
		{name: "пустая строка", in: "", wantErr: ErrInvalidSKU},
		{name: "не число", in: "SKU123", wantErr: ErrInvalidSKU},
		{name: "число с хвостом", in: "12abc", wantErr: ErrInvalidSKU},
		{name: "дробное", in: "1.5", wantErr: ErrInvalidSKU},
		{name: "ноль", in: "0", wantErr: ErrInvalidSKU},
		{name: "отрицательное", in: "-1", wantErr: ErrInvalidSKU},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSKU(tt.in)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ошибка = %v, ожидали %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got != tt.want {
				t.Errorf("sku = %d, ожидали %d", got, tt.want)
			}
		})
	}

}
