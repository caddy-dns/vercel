package vercel

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/vercel"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *vercel.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.vercel",
		New: func() caddy.Module { return &Provider{new(vercel.Provider)} },
	}
}

// Before using the provider config, resolve placeholders in the API token.
// Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.AuthAPIToken = repl.ReplaceAll(p.Provider.AuthAPIToken, "")
	p.Provider.TeamId = repl.ReplaceAll(p.Provider.TeamId, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// vercel [<api_token>] {
//     api_token <api_token>
// }
//
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		// if d.NextArg() {
		// 	p.Provider.AuthAPIToken = d.Val()
		// 	p.Provider.TeamId = d.Val()
		// }
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "auth_api_token":
				if p.Provider.AuthAPIToken != "" {
					return d.Err("API token already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.AuthAPIToken = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "team_id":
				if p.Provider.TeamId != "" {
					return d.Err("team ID already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.TeamId = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.AuthAPIToken == "" {
		return d.Err("missing API token")
	}
	if p.Provider.TeamId == "" {
		return d.Err("missing team ID")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
