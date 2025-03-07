# Quick Start Guide

This guide will help you get started with the Traefik Path Parameter Headers plugin quickly.

## Prerequisites

- Traefik v3
- Access to configure Traefik static and dynamic configuration

## Installation

### Option 1: Using from GitHub (Production)

1. Add the plugin to your Traefik static configuration:

```yaml
experimental:
  plugins:
    pathparamheaders:
      moduleName: github.com/sistemica/traefik-path-param-headers
      version: v1.0.0  # Use the appropriate version
```

### Option 2: Local Development

1. Clone the repository:

```bash
git clone https://github.com/sistemica/traefik-path-param-headers.git
cd traefik-path-param-headers
```

2. Configure Traefik to use the local plugin:

```yaml
experimental:
  localPlugins:
    pathparamheaders:
      moduleName: github.com/sistemica/traefik-path-param-headers
```

3. Mount the plugin source code in your Traefik container:

```yaml
volumes:
  - ./:/plugins-local/src/github.com/sistemica/traefik-path-param-headers
```

## Configuration

### Dynamic Configuration

Configure the middleware in your dynamic configuration:

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

## Example

When a request is made to `/products/electronics/12345`, the middleware adds:
- Header `X-Path-Category: electronics`
- Header `X-Path-Id: 12345`

Example response showing these headers:

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

## More Examples

### Nested Resources

```yaml
# Middleware configuration
pathPattern: "/api/users/{userId}/posts/{postId}"

# Router rule
rule: "Host(`test.localhost`) && PathRegexp(`^\\/api\\/users\\/([^\\/]+)\\/posts\\/([^\\/]+)$`)"
```

When a request is made to `/api/users/john123/posts/987`, the middleware adds:
- Header `X-Path-UserId: john123`
- Header `X-Path-PostId: 987`

### API Versioning

```yaml
# Middleware configuration
pathPattern: "/api/{version}/resources/{id}"

# Router rule
rule: "Host(`test.localhost`) && PathRegexp(`^\\/api\\/([^\\/]+)\\/resources\\/([^\\/]+)$`)"
```

When a request is made to `/api/v2/resources/42`, the middleware adds:
- Header `X-Path-Version: v2`
- Header `X-Path-Id: 42`

## Path Matching

The `PathRegexp` rule in the examples above is carefully crafted to:
- Match only the exact number of path segments needed
- Ensure each segment can contain any characters except a forward slash
- Anchor the expression to exactly match the entire path

For example, `^\\/products\\/([^\\/]+)\\/([^\\/]+)$` will:
- Match: `/products/electronics/12345`
- Not match: `/products/electronics/12345/details` (too many segments)
- Not match: `/products/electronics` (too few segments)

## Troubleshooting

- If the path doesn't match the pattern, no headers will be added
- Ensure your `PathRegexp` pattern matches the URLs you're testing with
- Check Traefik logs for any errors
- Verify the middleware is correctly attached to your router