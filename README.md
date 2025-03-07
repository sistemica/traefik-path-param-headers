# Path Parameter Headers Plugin for Traefik v3

This middleware plugin extracts path parameters from URLs based on a pattern and adds them as HTTP headers to the request.

## Features

- Extracts path parameters based on a pattern (e.g., `/products/{category}/{id}`)
- Adds each parameter as an HTTP header
- Customizable header prefix

## Example

With the configuration below, if a request is made to `/products/electronics/12345`, the plugin will add the following headers to the request:

- `X-Path-Category: electronics`
- `X-Path-Id: 12345`

Example response from a backend server showing the extracted headers:

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

## Installation

### Static Configuration

```yaml
# Enable experimental features for plugins
experimental:
  localPlugins:
    pathparamheaders:
      moduleName: github.com/sistemica/traefik-path-param-headers
```

### Dynamic Configuration

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

Note: The `PathRegexp` router rule uses a regular expression to precisely match the path pattern, ensuring that only valid paths with exactly two path segments after `/products/` are matched.

## Configuration Options

| Parameter    | Type   | Required | Default | Description                                                         |
|--------------|--------|----------|---------|---------------------------------------------------------------------|
| pathPattern  | String | Yes      | -       | URL pattern with parameters in curly braces (e.g., `/{param}/{id}`) |
| headerPrefix | String | No       | `X-Path-` | Prefix for the header names                                       |

## Development

See the [development documentation](dev/README.md) for instructions on setting up a local development environment.

## How It Works

The plugin uses regular expressions to extract parameter values from the actual request path based on the pattern provided in the configuration. It then adds these values as headers to the request before passing it to the next middleware or backend service.

Parameter names in the pattern are converted to header names with the first letter capitalized (e.g., `category` becomes `X-Path-Category`).

## Notes

- If the request path doesn't match the pattern, no headers will be added
- The plugin only processes path parameters, not query parameters
- Parameter names are case-sensitive in the pattern
- Using `PathRegexp` in the router rule provides precise path matching, but you can also use `PathPrefix` for simpler configurations