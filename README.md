# Linear TUI

A terminal user interface (TUI) application for interacting with Linear issues and projects.

## Prerequisites

- Go 1.19 or later
- A Linear account with API access

## Setup

### 1. Get your Linear API Key

1. Go to [Linear Settings > API](https://linear.app/settings/api)
2. Click "Create API Key"
3. Give it a descriptive name (e.g., "Linear TUI")
4. Copy the generated API key

### 2. Configure the API Key

You can provide your Linear API key in two ways:

#### Option A: Environment Variable (Recommended)
```bash
export LINEAR_API_KEY="your_api_key_here"
```

#### Option B: Configuration File
The application will also check for API keys in the configuration file (details TBD).

### 3. Build and Run

```bash
# Build the application
go build -o linear-tui cmd/linear-tui/main.go

# Run the application
./linear-tui
```

## Debug Logging

To enable detailed logging for troubleshooting API requests and responses:

```bash
# Enable debug logging
export DEBUG=1

# Run with debug logging
./linear-tui
```

Debug logs will be written to `debug.log` in the current directory and include:
- API key validation process
- All HTTP requests and responses
- Rate limiting information
- Retry attempts and backoff delays
- Method-specific operation logging

## Features

- Browse Linear issues and projects
- View issue details
- Team and user information
- Real-time data fetching with retry logic
- Rate limiting compliance

## Troubleshooting

### API Key Issues
If you're having trouble with API authentication:

1. Verify your API key is correct
2. Check that the key has appropriate permissions
3. Enable debug logging to see detailed error messages:
   ```bash
   DEBUG=1 ./linear-tui
   ```

### Network Issues
If experiencing network connectivity problems:
- Debug logging will show request/response details
- The application includes automatic retry logic with exponential backoff
- Check your network connection and firewall settings

## Development

### Building from Source
```bash
git clone <repository>
cd linear-tui
go mod tidy
go build cmd/linear-tui/main.go
```

### Running Tests
```bash
go test ./...
```