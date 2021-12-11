package wykop

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/xaxes/vikop-gorace/pkg/dedup"
	"go.uber.org/zap"
)

const baseURL = "https://a2.wykop.pl/"

type Period int

const (
	Six        Period = 6
	Twelve     Period = 12
	TwentyFour Period = 24
)

type Author struct {
	Login  string
	Color  int
	Gender string
	Avatar string
}

type Embed struct {
	Type string
	URL  string
}

type Entry struct {
	ID            int
	URL           string `json:"url"`
	Date          string
	Body          string
	Author        Author
	VoteCount     int `json:"vote_count"`
	CommentsCount int `json:"comments_count"`
	Embed         Embed
}

type Wykop struct {
	l       *zap.Logger
	appkey  string
	deduper dedup.Deduper
}

func New(l *zap.Logger, appkey string, deduper dedup.Deduper) Wykop {
	return Wykop{l, appkey, deduper}
}

func (w Wykop) url(fun string, namedParams map[string]string) string {
	namedParams["appkey"] = w.appkey
	namedParams["output"] = "clear"
	namedParams["data"] = "full"

	res := baseURL + fun

	for param, val := range namedParams {
		res += fmt.Sprintf("/%s/%s", param, val)
	}

	return res
}

var counter int

func (w Wykop) query(fun string, namedParams map[string]string) ([]Entry, error) {
	counter++
	if counter >= 100 {
		panic("counter 100")
	}

	finalURL := w.url(fun, namedParams)
	resp, err := http.Get(finalURL)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var entriesRaw struct {
		Entries []Entry `json:"data"`
	}
	if err := json.Unmarshal(body, &entriesRaw); err != nil {
		return nil, err
	}

	return entriesRaw.Entries, nil

}

func periodToTTL(p Period) time.Duration {
	return time.Duration(int(time.Hour) * int(p))
}

func (w Wykop) Hot(period Period, noOfEntries, maxPage int, blacklist []string, minVotes int) ([]Entry, error) {
	fg := FilterGroup{
		NewMinVotesFilter(minVotes),
		NewBlacklistFilter(blacklist),
		NewDuplicatesFilter(w.deduper),
	}

	var entries []Entry
	for page := 1; len(entries) < noOfEntries && page < maxPage; page++ {
		es, err := w.query("Entries/Hot", map[string]string{
			"page":   strconv.Itoa(page),
			"period": strconv.Itoa(int(period)),
		})
		if err != nil {
			return nil, err
		}
		w.l.Debug("collect entries from page", zap.Int("collected", len(es)), zap.Int("page", page), zap.Int("period", int(period)))

		es = fg.Filter(es...)
		w.l.Debug("filter entries", zap.Int("left", len(es)), zap.Int("period", int(period)))

		entries = append(entries, es...)
	}
	w.l.Debug("collect entries", zap.Int("collected", len(entries)), zap.Int("period", int(period)))

	if len(entries) < noOfEntries {
		return entries, nil
	}
	return entries[:noOfEntries], nil
}
