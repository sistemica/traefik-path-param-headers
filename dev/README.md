# Development Environment

This directory contains the development environment for the Traefik Path Parameter Headers plugin.

## Components

- **Traefik v3**: Running with plugin development mode enabled
- **Backend service**: A simple Go HTTP server that echoes back request information, including headers

## Getting Started

1. Start the environment:

```bash
docker-compose up -d
```

2. Test the plugin with a sample request:

```bash
curl -H "Host: test.localhost" http://localhost/products/electronics/12345
```

3. Check the response and backend logs:

```bash
docker-compose logs backend
```

You should see the path parameter headers (`X-Path-Category` and `X-Path-Id`) in the response and logs.

## Testing Different Patterns

You can modify the `pathPattern` in the `docker-compose.yml` file to test different path patterns:

```yaml
- "traefik.http.middlewares.path-params.plugin.pathparamheaders.pathPattern=/products/{category}/{id}"
```

After changing the configuration, restart the services:

```bash
docker-compose restart
```

## Debugging

- Use `docker-compose logs traefik` to see Traefik logs
- Use `docker-compose logs backend` to see backend logs
- Check the Traefik dashboard at `http://localhost:8080/dashboard/`

## Stopping the Environment

```bash
docker-compose down
```