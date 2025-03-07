# Development Environment for Traefik v3 Path Parameter Headers Plugin

This directory contains the development environment for testing the Traefik Path Parameter Headers plugin.

## Prerequisites

- Docker and Docker Compose
- Go 1.20 or later

## Directory Structure

```
traefik-path-param-headers/
├── .traefik.yml                # Plugin manifest
├── go.mod                      # Go module file
├── pathparamheaders.go         # Main plugin code (package traefik_path_param_headers)
├── dev/
│   ├── docker-compose.yml      # Docker Compose for development
│   ├── traefik/
│   │   ├── traefik.yml         # Traefik static configuration
│   │   └── dynamic/            # Directory for dynamic configurations
│   │       └── services.yml    # Dynamic configuration with routers and services
│   └── backend/
│       ├── Dockerfile          # Backend service Dockerfile
│       └── server.go           # Simple backend service
```

## Getting Started

1. Start the development environment:

```bash
cd dev
docker-compose up -d
```

2. Test the plugin with a sample request:

```bash
curl -H "Host: test.localhost" http://localhost/products/electronics/12345
```

3. Expected response:

```json
{
  "headers": {
    "Accept": "*/*",
    "Accept-Encoding": "gzip",
    "User-Agent": "curl/8.7.1",
    "X-Forwarded-For": "172.19.0.1",
    "X-Forwarded-Host": "test.localhost",
    "X-Forwarded-Port": "80",
    "X-Forwarded-Proto": "http",
    "X-Forwarded-Server": "c45947edd753",
    "X-Path-Category": "electronics",
    "X-Path-Id": "12345",
    "X-Real-Ip": "172.19.0.1"
  },
  "method": "GET",
  "path": "/products/electronics/12345",
  "pathParams": {
    "Category": "electronics",
    "Id": "12345"
  },
  "queryParams": {}
}
```

4. Check the Traefik dashboard:

```
http://localhost:8080/dashboard/
```

## Configuration

### Static Configuration (traefik.yml)

The static configuration sets up Traefik with local plugin support:

```yaml
experimental:
  localPlugins:
    pathparamheaders:
      moduleName: github.com/sistemica/traefik-path-param-headers
```

### Dynamic Configuration (dynamic/services.yml)

The dynamic configuration sets up routers, services, and uses the plugin:

```yaml
http:
  middlewares:
    path-params:
      plugin:
        pathparamheaders:
          pathPattern: "/products/{category}/{id}"
          headerPrefix: "X-Path-"
  
  services:
    backend:
      loadBalancer:
        servers:
          - url: "http://backend:8000"
  
  routers:
    test-router:
      rule: "Host(`test.localhost`) && PathRegexp(`^\\/products\\/([^\\/]+)\\/([^\\/]+)$`)"
      service: "backend"
      middlewares:
        - "path-params"
```

Note that we're using `PathRegexp` with a specific regular expression that matches the exact pattern we want to extract parameters from. This ensures that only valid paths with exactly two segments after `/products/` are matched.

## Docker Volumes

The Docker Compose setup mounts several volumes:

1. Static Traefik configuration:
   ```
   ./traefik/traefik.yml:/etc/traefik/traefik.yml
   ```

2. Dynamic configuration directory:
   ```
   ./traefik/dynamic:/etc/traefik/dynamic
   ```

3. Plugin source code:
   ```
   ../:/plugins-local/src/github.com/sistemica/traefik-path-param-headers
   ```

## Testing with Different Path Patterns

To test different path patterns:

1. Update the `pathPattern` in `traefik/dynamic/services.yml`:
   ```yaml
   pathPattern: "/api/users/{userId}/posts/{postId}"
   ```

2. Update the router rule to match your new path pattern:
   ```yaml
   rule: "Host(`test.localhost`) && PathRegexp(`^\\/api\\/users\\/([^\\/]+)\\/posts\\/([^\\/]+)$`)"
   ```

3. Restart the services:
   ```bash
   docker-compose restart
   ```

4. Test with the new path:
   ```bash
   curl -H "Host: test.localhost" http://localhost/api/users/john/posts/123
   ```
   
   Expected result with headers:
   - `X-Path-UserId: john`
   - `X-Path-PostId: 123`

## Another Example: API Version Routes

1. Update the configuration:
   ```yaml
   # Middleware configuration
   pathPattern: "/api/{version}/resources/{id}"
   
   # Router rule
   rule: "Host(`test.localhost`) && PathRegexp(`^\\/api\\/([^\\/]+)\\/resources\\/([^\\/]+)$`)"
   ```

2. Test with:
   ```bash
   curl -H "Host: test.localhost" http://localhost/api/v2/resources/42
   ```
   
   Expected result with headers:
   - `X-Path-Version: v2`
   - `X-Path-Id: 42`

## Troubleshooting

- Make sure your go.mod module name matches what's in the static configuration
- Ensure the package name in .go files is `traefik_path_param_headers`
- Check Traefik logs for errors: `docker-compose logs traefik`
- Verify your .traefik.yml manifest is present and correctly formatted
- Make sure the dynamic configuration file has a valid extension (.yml, .yaml, .toml, or .json)
- Ensure that your `PathRegexp` router rule matches the actual paths you want to process
- Check for syntax errors in your regular expressions by testing them separately