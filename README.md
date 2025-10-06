# Video Streamer

A real-time video streaming application built with Go backend and React frontend. The application processes video files using FFmpeg and streams compressed JPEG frames to connected clients via WebSocket.

## Features

- Real-time video streaming through WebSocket connections
- Automatic video processing with FFmpeg at 10 FPS
- JPEG compression for bandwidth optimization
- Modern React frontend with real-time statistics
- Automatic reconnection handling
- CORS support for development
- Health monitoring endpoints

## Architecture

The application consists of two main components:

**Backend (Go)**
- FFmpeg integration for video frame extraction
- WebSocket server for real-time communication
- JPEG compression pipeline
- Health check endpoints

**Frontend (React + TypeScript)**
- WebSocket client with automatic reconnection
- Real-time video display
- Connection status monitoring
- Frame statistics dashboard

## Prerequisites

- Go 1.24.5 or higher
- Node.js 18 or higher
- FFmpeg installed and available in PATH
- A video file (MP4 recommended)

## Installation

1. Install backend dependencies:
   ```bash
   cd api
   go mod tidy
   ```

2. Install frontend dependencies:
   ```bash
   cd frontend
   npm install
   ```

3. Place a video file in one of these locations:
   - `video/sample/FireForce-S1E3-360P.mp4`
   - `sample.mp4`
   - `test.mp4`
   - `video.mp4`

## Usage

### Development

Start the backend server:
```bash
cd api
go run main.go
```

Start the frontend development server:
```bash
cd frontend
npm run dev
```

Access the application at http://localhost:5173

### Production

Build the backend:
```bash
cd api
go build -o video-streamer main.go
```

Build the frontend:
```bash
cd frontend
npm run build
npm run preview
```

## API Endpoints

- `ws://localhost:8080/ws` - WebSocket streaming endpoint
- `http://localhost:8080/health` - Health check endpoint

## Configuration

### Backend Settings

- Server port: 8080
- Frame rate: 10 FPS
- JPEG quality: 30
- Buffer size: 30 frames

Modify these settings in the respective source files:
- Frame rate: `api/video/read.go`
- JPEG quality: `api/video/quant.go`
- Server port: `api/main.go`

### Frontend Settings

- Development port: 5173 (Vite default)
- WebSocket endpoint: `ws://localhost:8080/ws`
- Reconnection delay: 3 seconds

## Project Structure

```
├── api/                    # Go backend
│   ├── main.go            # Server entry point
│   ├── streaming/         # WebSocket handling
│   ├── video/             # Video processing
│   └── utils/             # Utilities
├── frontend/              # React frontend
│   ├── src/
│   │   ├── App.tsx       # Main application
│   │   └── main.tsx      # Entry point
│   ├── package.json      # Dependencies
│   └── vite.config.ts    # Build configuration
└── WARP.md               # Development guidance
```

## Development Commands

### Backend
```bash
go run main.go        # Run development server
go build main.go      # Build for production
go test ./...         # Run tests
go fmt ./...          # Format code
```

### Frontend
```bash
npm run dev           # Development server
npm run build         # Production build
npm run preview       # Preview production build
npm run lint          # Run linter
```

## Troubleshooting

**Video file not found**: Ensure your video file is placed in one of the supported locations.

**FFmpeg not found**: Install FFmpeg and ensure it's available in your system PATH.

**WebSocket connection failed**: Verify the backend server is running on port 8080.

**Build errors**: Clear node_modules and reinstall dependencies with `npm install`.

## License

This project is licensed under the MIT License.