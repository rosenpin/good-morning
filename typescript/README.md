# Good Morning Worker - TypeScript Migration

A TypeScript migration of the Go good-morning server to Cloudflare Workers.

## Migration Complete

✅ **Architecture**: Maintains same design as Go version  
✅ **Config**: Environment variables + secrets  
✅ **Caching**: KV storage instead of files  
✅ **Daily Quotas**: Durable Objects instead of in-memory counters  
✅ **Periodic Resets**: Cron triggers instead of goroutines  

## Deployment Steps

### 1. Install Dependencies
```bash
cd typescript
npm install
```

### 2. Configure Secrets
```bash
# Set your API credentials as secrets
wrangler secret put API_KEY
wrangler secret put API_CX
```

### 3. Create KV Namespace
```bash
# Create KV namespace for caching
wrangler kv:namespace create "CACHE_KV"
wrangler kv:namespace create "CACHE_KV" --preview

# Update the IDs in wrangler.toml
```

### 4. Deploy
```bash
wrangler deploy
```

## Environment Variables

Set these in `wrangler.toml` (already configured):
- `IMAGE_SIZE`: Image size filter
- `IMAGE_RIGHTS`: Usage rights filter  
- `IMAGE_LIFESPAN`: Cache duration in hours
- `SEARCH_RANDOMNESS`: Random search offset
- `SEARCH_FEATURES`: Comma-separated search features
- `SEARCH_SAFE`: Safe search setting
- `SEARCH_BASE_QUERY`: Base search query
- `MAX_DAILY_RELOAD`: Daily reload limit

## API Endpoints

- `GET /`: Returns cached daily image
- `POST /change`: Forces reload (respects daily quota)

## Key Features

- **Daily Quota System**: Uses Durable Objects for persistent counters
- **Smart Caching**: KV-based with configurable expiration
- **Auto Reset**: Cron triggers reset quota at midnight UTC
- **Error Handling**: Maintains same error responses as Go version

## File Structure

```
typescript/
├── src/
│   ├── caching/     # KV-based cache (replaces file cache)
│   ├── config/      # Environment-based config loading
│   ├── provider/    # Core business logic
│   ├── querier/     # HTTP requests to Google API
│   ├── quota/       # Durable Object quota management  
│   ├── result/      # Google API response parsing
│   ├── url/         # Google Images URL construction
│   └── index.ts     # Workers fetch handler
├── package.json
├── tsconfig.json
└── wrangler.toml
```