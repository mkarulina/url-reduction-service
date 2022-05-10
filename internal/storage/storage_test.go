package storage

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func Test_storage_GetAllUrls(t *testing.T) {
	type fields struct {
		mu   sync.Mutex
		urls map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Link
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				mu: sync.Mutex{},
				urls: map[string]string{
					"key1": "http://ya.ru/1",
					"key2": "http://ya.ru/2",
					"key3": "http://ya.ru/3",
				},
			},
			want: []Link{
				{Key: "key1", Link: "http://ya.ru/1"},
				{Key: "key2", Link: "http://ya.ru/2"},
				{Key: "key3", Link: "http://ya.ru/3"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				mu:   tt.fields.mu,
				urls: tt.fields.urls,
			}
			got, err := s.GetAllUrls()
			require.NoError(t, err)
			require.Equal(t, len(got), len(tt.want))

			for i := 0; i < len(tt.want); i++ {
				require.Contains(t, got, tt.want[i])
			}
		})
	}
}

func Test_storage_GetKeyByLink(t *testing.T) {
	type fields struct {
		mu   sync.Mutex
		urls map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		link   string
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				mu: sync.Mutex{},
				urls: map[string]string{
					"key1": "http://ya.ru/1",
					"key2": "http://ya.ru/2",
					"key3": "http://ya.ru/3",
				},
			},
			link: "http://ya.ru/2",
			want: "key2",
		},
		{
			name: "not existent link",
			fields: fields{
				mu: sync.Mutex{},
				urls: map[string]string{
					"key1": "http://ya.ru/1",
					"key2": "http://ya.ru/2",
					"key3": "http://ya.ru/3",
				},
			},
			link: "http://notExistent.ru",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				mu:   tt.fields.mu,
				urls: tt.fields.urls,
			}
			if got := s.GetKeyByLink(tt.link); got != tt.want {
				t.Errorf("GetKeyByLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storage_GetLinkByKey(t *testing.T) {
	type fields struct {
		mu   sync.Mutex
		urls map[string]string
	}

	tests := []struct {
		name   string
		fields fields
		key    string
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				mu: sync.Mutex{},
				urls: map[string]string{
					"key1": "http://ya.ru/1",
					"key2": "http://ya.ru/2",
					"key3": "http://ya.ru/3",
				},
			},
			key:  "key3",
			want: "http://ya.ru/3",
		},
		{
			name: "not existent key",
			fields: fields{
				mu: sync.Mutex{},
				urls: map[string]string{
					"key1": "http://ya.ru/4",
					"key2": "http://ya.ru/5",
					"key3": "http://ya.ru/6",
				},
			},
			key:  "notExistent",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				mu:   tt.fields.mu,
				urls: tt.fields.urls,
			}
			if got := s.GetLinkByKey(tt.key); got != tt.want {
				t.Errorf("GetLinkByKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
