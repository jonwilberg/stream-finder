# Stream Finder

> ⚠️ **Work in Progress** - This project is currently under active development. Features may be incomplete or subject to change.

A web application to help users find where to stream their favorite movies and TV shows across different platforms.

## Features

- 🔍 Search for movies and TV shows
- 📺 Find streaming availability across platforms
- 🎬 Clean, modern interface with responsive design
- ⚡ Fast search powered by Elasticsearch

## Tech Stack

**Frontend:**
- React 19 with TypeScript
- React Router 7
- Tailwind CSS 4
- Vite

**Backend:**
- Go 1.22
- Elasticsearch
- Google Cloud Firestore
- Netflix API integration

**Infrastructure:**
- Terraform for infrastructure as code
- Docker support

## Getting Started

### Prerequisites

- Node.js 18+
- Go 1.22+
- Elasticsearch
- Docker (optional)

### Installation

1. **Frontend setup**
   ```bash
   cd frontend
   npm install
   ```

2. **Backend setup**
   ```bash
   cd backend
   go mod download
   ```

3. **Environment variables**
   Set up your environment variables for:
   - Elasticsearch connection
   - Google Cloud credentials
   - Netflix API keys

### Running the Application

**Development:**

Frontend:
```bash
cd frontend
npm run dev
```

Backend:
```bash
cd backend
go run cmd/titles/main.go
```

**Production:**
```bash
# Frontend
cd frontend
npm run build
npm start

# Backend
cd backend
go build -o stream-finder cmd/titles/main.go
./stream-finder
```

## Project Structure

```
stream-finder/
├── frontend/              # React frontend application
│   ├── app/
│   │   ├── components/    # React components
│   │   ├── routes/        # Route handlers
│   │   └── welcome/       # Welcome page
│   ├── public/            # Static assets
│   └── package.json
├── backend/               # Go backend application
│   ├── cmd/titles/        # Main application entry point
│   ├── internal/          # Internal packages
│   │   ├── repos/         # Repository implementations
│   │   └── titles/        # Title-related logic
│   └── pkg/               # Shared packages
├── tf/                    # Terraform infrastructure
└── README.md
```