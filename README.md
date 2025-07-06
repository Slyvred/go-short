# Go Short

A minimalist URL shortener I built in Go, in order to learn the langage.

## Quick Start

```bash
# Clone the repository
git clone https://github.com/slyvred/go-short.git
cd go-short

# Build and run
go build -o go-short
./go-short
```

## API Endpoints

### Shorten URL
```json
POST /shorten
{
  "url": "https://example.com/very/long/url",
}
```

#### Response
```json
{
  "original": "https://example.com/very/long/url",
  "shortened": "cbk78yyk"
}
```

**Note**: URLs that have not been accessed for 60 days are automatically removed from the system.

### Redirect
```
GET /{short_code}
```

### Get Stats
```
GET /{short_code}/stats
```

#### Response
```json
{
  "accessCount": 0,
  "lastAccessed": "2025-07-06T09:52:38.212Z",
  "original": "https://example.com/very/long/url"
}
```

## Configuration

Set environment variables:
- `MONGO_URI`: MongoDB Connection string (can be set in a .env if you want)

## License

GPL-3.0
