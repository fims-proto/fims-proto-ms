package class

import "testing"

func TestValidate(t *testing.T) {
	type args struct {
		c Class
		g Group
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal success",
			args: args{
				c: 1,
				g: 101,
			},
			wantErr: false,
		},
		{
			name: "invalid class",
			args: args{
				c: 0,
				g: 101,
			},
			wantErr: true,
		},
		{
			name: "invalid group",
			args: args{
				c: 1,
				g: 201,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.args.c, tt.args.g); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
