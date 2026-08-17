package signater

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

// queryOf parses a recorded raw query string.
func queryOf(t *testing.T, raw string) url.Values {
	t.Helper()
	v, err := url.ParseQuery(raw)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestVaultCreate(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusCreated, `{"id":"b01eabec-dd0f-4c23-b537-65f40a1e216d"}`, nil)

	id, err := c.Vaults.Create(context.Background(), &VaultParams{
		Name:        "My Vault",
		Description: String("My Vault Description"),
		Type:        VaultTypeUserAccountGroup,
		MemberIDs:   []string{"136ad90e-2f5e-46d1-85cf-b0b94457f334"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != "b01eabec-dd0f-4c23-b537-65f40a1e216d" {
		t.Errorf("id = %q", id)
	}

	var body map[string]any
	if err := json.Unmarshal(tr.Snapshot()[0].Body, &body); err != nil {
		t.Fatal(err)
	}
	if body["type"] != "UserAccountGroup" {
		t.Errorf("body type = %v", body["type"])
	}
	members, ok := body["memberIds"].([]any)
	if !ok || len(members) != 1 {
		t.Errorf("memberIds = %v", body["memberIds"])
	}
}

func TestVaultGet(t *testing.T) {
	// Documented example response of Get vault.
	const fixture = `{
	  "id": "70494cbf-b075-476b-8fa4-974901488527",
	  "name": "Example Vault",
	  "type": "UserAccount",
	  "description": "This is an example of a vault used for demonstration purposes.",
	  "userAccountMembers": [
	    {"id": "69e9cc6a-574b-4cc4-b55d-93108d028082", "name": "John Doe", "avatar": null},
	    {"id": "e3233499-262f-4975-b6d4-dfef027983b2", "name": "Jane Smith", "avatar": "https://example.com/avatars/janesmith.png"}
	  ]
	}`
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusOK, fixture, nil)

	v, err := c.Vaults.Get(context.Background(), "70494cbf-b075-476b-8fa4-974901488527")
	if err != nil {
		t.Fatal(err)
	}
	if v.Type != VaultTypeUserAccount || len(v.UserAccountMembers) != 2 || v.UserAccountMembers[1].Name != "Jane Smith" {
		t.Errorf("vault = %+v", v)
	}
	if req := tr.Snapshot()[0]; req.Path != "/v1/ecm/vaults/70494cbf-b075-476b-8fa4-974901488527" {
		t.Errorf("path = %s", req.Path)
	}
}

func TestVaultUpdateMethods(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusNoContent, "", nil)
	tr.Enqueue(http.StatusNoContent, "", nil)
	tr.Enqueue(http.StatusNoContent, "", nil)

	ctx := context.Background()
	params := &VaultParams{Name: "New name", Type: VaultTypeAccount}
	if err := c.Vaults.Update(ctx, "v1", params); err != nil {
		t.Fatal(err)
	}
	if err := c.Vaults.UpdatePartial(ctx, "v1", params); err != nil {
		t.Fatal(err)
	}
	if err := c.Vaults.Remove(ctx, "v1"); err != nil {
		t.Fatal(err)
	}

	reqs := tr.Snapshot()
	if reqs[0].Method != http.MethodPut || reqs[1].Method != http.MethodPatch || reqs[2].Method != http.MethodDelete {
		t.Errorf("methods = %s, %s, %s; want PUT, PATCH, DELETE", reqs[0].Method, reqs[1].Method, reqs[2].Method)
	}
	for _, r := range reqs {
		if r.Path != "/v1/ecm/vaults/v1" {
			t.Errorf("path = %s", r.Path)
		}
	}
	// PUT and PATCH share the wire format; type must be present on both.
	for i := range 2 {
		var body map[string]any
		if err := json.Unmarshal(reqs[i].Body, &body); err != nil {
			t.Fatal(err)
		}
		if body["name"] != "New name" || body["type"] != "Account" {
			t.Errorf("request %d body = %v", i, body)
		}
		// memberIds has no omitempty: nil serializes as null so an empty
		// non-nil slice can clear the member list.
		if v, present := body["memberIds"]; !present || v != nil {
			t.Errorf("request %d memberIds = %v (present=%v), want explicit null", i, v, present)
		}
	}
}

func TestVaultListAutoPagingSnapshotsParams(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusOK, `{"items":[{"id":"v1"}],"pagination":{"totalItems":2,"pageSize":1,"pageNumber":1,"pageItems":1}}`, nil)
	tr.Enqueue(http.StatusOK, `{"items":[{"id":"v2"}],"pagination":{"totalItems":2,"pageSize":1,"pageNumber":2,"pageItems":1}}`, nil)

	params := &VaultListParams{ListParams: ListParams{PageSize: 1}, Search: "x"}
	var ids []string
	for v, err := range c.Vaults.ListAccountAutoPaging(context.Background(), params) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, v.ID)
		params.Search = "mutated-mid-iteration" // must not affect later fetches
	}
	if len(ids) != 2 || ids[1] != "v2" {
		t.Errorf("ids = %v", ids)
	}
	reqs := tr.Snapshot()
	q2 := queryOf(t, reqs[1].Query)
	if q2.Get("PageNumber") != "2" {
		t.Errorf("second request PageNumber = %q", q2.Get("PageNumber"))
	}
	if q2.Get("Search") != "x" {
		t.Errorf("second request Search = %q; params mutation leaked into the iterator", q2.Get("Search"))
	}
	if params.PageNumber != 0 {
		t.Errorf("caller params mutated: PageNumber = %d", params.PageNumber)
	}
}

func TestVaultOwnersAndMembers(t *testing.T) {
	const fixture = `{"items":[{"id":"f43cff0f-df08-4a71-b91d-ee64de8ef103","name":"John Doe","isActive":true,"isCurrent":true,"role":"Administrator","avatar":"https://example.com/avatar.jpg"}]}`
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusOK, fixture, nil)
	tr.Enqueue(http.StatusOK, fixture, nil)

	ctx := context.Background()
	owners, err := c.Vaults.Owners(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(owners) != 1 || owners[0].Role != UserRoleAdministrator || !owners[0].IsCurrent {
		t.Errorf("owners = %+v", owners)
	}
	if _, err := c.Vaults.Members(ctx); err != nil {
		t.Fatal(err)
	}

	reqs := tr.Snapshot()
	if reqs[0].Path != "/v1/ecm/vaults/owners" || reqs[1].Path != "/v1/ecm/vaults/members" {
		t.Errorf("paths = %s, %s", reqs[0].Path, reqs[1].Path)
	}
}

func TestVaultListEndpoints(t *testing.T) {
	// Documented example item of the vault list endpoints.
	const fixture = `{
	  "items": [{
	    "id": "f39a90a4-cc59-4097-9764-d549ea93648e",
	    "createdAtUtc": "2026-08-17T00:52:53.9840363Z",
	    "name": "Example Vault",
	    "type": "UserAccount",
	    "isFavorite": false,
	    "canBeRemoved": true,
	    "userAccountMembers": [{"id": "715c4a9d-cdf4-45e8-8428-32413d9a7c9e", "name": "John Doe", "email": "john.doe@example.com", "avatar": null, "isActive": true}]
	  }],
	  "pagination": {"totalItems": 1, "pageSize": 10, "pageNumber": 1, "pageItems": 1}
	}`
	c, tr := newTestClient(t)
	for range 3 {
		tr.Enqueue(http.StatusOK, fixture, nil)
	}

	ctx := context.Background()
	page, err := c.Vaults.ListAccount(ctx, &VaultListParams{IsFavorite: Bool(false), Search: "example"})
	if err != nil {
		t.Fatal(err)
	}
	v := page.Items[0]
	if !v.CanBeRemoved || v.UserAccountMembers[0].Email != "john.doe@example.com" || !v.UserAccountMembers[0].IsActive {
		t.Errorf("vault item = %+v", v)
	}
	if _, err := c.Vaults.ListMine(ctx, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Vaults.ListByUser(ctx, "user/1", nil); err != nil {
		t.Fatal(err)
	}

	reqs := tr.Snapshot()
	if reqs[0].Path != "/v1/ecm/vaults/accounts" || reqs[1].Path != "/v1/ecm/vaults/user-accounts" {
		t.Errorf("paths = %s, %s", reqs[0].Path, reqs[1].Path)
	}
	if got := queryOf(t, reqs[0].Query).Get("IsFavorite"); got != "false" {
		t.Errorf("IsFavorite = %q", got)
	}
	if reqs[2].Path != "/v1/ecm/vaults/user-accounts/user%2F1" {
		t.Errorf("ListByUser path = %s (id must be escaped)", reqs[2].Path)
	}
}

func TestVaultContents(t *testing.T) {
	const fixture = `{
	  "vault": {"id": "v9", "isSandbox": true, "name": "Example Vault", "type": "UserAccount", "canBeRemoved": false, "userAccountMembers": []},
	  "items": [
	    {"type": "Envelope", "envelope": {"id": "e1", "status": "Published", "name": "Contract", "canBeRemoved": false, "signers": [{"name": "John", "email": "j@x.com"}]}, "template": null},
	    {"type": "Template", "envelope": null, "template": {"id": "t1", "name": "NDA Template"}}
	  ],
	  "pagination": {"totalItems": 2, "pageSize": 10, "pageNumber": 1, "pageItems": 2}
	}`
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusOK, fixture, nil)

	contents, err := c.Vaults.Contents(context.Background(), "v9", &VaultContentsParams{
		ExcludeTemplates: Bool(false),
		Search:           "con",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !contents.Vault.IsSandbox || len(contents.Items) != 2 {
		t.Fatalf("contents = %+v", contents)
	}
	env := contents.Items[0]
	if env.Type != VaultItemTypeEnvelope || env.Envelope == nil || env.Envelope.Status != EnvelopeStatusPublished || env.Template != nil {
		t.Errorf("envelope item = %+v", env)
	}
	tpl := contents.Items[1]
	if tpl.Type != VaultItemTypeTemplate || tpl.Template == nil || tpl.Template.Name != "NDA Template" {
		t.Errorf("template item = %+v", tpl)
	}

	req := tr.Snapshot()[0]
	if req.Path != "/v1/ecm/vaults/v9/list" {
		t.Errorf("path = %s", req.Path)
	}
	q := queryOf(t, req.Query)
	if q.Get("ExcludeTemplates") != "false" || q.Get("Search") != "con" {
		t.Errorf("query = %v", q)
	}
}

func TestVaultContentsAutoPaging(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusOK, `{"vault":{"id":"v9"},"items":[{"type":"Template","template":{"id":"t1"}}],"pagination":{"totalItems":2,"pageSize":1,"pageNumber":1,"pageItems":1}}`, nil)
	tr.Enqueue(http.StatusOK, `{"vault":{"id":"v9"},"items":[{"type":"Template","template":{"id":"t2"}}],"pagination":{"totalItems":2,"pageSize":1,"pageNumber":2,"pageItems":1}}`, nil)

	var ids []string
	for item, err := range c.Vaults.ContentsAutoPaging(context.Background(), "v9", &VaultContentsParams{ListParams: ListParams{PageSize: 1}}) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, item.Template.ID)
	}
	if len(ids) != 2 || ids[1] != "t2" {
		t.Errorf("ids = %v", ids)
	}
	if got := queryOf(t, tr.Snapshot()[1].Query).Get("PageNumber"); got != "2" {
		t.Errorf("second request PageNumber = %q", got)
	}
}

func TestVaultFavoriteUnfavorite(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusNoContent, "", nil)
	tr.Enqueue(http.StatusNoContent, "", nil)

	ctx := context.Background()
	if err := c.Vaults.Favorite(ctx, "v1"); err != nil {
		t.Fatal(err)
	}
	if err := c.Vaults.Unfavorite(ctx, "v1"); err != nil {
		t.Fatal(err)
	}
	reqs := tr.Snapshot()
	if reqs[0].Path != "/v1/ecm/vaults/v1/favorite" || reqs[1].Path != "/v1/ecm/vaults/v1/unfavorite" {
		t.Errorf("paths = %s, %s", reqs[0].Path, reqs[1].Path)
	}
}
