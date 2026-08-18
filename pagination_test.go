package signater

import (
	"errors"
	"testing"
)

func TestListParamsDefaults(t *testing.T) {
	v := ListParams{}.values()
	if v.Get("PageSize") != "10" || v.Get("PageNumber") != "1" || v.Get("OrderByDirection") != "DESC" {
		t.Errorf("defaults = %v", v)
	}
}

func TestListParamsExplicit(t *testing.T) {
	v := ListParams{PageSize: 50, PageNumber: 3, OrderByDirection: OrderAsc}.values()
	if v.Get("PageSize") != "50" || v.Get("PageNumber") != "3" || v.Get("OrderByDirection") != "ASC" {
		t.Errorf("values = %v", v)
	}
}

// fakePages simulates a server with n items served in pages of size.
func fakePages(total, size int, calls *int) func(page int) (*Page[int], error) {
	return func(page int) (*Page[int], error) {
		*calls++
		start := (page - 1) * size
		var items []int
		for i := start; i < min(start+size, total); i++ {
			items = append(items, i)
		}
		return &Page[int]{
			Items: items,
			Pagination: Pagination{
				TotalItems: total, PageSize: size, PageNumber: page, PageItems: len(items),
			},
		}, nil
	}
}

func TestAutoPagingAllPages(t *testing.T) {
	var calls int
	var got []int
	for item, err := range autoPaging(1, fakePages(25, 10, &calls)) {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, item)
	}
	if len(got) != 25 || got[0] != 0 || got[24] != 24 {
		t.Errorf("items = %d, first %d, last %d", len(got), got[0], got[len(got)-1])
	}
	if calls != 3 {
		t.Errorf("fetch calls = %d, want 3", calls)
	}
}

func TestAutoPagingExactPageBoundary(t *testing.T) {
	// 20 items in pages of 10: totalItems tells the iterator page 2 is the
	// last one, so no extra empty fetch is issued.
	var calls int
	count := 0
	for _, err := range autoPaging(1, fakePages(20, 10, &calls)) {
		if err != nil {
			t.Fatal(err)
		}
		count++
	}
	if count != 20 {
		t.Errorf("items = %d, want 20", count)
	}
	if calls != 2 {
		t.Errorf("fetch calls = %d, want 2 (totalItems saves the empty fetch)", calls)
	}
}

func TestAutoPagingNilPage(t *testing.T) {
	fetch := func(_ int) (*Page[int], error) { return nil, nil }
	count := 0
	for _, err := range autoPaging(1, fetch) {
		if err != nil {
			t.Fatal(err)
		}
		count++
	}
	if count != 0 {
		t.Errorf("items = %d, want 0 (nil page must terminate, not panic)", count)
	}
}

func TestAutoPagingMissingPaginationMetadata(t *testing.T) {
	// A non-empty page with a zero-valued pagination block must stop after
	// one fetch instead of looping forever.
	var calls int
	fetch := func(_ int) (*Page[int], error) {
		calls++
		return &Page[int]{Items: []int{1, 2, 3}}, nil
	}
	count := 0
	for _, err := range autoPaging(1, fetch) {
		if err != nil {
			t.Fatal(err)
		}
		count++
	}
	if count != 3 || calls != 1 {
		t.Errorf("items = %d, calls = %d; want 3 items from exactly 1 call", count, calls)
	}
}

func TestAutoPagingEarlyBreak(t *testing.T) {
	var calls int
	count := 0
	for _, err := range autoPaging(1, fakePages(100, 10, &calls)) {
		if err != nil {
			t.Fatal(err)
		}
		count++
		if count == 5 {
			break
		}
	}
	if calls != 1 {
		t.Errorf("fetch calls = %d, want 1 (break must stop fetching)", calls)
	}
}

func TestAutoPagingErrorPropagates(t *testing.T) {
	boom := errors.New("boom")
	fetch := func(page int) (*Page[int], error) {
		if page == 2 {
			return nil, boom
		}
		return &Page[int]{
			Items:      []int{1, 2},
			Pagination: Pagination{PageSize: 2, PageItems: 2, PageNumber: page},
		}, nil
	}
	var seen int
	var lastErr error
	for _, err := range autoPaging(1, fetch) {
		if err != nil {
			lastErr = err
			break
		}
		seen++
	}
	if !errors.Is(lastErr, boom) {
		t.Errorf("lastErr = %v, want boom", lastErr)
	}
	if seen != 2 {
		t.Errorf("items before error = %d, want 2", seen)
	}
}

func TestAutoPagingStartPage(t *testing.T) {
	fetch := func(page int) (*Page[int], error) {
		return &Page[int]{
			Items:      []int{page},
			Pagination: Pagination{PageSize: 10, PageItems: 1, PageNumber: page},
		}, nil
	}
	for item, err := range autoPaging(4, fetch) {
		if err != nil {
			t.Fatal(err)
		}
		if item != 4 {
			t.Errorf("first item = %d, want 4 (must start at firstPage)", item)
		}
		break
	}
}
