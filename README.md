# Image Processing Application

A full-stack application for asynchronous image processing. The project provides a React-based frontend and a robust Go-based backend (API + Background Worker) that leverages RabbitMQ for message queuing, PostgreSQL for metadata storage, and Cloudflare R2 (or any S3-compatible storage) for object storage.

## Architecture

The system consists of the following components:

- **Frontend (`/frontend`)**: A modern web interface built with React 19, TypeScript, and Vite.
- **Backend API (`/backend/cmd/api`)**: A Go REST API that handles client requests, stores metadata in PostgreSQL, uploads initial images to S3, and pushes transformation jobs to RabbitMQ.
- **Backend Worker (`/backend/cmd/worker`)**: A Go background worker that consumes jobs from RabbitMQ, downloads images from S3, processes them using `libvips` (fast image processing library), and uploads the transformed results back to S3.
- **Database**: PostgreSQL for persisting image records and transformation statuses.
- **Message Broker**: RabbitMQ for reliable asynchronous task queuing between the API and Worker.

## Technologies Used

- **Frontend**: React 19, TypeScript, Vite
- **Backend**: Go 1.22
- **Image Processing**: `libvips` (via CGO)
- **Infrastructure/Services**: 
  - PostgreSQL
  - RabbitMQ
  - S3-compatible Storage (e.g., Cloudflare R2, AWS S3)
  - Docker & Docker Compose

## Prerequisites

To run this project locally, you will need:
- [Docker](https://www.docker.com/) and Docker Compose
- [Node.js](https://nodejs.org/) (for local frontend development)
- [Go 1.22+](https://golang.org/) (optional, if you want to run the backend without Docker; requires `libvips` installed locally)

## Local Development

The easiest way to run the entire backend infrastructure (Database, RabbitMQ, API, and Worker) is using Docker Compose.

### 1. Environment Setup

Create a `.env` file in the root directory (or use `.env.example` in the `backend` folder as a reference) to provide your S3/R2 credentials and other configuration values.

```env
# Example .env configuration
DATABASE_URL=postgres://user:password@postgres:5432/image_processing?sslmode=disable
RABBITMQ_URL=amqp://user:password@rabbitmq:5672/
S3_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
S3_ACCESS_KEY=your_access_key
S3_SECRET_KEY=your_secret_key
S3_BUCKET=your_bucket_name
```

### 2. Start the Backend Infrastructure

From the root of the project, run:

```bash
docker-compose up --build
```

This will spin up:
- **PostgreSQL** on port `5432`
- **RabbitMQ** (with Management UI) on ports `5672` and `15672`
- **API Service** on port `8080`
- **Worker Service** (runs in the background)

### 3. Start the Frontend

In a new terminal window, navigate to the `frontend` directory, install dependencies, and start the Vite development server:

```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at `http://localhost:5173` (or the port Vite provides).

## Deployment

This project is Docker-ready and can be easily deployed to container orchestration platforms or PaaS providers like **Render**:

1. **API**: Deploy as a Web Service using `Dockerfile` (default command).
2. **Worker**: Deploy as a Background Worker using the same `Dockerfile` by overriding the startup command to `/bin/worker`.
3. **Services**: Connect to a managed PostgreSQL database and a managed RabbitMQ instance (e.g., CloudAMQP).
