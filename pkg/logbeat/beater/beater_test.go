package beater

import (
	"github.com/elastic/beats/libbeat/beat"
	"sync"
	"testing"
	"time"
)

func TestLogbeat_Run(t *testing.T) {
	type fields struct {
		period  time.Duration
		client  beat.Client
		handler Handler
		done    chan struct{}
	}
	type args struct {
		b *beat.Beat
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal test",
			fields: fields{
				period:  time.Second,
				client:  new(mockBeatClient),
				handler: newMockHandler(),
				done:    nil,
			},
			args: args{b: &mockBeat},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := &Logbeat{
				period:  tt.fields.period,
				client:  tt.fields.client,
				handler: tt.fields.handler,
				done:    tt.fields.done,
			}
			wg := new(sync.WaitGroup)
			wg.Add(1)
			go func() {
				if err := lb.Run(tt.args.b); (err != nil) != tt.wantErr {
					t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				}
				wg.Done()
			}()
			time.Sleep(tt.fields.period * 2)
			lb.Stop()
			wg.Wait()
		})
	}
}
