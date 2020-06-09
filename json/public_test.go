package json

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

func TestParse_nonObject(t *testing.T) {
	src := `true`
	file, diags := Parse([]byte(src), "", hcl.Pos{Byte: 0, Line: 1, Column: 1})
	if len(diags) != 1 {
		t.Errorf("got %d diagnostics; want 1", len(diags))
	}
	if file == nil {
		t.Errorf("got nil File; want actual file")
	}
	if file.Body == nil {
		t.Fatalf("got nil Body; want actual body")
	}
	if file.Body.(*body).val == nil {
		t.Errorf("got nil Body object; want placeholder object")
	}
}

func TestParseTemplate(t *testing.T) {
	src := `{"greeting": "hello ${\"world\"}"}`
	file, diags := Parse([]byte(src), "", hcl.Pos{Byte: 0, Line: 1, Column: 1})
	if len(diags) != 0 {
		t.Errorf("got %d diagnostics on parse; want 0", len(diags))
		for _, diag := range diags {
			t.Logf("- %s", diag.Error())
		}
	}
	if file == nil {
		t.Errorf("got nil File; want actual file")
	}
	if file.Body == nil {
		t.Fatalf("got nil Body; want actual body")
	}
	attrs, diags := file.Body.JustAttributes()
	if len(diags) != 0 {
		t.Errorf("got %d diagnostics on decode; want 0", len(diags))
		for _, diag := range diags {
			t.Logf("- %s", diag.Error())
		}
	}

	val, diags := attrs["greeting"].Expr.Value(&hcl.EvalContext{})
	if len(diags) != 0 {
		t.Errorf("got %d diagnostics on eval; want 0", len(diags))
		for _, diag := range diags {
			t.Logf("- %s", diag.Error())
		}
	}

	if !val.RawEquals(cty.StringVal("hello world")) {
		t.Errorf("wrong result %#v; want %#v", val, cty.StringVal("hello world"))
	}
}

func TestParseTemplateUnwrap(t *testing.T) {
	src := `{"greeting": "${true}"}`
	file, diags := Parse([]byte(src), "", hcl.Pos{Byte: 0, Line: 1, Column: 1})
	if len(diags) != 0 {
		t.Errorf("got %d diagnostics on parse; want 0", len(diags))
		for _, diag := range diags {
			t.Logf("- %s", diag.Error())
		}
	}
	if file == nil {
		t.Errorf("got nil File; want actual file")
	}
	if file.Body == nil {
		t.Fatalf("got nil Body; want actual body")
	}
	attrs, diags := file.Body.JustAttributes()
	if len(diags) != 0 {
		t.Errorf("got %d diagnostics on decode; want 0", len(diags))
		for _, diag := range diags {
			t.Logf("- %s", diag.Error())
		}
	}

	val, diags := attrs["greeting"].Expr.Value(&hcl.EvalContext{})
	if len(diags) != 0 {
		t.Errorf("got %d diagnostics on eval; want 0", len(diags))
		for _, diag := range diags {
			t.Logf("- %s", diag.Error())
		}
	}

	if !val.RawEquals(cty.True) {
		t.Errorf("wrong result %#v; want %#v", val, cty.True)
	}
}

func TestParse_malformed(t *testing.T) {
	src := `{
  "http_proxy_url: "http://xxxxxx",
}`
	file, diags := Parse([]byte(src), "", hcl.Pos{Byte: 0, Line: 1, Column: 1})
	if got, want := len(diags), 2; got != want {
		t.Errorf("got %d diagnostics; want %d", got, want)
	}
	if err, want := diags.Error(), `Missing property value colon`; !strings.Contains(err, want) {
		t.Errorf("diags are %q, but should contain %q", err, want)
	}
	if file == nil {
		t.Errorf("got nil File; want actual file")
	}
}
