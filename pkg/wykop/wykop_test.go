package wykop

import (
	"reflect"
	"testing"
)

func Test_filter(t *testing.T) {
	type args struct {
		blacklist []string
		entries   []Entry
	}
	tests := []struct {
		name string
		args args
		want []Entry
	}{
		{
			name: "has tag",
			args: args{
				blacklist: []string{"blacklisted1"},
				entries: []Entry{{
					ID:            0,
					URL:           "",
					Date:          "",
					Body:          "whatever #someTag whatever",
					Author:        Author{},
					VoteCount:     0,
					CommentsCount: 0,
					Embed:         Embed{},
				}},
			},
			want: nil,
		},
		{
			name: "has no blacklisted tag",
			args: args{
				blacklist: []string{"rozdajo"},
				entries: []Entry{{
					ID:            0,
					URL:           "",
					Date:          "",
					Body:          "whatever #someTag whatever",
					Author:        Author{},
					VoteCount:     0,
					CommentsCount: 0,
					Embed:         Embed{},
				}},
			},
			want: []Entry{{
				ID:            0,
				URL:           "",
				Date:          "",
				Body:          "whatever #someTag whatever",
				Author:        Author{},
				VoteCount:     0,
				CommentsCount: 0,
				Embed:         Embed{},
			}},
		},
		{
			name: "mixed",
			args: args{
				blacklist: []string{"rozdajo"},
				entries: []Entry{{
					ID:            0,
					URL:           "",
					Date:          "",
					Body:          "whatever #someTag whatever #anotherTag ipsum",
					Author:        Author{},
					VoteCount:     0,
					CommentsCount: 0,
					Embed:         Embed{},
				}},
			},
			want: nil,
		},
		{
			name: "mixed many",
			args: args{
				blacklist: []string{"f1"},
				entries: []Entry{
					{
						ID:            0,
						URL:           "",
						Date:          "",
						Body:          "whatever #someTag whatever #anotherTag ipsum",
						Author:        Author{},
						VoteCount:     0,
						CommentsCount: 0,
						Embed:         Embed{},
					},
					{
						ID:            0,
						URL:           "",
						Date:          "",
						Body:          "lorem #okTag ipsum",
						Author:        Author{},
						VoteCount:     0,
						CommentsCount: 0,
						Embed:         Embed{},
					},
				},
			},
			want: []Entry{
				{
					ID:            0,
					URL:           "",
					Date:          "",
					Body:          "whatever #someTag whatever #anotherTag ipsum",
					Author:        Author{},
					VoteCount:     0,
					CommentsCount: 0,
					Embed:         Embed{},
				}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewBlacklistFilter(tt.args.blacklist)
			if got := filter(tt.args.entries...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
