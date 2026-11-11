package pagination

import "testing"

func TestOffset(t *testing.T) {
	tests := []struct {
		name   string
		params Params
		want   int
	}{
		{name: "first page starts at 0", params: Params{Page: 1, Limit: 10}, want: 0},
		{name: "second page", params: Params{Page: 2, Limit: 10}, want: 10},
		{name: "third page with limit 10", params: Params{Page: 3, Limit: 10}, want: 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.Offset(); got != tt.want {
				t.Errorf("Offset() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		page    string
		limit   string
		want    Params
		wantErr bool
	}{
		{name: "defaults when empty", want: Params{Page: DefaultPage, Limit: DefaultLimit}},
		{name: "valid values", page: "3", limit: "25", want: Params{Page: 3, Limit: 25}},
		{name: "limit at the maximum", limit: "100", want: Params{Page: DefaultPage, Limit: 100}},
		{name: "page is not a number", page: "abc", wantErr: true},
		{name: "page is zero", page: "0", wantErr: true},
		{name: "limit is zero", limit: "0", wantErr: true},
		{name: "limit above the maximum", limit: "101", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.page, tt.limit)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}