package vercel

import (
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	vercelLib "github.com/libdns/vercel"
)

func TestUnmarshalCaddyfile_InlineToken(t *testing.T) {
	input := `vercel my-secret-token`
	d := caddyfile.NewTestDispenser(input)
	p := &Provider{new(vercelLib.Provider)}
	err := p.UnmarshalCaddyfile(d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Provider.AuthAPIToken != "my-secret-token" {
		t.Fatalf("expected 'my-secret-token', got '%s'", p.Provider.AuthAPIToken)
	}
}

func TestUnmarshalCaddyfile_BlockToken(t *testing.T) {
	input := `vercel {
	api_token my-block-token
}`
	d := caddyfile.NewTestDispenser(input)
	p := &Provider{new(vercelLib.Provider)}
	err := p.UnmarshalCaddyfile(d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Provider.AuthAPIToken != "my-block-token" {
		t.Fatalf("expected 'my-block-token', got '%s'", p.Provider.AuthAPIToken)
	}
}

func TestUnmarshalCaddyfile_MissingToken(t *testing.T) {
	input := `vercel`
	d := caddyfile.NewTestDispenser(input)
	p := &Provider{new(vercelLib.Provider)}
	err := p.UnmarshalCaddyfile(d)
	if err == nil {
		t.Fatal("expected error for missing token, got nil")
	}
}

func TestUnmarshalCaddyfile_BlockMissingValue(t *testing.T) {
	input := `vercel {
	api_token
}`
	d := caddyfile.NewTestDispenser(input)
	p := &Provider{new(vercelLib.Provider)}
	err := p.UnmarshalCaddyfile(d)
	if err == nil {
		t.Fatal("expected error for missing token value, got nil")
	}
}

func TestUnmarshalCaddyfile_DuplicateToken(t *testing.T) {
	input := `vercel my-token {
	api_token another-token
}`
	d := caddyfile.NewTestDispenser(input)
	p := &Provider{new(vercelLib.Provider)}
	err := p.UnmarshalCaddyfile(d)
	if err == nil {
		t.Fatal("expected error for duplicate token, got nil")
	}
}

func TestUnmarshalCaddyfile_TooManyArgs(t *testing.T) {
	input := `vercel token1 token2`
	d := caddyfile.NewTestDispenser(input)
	p := &Provider{new(vercelLib.Provider)}
	err := p.UnmarshalCaddyfile(d)
	if err == nil {
		t.Fatal("expected error for too many args, got nil")
	}
}

func TestUnmarshalCaddyfile_UnrecognizedSubdirective(t *testing.T) {
	input := `vercel {
	unknown_directive value
}`
	d := caddyfile.NewTestDispenser(input)
	p := &Provider{new(vercelLib.Provider)}
	err := p.UnmarshalCaddyfile(d)
	if err == nil {
		t.Fatal("expected error for unrecognized subdirective, got nil")
	}
}
