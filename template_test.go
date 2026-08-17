package signater

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestTemplateCreate(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusCreated, `{"id":"tpl-1"}`, nil)

	id, err := c.Templates.Create(context.Background(), &TemplateParams{
		VaultID:     "1900bc26-7e01-40a1-aa83-d2abc112401c",
		Name:        "Contract Template",
		Description: String("A template description for contracts between parties."),
		Payload:     String(`{ "key": "value" }`),
		Fields: []TemplateFieldParams{
			{Name: "First Name", Alias: "first_name", IsRequired: true, Type: TemplateFieldText},
			{Name: "Age", Alias: "age", Type: TemplateFieldNumber},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != "tpl-1" {
		t.Errorf("id = %q", id)
	}

	req := tr.Snapshot()[0]
	if req.Method != http.MethodPost || req.Path != "/v1/ecm/templates" {
		t.Errorf("request = %s %s", req.Method, req.Path)
	}
	var body struct {
		VaultID string `json:"vaultId"`
		Fields  []map[string]any
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		t.Fatal(err)
	}
	if body.VaultID == "" || len(body.Fields) != 2 {
		t.Errorf("body = %+v", body)
	}
	if body.Fields[0]["alias"] != "first_name" || body.Fields[0]["type"] != "Text" {
		t.Errorf("field 0 = %v", body.Fields[0])
	}
	if _, present := body.Fields[0]["id"]; present {
		t.Error("nil field id must be omitted on create")
	}
}

func TestTemplateGet(t *testing.T) {
	const fixture = `{
	  "id": "tpl-1",
	  "vaultId": "v-1",
	  "vaultName": "Contracts",
	  "name": "Contract Template",
	  "description": "desc",
	  "fields": [{
	    "id": "f-1",
	    "alias": "first_name",
	    "isRequired": true,
	    "type": "Dropdown",
	    "schema": {
	      "label": "First name",
	      "options": ["a", "b"],
	      "constraints": {"minLength": 2, "maxLength": 10, "allowNegative": null, "profile": "None"}
	    },
	    "index": 0
	  }]
	}`
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusOK, fixture, nil)

	tpl, err := c.Templates.Get(context.Background(), "tpl-1")
	if err != nil {
		t.Fatal(err)
	}
	if tpl.VaultName != "Contracts" || len(tpl.Fields) != 1 {
		t.Fatalf("template = %+v", tpl)
	}
	f := tpl.Fields[0]
	if f.Type != TemplateFieldDisplayDropdown || !f.IsRequired || f.Schema == nil {
		t.Errorf("field = %+v", f)
	}
	cons := f.Schema.Constraints
	if cons == nil || cons.MinLength == nil || *cons.MinLength != 2 || cons.AllowNegative != nil {
		t.Errorf("constraints = %+v", cons)
	}
	if req := tr.Snapshot()[0]; req.Path != "/v1/ecm/templates/tpl-1" {
		t.Errorf("path = %s", req.Path)
	}
}

func TestTemplateLifecycleVerbs(t *testing.T) {
	c, tr := newTestClient(t)
	for range 5 {
		tr.Enqueue(http.StatusNoContent, "", nil)
	}

	ctx := context.Background()
	if err := c.Templates.Update(ctx, "tpl-1", &TemplateParams{VaultID: "v-1", Name: "New"}); err != nil {
		t.Fatal(err)
	}
	if err := c.Templates.Rename(ctx, "tpl-1", &TemplateRenameParams{Name: "New name", Description: String("New description")}); err != nil {
		t.Fatal(err)
	}
	if err := c.Templates.Move(ctx, "tpl-1", "v-2"); err != nil {
		t.Fatal(err)
	}
	if err := c.Templates.Restore(ctx, "tpl-1", "v-3"); err != nil {
		t.Fatal(err)
	}
	if err := c.Templates.Remove(ctx, "tpl-1"); err != nil {
		t.Fatal(err)
	}

	reqs := tr.Snapshot()
	wantReqs := []struct{ method, path string }{
		{http.MethodPut, "/v1/ecm/templates/tpl-1"},
		{http.MethodPost, "/v1/ecm/templates/tpl-1/rename"},
		{http.MethodPost, "/v1/ecm/templates/tpl-1/move"},
		{http.MethodPost, "/v1/ecm/templates/tpl-1/restore"},
		{http.MethodDelete, "/v1/ecm/templates/tpl-1"},
	}
	for i, want := range wantReqs {
		if reqs[i].Method != want.method || reqs[i].Path != want.path {
			t.Errorf("request %d = %s %s, want %s %s", i, reqs[i].Method, reqs[i].Path, want.method, want.path)
		}
	}
	var updateBody, renameBody, moveBody, restoreBody map[string]any
	for i, dst := range []*map[string]any{&updateBody, &renameBody, &moveBody, &restoreBody} {
		if err := json.Unmarshal(reqs[i].Body, dst); err != nil {
			t.Fatalf("request %d body: %v", i, err)
		}
	}
	if updateBody["vaultId"] != "v-1" || updateBody["name"] != "New" {
		t.Errorf("update body = %v", updateBody)
	}
	if renameBody["name"] != "New name" || renameBody["description"] != "New description" {
		t.Errorf("rename body = %v", renameBody)
	}
	if moveBody["vaultId"] != "v-2" {
		t.Errorf("move body = %v", moveBody)
	}
	// Restore requires the destination vault in the body per the OpenAPI
	// fragment (RestoreTemplateApiRequest, required: [vaultId]).
	if restoreBody["vaultId"] != "v-3" {
		t.Errorf("restore body = %v", restoreBody)
	}
}
