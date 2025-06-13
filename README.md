# jwt-auth-validator
Golang app with jwt validator endpoint

This app provides signing validation of a JWT key. When used with an nginx proxy, this can provide verification the JWT was signed by a system a specific JWT_SECRET. If the JWT_SECRET was not used during the signing process, the key will be seen as invalid.


## Nginx Proxy Example

The goal of this proxy example is to serve as a gatekeeper to an upstream Api server.

The following is a starting example. It will require some modications and understanding of nginx configurations to make it work for your scenario. It will not work "out of the box."

**Assumptions**
- HTTP header is sent as "Cookie" in a value format of "token=<jwt>"
- Api server exists as "api.domain.com" on port 8080 (referenced in upstream)
- Domain name points to this server (referenced as server_name)


```conf
upstream apiUpstream {
  server api.domain.com:8080 max_fails=0 fail_timeout=10s;
  keepalive 512;
}

upstream authUpstream {
  server auth:8080 max_fails=0 fail_timeout=10s;
  keepalive 512;
}

server {
  listen 80;
  server_name my.domain.com;

  location ~ ^/(.*) {
    auth_request /auth-server/validate;
    auth_request_set $auth_status $upstream_status;

    error_page 401 = @fallback;
    error_page 404 = @notfound;
    add_header X-Auth-Status $auth_status;

    proxy_pass http://apiUpstream;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }

  location /auth-server/ {
    internal;
    proxy_pass http://authUpstream/;
    proxy_buffers 7 16k;
    proxy_buffer_size 31k;
  }

  location /   location @fallback {
    return 302 https://<my.domain.com>/auth?redirect=$uri?$query_string;
  }
}
```
