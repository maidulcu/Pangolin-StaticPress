# StaticPress - Project Plan

## Overview
StaticPress is a CLI tool to export WordPress sites to static HTML for deployment to S3, Netlify, or other static hosting providers.

## Architecture

### Components
1. **CLI (Go)** - Main engine: `init`, `export`, `deploy`
2. **Connector (PHP)** - WordPress Plugin for Auth & Webhooks
3. **Dashboard (Go/HTMX)** - Optional GUI for Pro users

### Technology Stack
- **Language:** Go 1.21+
- **CLI Framework:** Cobra
- **HTML Parsing:** goquery
- **Config:** Viper
- **AWS SDK:** AWS SDK v2 for S3
- **WP Plugin:** PHP (WordPress)

## Current Status: MVP Complete ✅

### Completed Features
- [x] CLI with Cobra (init, export, deploy commands)
- [x] Sitemap fetching (sitemap.xml / wp-sitemap.xml)
- [x] Concurrent page crawling with goroutines
- [x] Link rewriting (absolute → relative URLs)
- [x] Static HTML export to local folder
- [x] Config management with Viper (saves to ~/.staticpress/)
- [x] S3 deployment with content-type detection

### Usage
```bash
# Initialize with WordPress site
staticpress init -u https://example.com -k YOUR_API_KEY

# Export to static HTML
staticpress export -c 5 -d dist

# Deploy to S3
staticpress deploy -b my-bucket -r us-east-1
```

## Future Enhancements

### Phase 1: WordPress Plugin
- [ ] PHP WordPress plugin
- [ ] API key generation endpoint
- [ ] Auth via Bearer token
- [ ] Webhook support for auto-export

### Phase 2: Enhanced Features
- [ ] Netlify deployment support
- [ ] Image optimization
- [ ] CSS/JS asset bundling
- [ ] Incremental export (only changed pages)
- [ ] Preview mode (local server)

### Phase 3: Pro Features (Dashboard)
- [ ] Go/HTMX dashboard
- [ ] Auto-sync on content change
- [ ] CDN cache invalidation
- [ ] Multi-site support

## Security Model
- No admin access required
- Binary runs locally (not on WP server)
- API key requires only `edit_posts` capability
- Config stored encrypted locally

## Project Structure
```
├── main.go                 # Entry point
├── go.mod                  # Go dependencies
├── cmd/
│   ├── init.go             # init command
│   ├── export.go           # export command
│   ├── deploy.go           # deploy command
│   └── internal/
│       ├── config/         # Config management
│       ├── sitemap/       # Sitemap fetching
│       ├── crawler/       # Page fetching & link rewriting
│       └── exporter/      # Export & S3 upload
└── wp-plugin/              # WordPress plugin (future)
```

## Configuration
Config is saved to `~/.staticpress/staticpress.yaml`:
```yaml
site_url: "https://example.com"
api_key: "your-api-key"
s3_bucket: ""
s3_region: "us-east-1"
```

## Environment Variables (for S3 deploy)
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

## Roadmap
1. **MVP** (Current) - CLI export to local folder ✅
2. **Phase 1** - WordPress Plugin for auth
3. **Phase 2** - Enhanced features (Netlify, image optimization)
4. **Phase 3** - Dashboard & Pro features
