package wykop

import (
	"strings"

	"github.com/xaxes/vikop-gorace/pkg/dedup"
)

type FilterGroup []Filter

func (fg FilterGroup) Filter(es ...Entry) []Entry {
	for _, f := range fg {
		es = f(es...)
	}

	return es
}

type Filter func(...Entry) []Entry

func blacklisted(blacklist []string, e Entry) bool {
	for _, tag := range blacklist {
		if strings.Contains(e.Body, "#"+tag) {
			return true
		}
	}
	return false
}

func NewBlacklistFilter(blacklist []string) Filter {
	return func(es ...Entry) []Entry {
		var filtered []Entry
		for _, e := range es {
			if !blacklisted(blacklist, e) {
				filtered = append(filtered, e)
			}
		}

		return filtered
	}
}

func NewDuplicatesFilter(deduper dedup.Deduper) Filter {
	return func(es ...Entry) []Entry {
		var deduped []Entry
		for _, e := range es {
			if !deduper.Sent(e.ID) {
				deduped = append(deduped, e)

				if err := deduper.SetSent(e.ID, periodToTTL(TwentyFour)); err != nil {
					return nil
				}
			}
		}

		return deduped
	}
}

func NewMinVotesFilter(minVotes int) Filter {
	return func(es ...Entry) []Entry {
		var res []Entry
		for _, e := range es {
			if e.VoteCount > minVotes {
				res = append(res, e)
			}
		}
		return res
	}
}
