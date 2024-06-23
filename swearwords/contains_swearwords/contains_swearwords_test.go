package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
)

func TestHandleRequest(t *testing.T) {

	godotenv.Load("../../.env")

	type args struct {
		ctx   context.Context
		event *ContainsSwearwordsEvent
	}
	tests := []struct {
		name    string
		args    args
		want    *[]string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "A text without swearwords",
			args: args{
				ctx: context.Background(),
				event: &ContainsSwearwordsEvent{
					Text: "This is a test without swearwords",
				},
			},
			want:    &[]string{},
			wantErr: false,
		},
		{
			name: "A text with one swearword",
			args: args{
				ctx: context.Background(),
				event: &ContainsSwearwordsEvent{
					Text: "This is a test with a swearword: Depp",
				},
			},
			want:    &[]string{"Depp"},
			wantErr: false,
		},
		{
			name: "A text with one swearword with punctuation",
			args: args{
				ctx: context.Background(),
				event: &ContainsSwearwordsEvent{
					Text: "This is a test with a swearword: Depp!",
				},
			},
			want:    &[]string{"Depp"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandleRequest(tt.args.ctx, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removePunctuation(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "A word without punctuation",
			args: args{
				word: "Word",
			},
			want: "Word",
		},
		{
			name: "A word with punctuation at the start",
			args: args{
				word: "!Word",
			},
			want: "Word",
		},
		{
			name: "A word with punctuation at the end",
			args: args{
				word: "Word!",
			},
			want: "Word",
		},
		{
			name: "A word with punctuation at the start and end",
			args: args{
				word: "&)(&!Word))((!",
			},
			want: "Word",
		},
		{
			name: "A word with punctuation in the middle",
			args: args{
				word: "Wo!rd",
			},
			want: "Wo!rd",
		},
		{
			name: "A word with whitespace at the start",
			args: args{
				word: " Word",
			},
			want: "Word",
		},
		{
			name: "An empty string",
			args: args{
				word: "",
			},
			want: "",
		},
		{
			name: "Only whitespace",
			args: args{
				word: "   ",
			},
			want: "",
		},
		{
			name: "Only punctuation",
			args: args{
				word: "!!/&&%==))",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removePunctuation(tt.args.word); got != tt.want {
				t.Errorf("removePunctuation() = %v, want %v", got, tt.want)
			}
		})
	}
}
