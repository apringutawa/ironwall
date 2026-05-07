# IronWall Frontend

Next.js dashboard for IronWall security monitoring system.

## Setup

```bash
cd frontend
npm install
cp .env.example .env.local
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Pages

- `/` - Dashboard overview
- `/events` - Security events timeline
- `/modules` - Module management
- `/firewall` - Firewall rules management
- `/settings` - Alert configuration

## Tech Stack

- Next.js 15 (App Router)
- React 19
- TypeScript
- Tailwind CSS 4
- Axios
- Zustand

## Environment Variables

```env
NEXT_PUBLIC_API_URL=http://localhost:8001
```

## Build

```bash
npm run build
npm start
```

## Docker

```bash
docker build -t ironwall-frontend .
docker run -p 3000:3000 ironwall-frontend
```
