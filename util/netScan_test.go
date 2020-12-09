package util

import (
	"net"
	"net/http"
	"testing"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestScanNetwork(t *testing.T) {
	http.HandleFunc("/", handler)
	go http.ListenAndServe(":8080", nil)
	type args struct {
		port string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "find ip",
			args:    args{":8080"},
			want:    net.Dial,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ScanNetwork(tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ScanNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}
