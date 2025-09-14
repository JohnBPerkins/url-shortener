# URL Shortener Frontend

A modern Next.js frontend for the URL Shortener service.

## Quick Start

1. **Install dependencies**:
   ```bash
   npm install
   ```

2. **Start development server**:
   ```bash
   npm run dev
   ```

3. **Open your browser** and navigate to `http://localhost:3001`

## Environment Setup

- Copy `.env.local` and configure the API endpoint:
  ```env
  NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
  ```

## Available Scripts

- `npm run dev` - Start development server with hot reloading
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint

## Features

- ✅ Modern React 19 + Next.js 15
- ✅ TypeScript for type safety
- ✅ Tailwind CSS for styling
- ✅ Responsive design
- ✅ Real-time metrics dashboard
- ✅ Glass morphism UI design
- ✅ CORS-enabled API communication

## API Integration

The frontend communicates with the Go backend via:
- `POST /api/shorten` - Create shortened URLs
- `GET /{code}` - Redirect to original URLs (handled by backend)

## Deployment

### Docker
```bash
docker build -t url-shortener-frontend .
docker run -p 3001:3001 url-shortener-frontend
```

### Vercel/Netlify
The app is ready for deployment to modern hosting platforms. Configure the `NEXT_PUBLIC_API_BASE_URL` environment variable to point to your production API.

---

This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).