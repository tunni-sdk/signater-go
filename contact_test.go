package signater

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

// listContactsFixture is the documented example response of List contacts.
const listContactsFixture = `{
  "items": [
    {
      "id": "99c8a78c-71eb-4e20-a758-99082ee63ef9",
      "createdAtUtc": "2026-08-17T00:52:53.7112783Z",
      "createdByName": "John Doe",
      "createdByAvatar": "https://example.com/avatar.jpg",
      "updatedAtUtc": "2026-08-17T00:52:53.7112785Z",
      "updatedByName": "John Doe",
      "updatedByAvatar": "https://example.com/avatar.jpg",
      "name": "John Doe",
      "title": "Company Director",
      "email": "john.doe@example.com",
      "phoneIdd": 1,
      "phoneNumber": "5555555555",
      "documentType": "GenericIdentification",
      "documentValue": "123456789",
      "isFavorite": true
    }
  ],
  "pagination": {"totalItems": 1, "pageSize": 10, "pageNumber": 1, "pageItems": 1}
}`

func TestContactList(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusOK, listContactsFixture, nil)

	page, err := c.Contacts.List(context.Background(), &ContactListParams{
		IsFavorite: Bool(true),
		Search:     "john",
	})
	if err != nil {
		t.Fatal(err)
	}

	req := tr.Snapshot()[0]
	if req.Method != http.MethodGet || req.Path != "/v1/ecm/contacts" {
		t.Errorf("request = %s %s", req.Method, req.Path)
	}
	for k, want := range map[string]string{
		"IsFavorite": "true", "Search": "john",
		"PageSize": "10", "PageNumber": "1", "OrderByDirection": "DESC",
	} {
		if got := queryOf(t, req.Query).Get(k); got != want {
			t.Errorf("query %s = %q, want %q", k, got, want)
		}
	}

	if len(page.Items) != 1 {
		t.Fatalf("items = %d", len(page.Items))
	}
	got := page.Items[0]
	if got.ID != "99c8a78c-71eb-4e20-a758-99082ee63ef9" || got.Name != "John Doe" ||
		got.PhoneIDD != 1 || got.DocumentType != IdentityDocumentGenericIdentification ||
		!got.IsFavorite || got.CreatedAt.IsZero() {
		t.Errorf("decoded contact = %+v", got)
	}
	if page.Pagination.TotalItems != 1 {
		t.Errorf("pagination = %+v", page.Pagination)
	}
}

func TestContactListAutoPaging(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusOK, `{"items":[{"id":"a"},{"id":"b"}],"pagination":{"totalItems":3,"pageSize":2,"pageNumber":1,"pageItems":2}}`, nil)
	tr.Enqueue(http.StatusOK, `{"items":[{"id":"c"}],"pagination":{"totalItems":3,"pageSize":2,"pageNumber":2,"pageItems":1}}`, nil)

	var ids []string
	for contact, err := range c.Contacts.ListAutoPaging(context.Background(), &ContactListParams{ListParams: ListParams{PageSize: 2}}) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, contact.ID)
	}
	if len(ids) != 3 || ids[2] != "c" {
		t.Errorf("ids = %v", ids)
	}
	reqs := tr.Snapshot()
	if len(reqs) != 2 {
		t.Fatalf("requests = %d, want 2", len(reqs))
	}
	if got := queryOf(t, reqs[1].Query).Get("PageNumber"); got != "2" {
		t.Errorf("second request PageNumber = %q", got)
	}
}

func TestContactCreate(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusCreated, `{"id":"e303178d-de2c-42de-8f89-4f7dcac95beb"}`, nil)

	id, err := c.Contacts.Create(context.Background(), &ContactParams{
		Name:         "John Doe",
		Title:        String("Project Manager"),
		Email:        String("john.doe@example.com"),
		PhoneIDD:     Int(1),
		PhoneNumber:  String("1234567890"),
		DocumentType: IdentityDocumentGenericIdentification,
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != "e303178d-de2c-42de-8f89-4f7dcac95beb" {
		t.Errorf("id = %q", id)
	}

	req := tr.Snapshot()[0]
	if req.Method != http.MethodPost || req.Path != "/v1/ecm/contacts" {
		t.Errorf("request = %s %s", req.Method, req.Path)
	}
	var body map[string]any
	if err := json.Unmarshal(req.Body, &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "John Doe" || body["phoneIdd"] != float64(1) || body["documentType"] != "GenericIdentification" {
		t.Errorf("body = %v", body)
	}
	if _, present := body["documentValue"]; present {
		t.Error("nil optional field must be omitted from the body")
	}
}

func TestContactGetUpdateRemoveFavorite(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusOK, `{"id":"x1","name":"Jane"}`, nil)
	tr.Enqueue(http.StatusNoContent, "", nil)
	tr.Enqueue(http.StatusNoContent, "", nil)
	tr.Enqueue(http.StatusNoContent, "", nil)
	tr.Enqueue(http.StatusNoContent, "", nil)

	ctx := context.Background()
	contact, err := c.Contacts.Get(ctx, "x1")
	if err != nil || contact.Name != "Jane" {
		t.Fatalf("get = %+v, %v", contact, err)
	}
	if err := c.Contacts.Update(ctx, "x1", &ContactParams{Name: "Jane B", Email: String("jane@x.com")}); err != nil {
		t.Fatal(err)
	}
	if err := c.Contacts.Remove(ctx, "x1"); err != nil {
		t.Fatal(err)
	}
	if err := c.Contacts.Favorite(ctx, "x1"); err != nil {
		t.Fatal(err)
	}
	if err := c.Contacts.Unfavorite(ctx, "x1"); err != nil {
		t.Fatal(err)
	}

	reqs := tr.Snapshot()
	wantReqs := []struct{ method, path string }{
		{http.MethodGet, "/v1/ecm/contacts/x1"},
		{http.MethodPut, "/v1/ecm/contacts/x1"},
		{http.MethodDelete, "/v1/ecm/contacts/x1"},
		{http.MethodPost, "/v1/ecm/contacts/x1/favorite"},
		{http.MethodPost, "/v1/ecm/contacts/x1/unfavorite"},
	}
	for i, want := range wantReqs {
		if reqs[i].Method != want.method || reqs[i].Path != want.path {
			t.Errorf("request %d = %s %s, want %s %s", i, reqs[i].Method, reqs[i].Path, want.method, want.path)
		}
	}
	var updateBody map[string]any
	if err := json.Unmarshal(reqs[1].Body, &updateBody); err != nil {
		t.Fatal(err)
	}
	if updateBody["name"] != "Jane B" || updateBody["email"] != "jane@x.com" {
		t.Errorf("update body = %v", updateBody)
	}
}

func TestContactIDIsPathEscaped(t *testing.T) {
	// Dot and slash segments must never collapse the URL onto another route.
	cases := map[string]string{
		"a/b":      "/v1/ecm/contacts/a%2Fb",
		"../admin": "/v1/ecm/contacts/..%2Fadmin",
		".":        "/v1/ecm/contacts/%2E",
		"..":       "/v1/ecm/contacts/%2E%2E",
	}
	for id, wantPath := range cases {
		c, tr := newTestClient(t)
		tr.Enqueue(http.StatusNotFound, "", nil)
		if _, err := c.Contacts.Get(context.Background(), id); err == nil {
			t.Fatalf("id %q: expected 404 error", id)
		}
		if got := tr.Snapshot()[0].Path; got != wantPath {
			t.Errorf("id %q: path = %s, want %s", id, got, wantPath)
		}
	}
}

func TestEmptyIDIsRejectedClientSide(t *testing.T) {
	c, tr := newTestClient(t)
	_, err := c.Contacts.Get(context.Background(), "")
	if err != errEmptyResourceID {
		t.Fatalf("err = %v, want errEmptyResourceID", err)
	}
	// The empty-id error must fire before any request reaches the wire.
	if tr.Count() != 0 {
		t.Errorf("requests = %d, want 0", tr.Count())
	}
}
