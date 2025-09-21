# File Storage API

A REST API for uploading and managing files using Express.js, TypeScript, and UploadThing.

## Setup

1. Install dependencies:
```bash
pnpm install
```

2. Copy environment variables:
```bash
cp .env.example .env
```

3. Configure your UploadThing credentials in `.env`

4. Start development server:
```bash
pnpm dev
```

## API Documentation

Interactive API documentation is available at `/api/v1/swagger` when the server is running.

### API Endpoints

#### Upload Files
- **POST** `/api/v1/files/upload`
- Content-Type: `multipart/form-data`
- Form fields:
  - `files`: One or more files to upload
  - `domain`: Domain for all uploaded files (string)
- Response: `{ files: [{ id: string, name: string, url: string }] }`

#### Delete Files
- **DELETE** `/api/v1/files/delete`
- Content-Type: `application/json`
- Body: `{ fileIds: string[] }`
- Response: `204 No Content`

## Scripts

- `pnpm dev` - Start development server
- `pnpm build` - Build for production
- `pnpm start` - Start production server
- `pnpm type-check` - Check TypeScript types