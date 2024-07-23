package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
)

func TestHandleRequest(t *testing.T) {

	err := godotenv.Load("../../.TEST.env")
	if err != nil {
		t.Error(err)
		return
	}

	type args struct {
		ctx   context.Context
		event *ContainsSwearwordsEvent
	}
	tests := []struct {
		name    string
		args    args
		want    []string
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
			want:    []string{},
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
			want:    []string{"Depp"},
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
			want:    []string{"Depp"},
			wantErr: false,
		},
		{
			name: "A text with several hundred words",
			args: args{
				ctx: context.Background(),
				event: &ContainsSwearwordsEvent{
					Text: getSeveralHundredWords(),
				},
			},
			want:    []string{},
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
			if !reflect.DeepEqual(got.Swearwords, tt.want) {
				t.Errorf("HandleRequest() = %v, want %v", got.Swearwords, tt.want)
			}
		})
	}
}

func getSeveralHundredWords() string {
	text := "This is a test with several hundred words: "
	text += `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur feugiat ligula eu magna bibendum aliquam. Proin semper nisi sit amet justo efficitur, ac ullamcorper augue mollis. Aliquam erat volutpat. Sed auctor felis a libero ultrices, vel vehicula nisl ultrices. Ut vitae elit sed nulla blandit sodales sed ac felis. Nunc faucibus ante tellus, vitae auctor felis dignissim nec. Aenean malesuada elit sed purus blandit, eget accumsan felis pretium. Suspendisse sit amet porttitor dui. In at metus venenatis, sagittis lacus eget, posuere diam. Vestibulum blandit id sapien mollis mollis. Nulla quam lorem, scelerisque in risus vel, eleifend semper mi. Quisque in leo sit amet tortor semper luctus id at eros. Curabitur rhoncus enim vel augue placerat, non tincidunt dui aliquam. Nam at risus vitae elit auctor dignissim ut et urna.

Integer tincidunt ac erat in tempor. Ut lacinia neque at lectus ultricies, at bibendum mauris tempus. Curabitur ut lectus gravida, blandit orci non, imperdiet dui. Aenean luctus sed libero consectetur condimentum. Cras pellentesque odio felis, venenatis cursus urna varius sed. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Duis ac metus est. Praesent commodo justo non lobortis aliquet. Nulla facilisi. Etiam egestas, enim eget hendrerit ultrices, justo turpis tincidunt sapien, at pellentesque risus massa a ipsum. Suspendisse mollis mauris eu odio vestibulum, hendrerit accumsan lorem dapibus. Integer quis dolor lacus. Aenean erat eros, tempor quis tempor eu, mollis at nunc. Phasellus nunc leo, interdum ac sem eget, faucibus tempus massa.

Etiam ac quam pharetra, dapibus diam ut, iaculis sem. Suspendisse ac cursus elit, non consectetur magna. Phasellus condimentum magna libero, ut imperdiet neque lobortis eget. Etiam rutrum a mi vel feugiat. Pellentesque imperdiet justo in enim gravida ultrices. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur facilisis suscipit pulvinar.

Curabitur quis efficitur sapien. Etiam consectetur ut odio sed ullamcorper. Proin a lorem dictum, viverra elit et, commodo ligula. Sed eget nisi facilisis, bibendum justo vel, lacinia purus. Ut id euismod eros. Nunc sed tellus id odio feugiat lobortis in sed enim. Integer porta odio eu fermentum malesuada. Nulla facilisi. Aliquam vitae elit risus. Interdum et malesuada fames ac ante ipsum primis in faucibus.

Morbi a semper mi. Duis at fringilla ipsum, sed vehicula tellus. Phasellus porttitor tellus non mi sollicitudin efficitur. Quisque eget lacus in augue venenatis pharetra convallis euismod nunc. Mauris a magna malesuada, malesuada nibh sed, ornare odio. Maecenas tellus sem, varius in nunc ut, luctus volutpat mi. Nulla facilisi. Aenean tincidunt eget diam ac pretium.`
	return text
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
