# Quick Start Guide

This guide will help you get started with the Traefik Path Parameter Headers plugin quickly.

## Prerequisites

- Docker and Docker Compose
- Go 1.18 or later (for development)
- Git

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/traefik-path-param-headers.git
cd traefik-path-param-headers
```

2. Start the development environment:

```bash
cd dev
docker-compose up -d
```

3. Test the plugin:

```bash
curl -H "Host: test.localhost" http://localhost/products/electronics/12345
```

You should see a JSON response that includes the extracted path parameters as headers.

## Using in Your Traefik v3 Setup

### Static Configuration

Add the plugin to your Traefik static configuration:

```yaml
# traefik.yml
plugins:
  pathparamheaders:
    moduleName: "github.com/yourusername/traefik-path-param-headers"
    version: "v1.0.0"
```

### Dynamic Configuration

Configure the middleware in your dynamic configuration:

```yaml
# dynamic-config.yml
http:
  middlewares:
    product-params:
      plugin:
        pathparamheaders:
          pathPattern: "/products/{category}/{id}"
          headerPrefix: "X-Path-"
```

Apply the middleware to a router:

```yaml
http:
  routers:
    products:
      rule: "Host(`example.com`) && Path(`/products/{category}/{id}`)"
      middlewares:
        - "product-params"
      service: "product-service"
```

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `pathPattern` | The URL pattern with parameters in curly braces | (required) |
| `headerPrefix` | Prefix for the header names | "X-Path-" |

## Examples

### Basic Example

```yaml
pathPattern: "/products/{category}/{id}"
```

When a request is made to `/products/electronics/12345`, the middleware adds:
- Header `X-Path-Category: electronics`
- Header `X-Path-Id: 12345`

### Nested Resources Example

```yaml
pathPattern: "/api/users/{userId}/posts/{postId}"
```

When a request is made to `/api/users/john123/posts/987`, the middleware adds:
- Header `X-Path-UserId: john123`
- Header `X-Path-PostId: 987`

## Troubleshooting

- If the path doesn't match the pattern, no headers will be added
- Check Traefik logs for any errors
- Verify the middleware is correctly attached to your router

For more details, see the [full documentation](README.md).