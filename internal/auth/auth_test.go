package auth

import "testing"

func Test_decomposeBasic(t *testing.T) {

	tests := []struct {
		name    string
		header  string
		want    string
		want1   string
		wantErr bool
	}{
		{"empty", "", "", "", true},
		{"not basic", "Bad", "", "", true},
		{"Basic but empty", "Basic 34", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := decomposeBasic(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("decomposeBasic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("decomposeBasic() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("decomposeBasic() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
