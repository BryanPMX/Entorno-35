# CI/CD and Deployment

This document describes how commits trigger automatic updates for the Vercel frontend and the self-hosted backend (Docker on Portainer), and how to investigate Cloudflare. It covers required GitHub secrets, Portainer stack setup, and a step-by-step Cloudflare and Portainer checklist.

**Production branch**: `develop`. Vercel production deploys and the backend image tag used in production both come from `develop`.

**Last updated**: 2026-01-31

---

## Overview

| Component | Hosting | Trigger | What happens on push |
|-----------|---------|---------|----------------------|
| Frontend  | Vercel  | Git push to `develop` (or enabled branch) | Vercel builds and deploys from `web/frontend/` |
| Backend   | Self-hosted (Docker on Portainer) | Git push to `develop` when backend paths change | GitHub Actions builds image (tag `develop`), pushes to Docker Hub, then calls Portainer API to update the stack (pull new image and restart) |

DNS and optional proxy/WAF are often handled by Cloudflare. See "Investigating your Cloudflare setup" below to see how your domain and API are configured.

## Vercel (Frontend)

- **Connection**: Link the Vercel project to this Git repository. In Vercel project settings, set **Root Directory** to `web/frontend` so builds use the Next.js app there.
- **Auto-deploy**: Pushes to branches with deployment enabled (see `web/frontend/vercel.json` under `git.deploymentEnabled`) trigger a new build and deploy. No GitHub Action is required; Vercel uses its own Git integration.
- **Branches**: Production is `develop`. `vercel.json` enables deployment for `develop` and `main`; set Vercel's production branch to `develop` in project settings.
- **Environment**: In Vercel, set `NEXT_PUBLIC_API_URL` to your backend API URL (e.g. `https://api.entorno35.com`). That URL must be reachable from the browser (and allowed in CSP in `vercel.json` if you use strict CSP).

## GitHub Actions (Backend image and Portainer)

Workflow file: `.github/workflows/docker-build.yml`.

- **When it runs**: On push to `develop` or `main` when any of these paths change: `cmd/**`, `internal/**`, `pkg/**`, `go.mod`, `go.sum`, `Dockerfile`, `.dockerignore`, `fonts/**`. Also on `workflow_dispatch` for manual runs. For production you push to `develop`; the workflow builds and pushes tag `develop` and updates Portainer if secrets are set.
- **Steps**:
  1. Checkout, set up Docker Buildx.
  2. Log in to Docker Hub (secrets).
  3. Build the image from the repo root (Dockerfile) and push to Docker Hub with tags: branch name, `{branch}-{sha}`, and `latest` (only on the default branch).
  4. If Portainer secrets are set: PUT to Portainer API to update the stack with `docker-compose.prod.yml` and `pullImage: true`, so Portainer pulls the new image and restarts the stack.

### Required GitHub secrets

| Secret | Purpose |
|--------|---------|
| `DOCKERHUB_USERNAME` | Docker Hub login (used for image name and login). |
| `DOCKERHUB_TOKEN` | Docker Hub password or access token (recommended). |
| `PORTAINER_URL` | Portainer base URL (e.g. `https://portainer.example.com`), no trailing slash. |
| `PORTAINER_API_KEY` | Portainer API key (or token) with permission to update stacks. |
| `PORTAINER_STACK_ID` | Numeric stack ID to update (from Portainer UI or API). |
| `PORTAINER_ENDPOINT_ID` | (Optional) Endpoint ID; default is `1`. |

If `PORTAINER_URL`, `PORTAINER_API_KEY`, or `PORTAINER_STACK_ID` are missing, the workflow still builds and pushes the image but skips the Portainer update step.

## Portainer (Self-hosted backend)

- **Stack**: Create a stack from `docker-compose.prod.yml` (or paste the same content). Set the stack’s environment variables (e.g. `DB_PASSWORD`, `JWT_SECRET`, `CORS_ORIGIN`, `DOCKERHUB_USERNAME`, optionally `IMAGE_TAG`).
- **Image tag**: The compose file uses `${DOCKERHUB_USERNAME:-brpmx}/entorno35-backend:${IMAGE_TAG:-develop}`. For production (develop branch) leave `IMAGE_TAG` unset or set to `develop`. The workflow pushes tag `develop` on pushes to develop.
- **After CI runs**: When the workflow updates the stack via API, Portainer receives the same `docker-compose.prod.yml` content and `pullImage: true`, so it pulls the updated image and restarts the API container. Database and Redis are unchanged unless you change the compose file.

## Portainer stack step-by-step (entorno35)

Use this to verify and align your Portainer stack with the repo and with your frontend (www.entorno35.com).

### 1. Stack and services

- **Stack**: One stack named e.g. `entorno35` with compose content from `docker-compose.prod.yml`.
- **Services**: Three services from that compose:
  - **api** – Go backend, image `brpmx/entorno35-backend:develop`, port 8080.
  - **postgres** – PostgreSQL 15, internal only (no host port needed unless you access DB from outside).
  - **redis** – Redis 7, internal only.
- A fourth healthy container (e.g. `cloudflared` for the tunnel) may be in the same stack or a separate one; the repo compose only defines api, postgres, redis.

### 2. Do not put secrets in the repo

The repo’s `docker-compose.prod.yml` uses **placeholders** only: `${DB_PASSWORD}`, `${JWT_SECRET}`, `${CORS_ORIGIN}`, etc. Real values must **not** be committed. Set them in Portainer:

- **Portainer** > your stack > **Editor** (or **Stack** > **Environment variables**): define `DB_PASSWORD`, `JWT_SECRET`, `CORS_ORIGIN`, `DB_USER`, `DB_NAME`, and optionally `SMTP_*`, `DOCKERHUB_USERNAME`, `IMAGE_TAG`. Portainer injects these when deploying.

If your current stack file in Portainer has plain-text passwords, that’s only in Portainer (not in Git). For the next CI update, the workflow will send the repo’s compose (with `${VAR}` only); Portainer keeps your existing env vars, so your secrets stay in Portainer.

### 3. CORS must match the frontend URL

Your frontend is at **https://www.entorno35.com** (Vercel). The browser sends that as the `Origin` header. The backend must allow it.

- In Portainer, set **CORS_ORIGIN** = `https://www.entorno35.com` (no trailing slash).
- If you use `https://entorno35.com` as the main URL (after redirect or root pointing to Vercel), use that instead. Origin is the URL in the browser address bar.

Wrong CORS (e.g. `https://entorno35.com` when users use www) will cause login or API requests to be blocked by the browser.

### 4. Service names: postgres (not db)

The repo’s `docker-compose.prod.yml` uses service name **postgres** (not `db`) so it matches typical Portainer setups. The API has `DB_HOST=postgres` and `depends_on: postgres, redis`. When CI updates the stack, Portainer gets this same structure; no extra `db` service is created.

### 5. Optional: expose Postgres/Redis on the host

The repo compose does **not** publish `5432` or `6379` on the host, so only the `api` container can reach them. If you need to connect from the host (e.g. backups, debugging), add to the stack in Portainer:

```yaml
postgres:
  ports:
    - "127.0.0.1:5432:5432"
redis:
  ports:
    - "127.0.0.1:6379:6379"
```

Binding to `127.0.0.1` keeps them off the public network. Remove these if you don’t need host access.

### 6. Compose correctness (what to check)

Your stack compose in Portainer should have:

- **version**: `'3.8'`
- **api**: `image: brpmx/entorno35-backend:develop`, `depends_on: postgres, redis`, `ports: "8080:8080"`, `DB_HOST=postgres`, `REDIS_HOST=redis`, `CORS_ORIGIN=https://www.entorno35.com` (no trailing slash). Healthcheck with wget to `http://localhost:8080/health`.
- **postgres**: service name `postgres`, `image: postgres:15-alpine`, env `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, volume `postgres_data`. Healthcheck `pg_isready -U postgres`.
- **redis**: service name `redis`, `image: redis:7-alpine`, volume `redis_data`. Healthcheck `redis-cli ping`.
- **volumes**: `postgres_data`, `redis_data` defined at bottom.

Optional: `container_name` (e.g. `entorno35-postgres`, `entorno35-redis`) and exposing Postgres/Redis only on the host loopback (e.g. `127.0.0.1:5432:5432`, `127.0.0.1:6379:6379`) if you need host access. Do not commit real `DB_PASSWORD` or `JWT_SECRET` in the repo; set them in Portainer env vars and use `${DB_PASSWORD}` etc. in the compose if you use the repo template.

### 7. Checklist

| Item | Where | Value / action |
|------|--------|----------------|
| Stack name | Portainer | e.g. entorno35 |
| Image | api service | brpmx/entorno35-backend:develop |
| DB_HOST | api env | postgres |
| CORS_ORIGIN | api env or stack env | https://www.entorno35.com |
| DB_PASSWORD, JWT_SECRET | stack env only | Never in Git |
| All 3 containers healthy | Portainer | api, postgres, redis green |

### 8. After a push to develop

When GitHub Actions runs (backend paths or manual trigger) and Portainer secrets are set, the workflow will PUT the repo’s `docker-compose.prod.yml` to your stack with `pullImage: true`. Portainer will pull the new `develop` image and restart only the **api** container; postgres and redis keep running and keep their data.

## Cloudflare

- **Typical use**: DNS for your domain(s), and optionally a proxy (orange cloud) in front of the backend so traffic goes Cloudflare -> your server. The frontend can be Vercel-only (Vercel’s DNS) or go through Cloudflare too.
- **Backend**: The API hostname (e.g. `api.entorno35.com`) must resolve to your Ubuntu server (A or CNAME). If the record is "proxied" (orange cloud), Cloudflare sits in front and forwards to your server.
- **CORS**: Backend `CORS_ORIGIN` and Vercel’s `NEXT_PUBLIC_API_URL` must match the frontend origin users see (e.g. `https://yourdomain.com` or the Vercel URL).

## Investigating your Cloudflare setup

Use this to see how your domain and API are wired: Cloudflare DNS, proxy status, and that the backend is reachable from the server.

### 1. Cloudflare dashboard

1. Log in at [dash.cloudflare.com](https://dash.cloudflare.com).
2. Select the domain you use for Entorno35 (e.g. `entorno35.com` or the Vercel custom domain).
3. Go to **DNS** > **Records** and note:
   - **Frontend (website)**: Is there an A or CNAME for the main domain (e.g. `entorno35.com` or `www`)? What is the target? (Vercel often: `cname.vercel-dns.com` or similar; or an A record to Vercel IPs.) Is the cloud **orange** (proxied) or **grey** (DNS only)?
   - **Backend (API)**: Is there a record for the API hostname (e.g. `api.entorno35.com`)? Target is usually your Ubuntu server IP or a hostname that resolves to it. Note if the cloud is orange (proxied) or grey (DNS only).
4. **SSL/TLS**: Go to **SSL/TLS**. Note the mode (e.g. Full or Full (strict)). If the API is proxied, Cloudflare terminates HTTPS and can connect to your server on HTTP or HTTPS (your choice).

Write down: domain, frontend record + target + proxy (yes/no), API record + target + proxy (yes/no). That is your current Cloudflare setup.

### 2. From your Ubuntu server (SSH)

These checks confirm the backend is listening and that the server can reach itself; they do not prove Cloudflare or external reachability by themselves.

```bash
# Is the API listening on 8080?
sudo ss -tlnp | grep 8080
# or
sudo netstat -tlnp | grep 8080

# Local health check (backend only)
curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8080/health
# Expect 200.

# If your API is on a different hostname/IP on this server:
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health
```

If the health check returns 200, the backend container is up and reachable on the host. If not, check Portainer: stack running, API container healthy, no port conflict.

### 3. From the public internet (your laptop or a browser)

- **Frontend**: Open the site URL (e.g. `https://entorno35.com`). If it loads, DNS and (if used) Cloudflare for the frontend are working.
- **Backend**: In a browser or with curl, open `https://api.entorno35.com/health` (replace with your API hostname). You should get a 200 and a small JSON body. If you get certificate errors, the API hostname might not be in Cloudflare or the SSL mode may need adjustment. If the request times out, the API record may point to the wrong IP or a firewall may be blocking.

### 4. CORS and API URL

- In Vercel, check **Project** > **Settings** > **Environment Variables**: `NEXT_PUBLIC_API_URL` should be the exact API base URL users’ browsers will call (e.g. `https://api.entorno35.com`).
- On the backend (Portainer stack env or server env), `CORS_ORIGIN` must be the exact frontend origin (e.g. `https://entorno35.com` or your Vercel URL). No trailing slash. If the frontend is behind Cloudflare, the origin is still the URL in the browser (e.g. `https://entorno35.com`).

### 5. Summary table (fill after investigation)

| Item | Where to check | Your value |
|------|----------------|------------|
| Production branch | GitHub / Vercel | `develop` |
| Frontend URL | Browser / Vercel | e.g. `https://entorno35.com` |
| API URL | Vercel env, CSP in `vercel.json` | e.g. `https://api.entorno35.com` |
| API DNS record | Cloudflare DNS | e.g. `api` A/CNAME → server IP |
| API proxied (orange cloud) | Cloudflare DNS | Yes / No |
| Backend listening | SSH: `curl http://127.0.0.1:8080/health` | 200 |
| Backend reachable from internet | Browser: `https://api.../health` | 200 |
| CORS_ORIGIN on backend | Portainer stack env | Same as frontend URL |

Once this table is filled, you know how Cloudflare is involved (DNS only vs proxy for API and/or frontend) and can adjust SSL mode or firewall rules if needed.

## Current Cloudflare DNS (entorno35.com)

From your Cloudflare DNS records:

| Record | Name | Target | Proxy | Meaning |
|--------|------|--------|-------|--------|
| A | entorno35.com | 216.198.79.1 | DNS only | Root domain goes directly to this IP (not Vercel). |
| CNAME | www | f61e3827aa64ab5d.vercel-dns-017.com | DNS only | www.entorno35.com goes to Vercel (Next.js frontend). |
| CNAME | api | ...cfargotunnel.com | Proxied | api.entorno35.com goes through a Cloudflare Tunnel to your Ubuntu server. |

So today:

- **Frontend**: Only **www.entorno35.com** serves the Vercel app. **entorno35.com** (no www) goes to 216.198.79.1, so users opening `https://entorno35.com` see whatever is on that IP, not the Vercel site.
- **Backend**: **api.entorno35.com** is correctly routed via a Cloudflare Tunnel (proxied). No need to expose port 8080 to the internet; `cloudflared` on the server connects out to Cloudflare and receives API traffic.

**CORS and API URL**: Use the origin users actually see in the browser. If the canonical frontend URL is **https://www.entorno35.com**, set:

- Vercel: `NEXT_PUBLIC_API_URL=https://api.entorno35.com`
- Backend (Portainer): `CORS_ORIGIN=https://www.entorno35.com` (no trailing slash)

If you later make the root domain serve the Vercel app (see below), you can use `CORS_ORIGIN=https://entorno35.com` or allow both origins in the backend.

**Optional: serve the app on the root domain (entorno35.com)**  
If you want `https://entorno35.com` to show the same Vercel app as www:

1. In Vercel: add the domain **entorno35.com** (without www) in Project > Settings > Domains.
2. Vercel will show the required DNS (often a CNAME for the root, e.g. `cname.vercel-dns.com`, or A records). Some registrars/DNS support CNAME flattening at the root; Cloudflare does (they resolve root CNAME to an A).
3. In Cloudflare DNS: either replace the current A record for entorno35.com with the CNAME/A Vercel gives you, or add the CNAME and remove the A. After propagation, root and www will both serve the Vercel app.

If you prefer to keep the root on 216.198.79.1 (e.g. a landing page there), you can instead add a redirect in Cloudflare (Rules > Redirect Rules) so `https://entorno35.com` redirects to `https://www.entorno35.com`. Then the single canonical frontend URL is www, and CORS_ORIGIN stays `https://www.entorno35.com`.

## Checklist before first push

1. **Vercel**: Project linked to repo, Root Directory = `web/frontend`, `NEXT_PUBLIC_API_URL` set to backend URL, `git.deploymentEnabled` in `vercel.json` includes your deployment branches.
2. **GitHub**: `DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN` set; optionally `PORTAINER_*` secrets for automatic stack updates.
3. **Portainer**: Stack created from `docker-compose.prod.yml`, env vars set. For production (develop) leave `IMAGE_TAG` unset or set to `develop`.
4. **Cloudflare**: Use "Investigating your Cloudflare setup" above; ensure DNS (and proxy if used) point correctly and CORS matches the frontend origin.

After that, a commit to an enabled branch will refresh the Vercel site (for any push), and the Docker container will refresh when the workflow runs (backend-path or manual trigger) and Portainer secrets are configured.
