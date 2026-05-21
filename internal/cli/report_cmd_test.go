package cli

import "testing"

func TestSplitCSVTrimsEmptyValues(t *testing.T) {
	got := splitCSV(" users, sessions, ,screenPageViews ")
	want := []string{"users", "sessions", "screenPageViews"}
	if len(got) != len(want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	}
}

func TestParseOrderBys(t *testing.T) {
	got := parseOrderBys("-sessions,eventCount")
	if len(got) != 2 {
		t.Fatalf("got %#v", got)
	}
	if got[0].Field != "sessions" || !got[0].Desc {
		t.Fatalf("unexpected first order by: %#v", got[0])
	}
	if got[1].Field != "eventCount" || got[1].Desc {
		t.Fatalf("unexpected second order by: %#v", got[1])
	}
}
