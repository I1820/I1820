package request

import "testing"

func TestThing_Validate(t *testing.T) {
	tests := []struct {
		name    string
		fields  Thing
		wantErr bool
	}{
		{
			name: "with model",
			fields: Thing{
				Name:  "raha",
				Model: "aolab",
			},
			wantErr: false,
		},
		{
			name: "without model",
			fields: Thing{
				Name: "raha",
			},
			wantErr: false,
		},
		{
			name: "invalid location",
			fields: Thing{
				Name: "raha",
				Location: Location{
					Latitude:  1024,
					Longitude: 30,
				},
			},
			wantErr: true,
		},
		{
			name: "valid location",
			fields: Thing{
				Name: "raha",
				Location: Location{
					Latitude:  35.807657,
					Longitude: 51.398408,
				},
			},
			wantErr: false,
		},
	}

	// nolint: scopelint
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thing := Thing{
				Name:     tt.fields.Name,
				Model:    tt.fields.Model,
				Location: tt.fields.Location,
			}
			if err := thing.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
