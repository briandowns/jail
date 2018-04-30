package jail

import (
	"reflect"
	"testing"
)

func TestOpts_validate(t *testing.T) {
	type fields struct {
		Version  int
		Path     string
		Name     string
		Hostname string
		IP4      string
		Chdir    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Opts{
				Version:  tt.fields.Version,
				Path:     tt.fields.Path,
				Name:     tt.fields.Name,
				Hostname: tt.fields.Hostname,
				IP4:      tt.fields.IP4,
				Chdir:    tt.fields.Chdir,
			}
			if err := o.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Opts.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJail(t *testing.T) {
	type args struct {
		o *Opts
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Jail(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("Jail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Jail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewParams(t *testing.T) {
	tests := []struct {
		name string
		want Params
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewParams(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttach(t *testing.T) {
	type args struct {
		jailID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Attach(tt.args.jailID); (err != nil) != tt.wantErr {
				t.Errorf("Attach() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	type args struct {
		jailID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Remove(tt.args.jailID); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
