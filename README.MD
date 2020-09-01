### nginx-ldap-auth

A simple example of LDAP authentication using ngx_http_auth_request_module by golang.

[中文](https://github.com/xiaojun207/nginx-ldap-auth/blob/master/README_CN.MD)

### usage
It is recommended that you use docker for reference. Like this:
```
docker run -d --name nginx-ldap-auth -v /data/docker/nginx-ldap-auth/cfg.json:/app/cfg.json -p 8080:8080 xiaojun207/nginx-ldap-auth:latest
```

### cfg.json example
```
{
     // ldap config
    "ldap": {
        "addr": "ldap.example.org:389",
        "baseDn": "dc=example,dc=org",
        "bindDn": "cn=manager,dc=example,dc=org",
        "bindPass": "password",
        "authFilter": "(&(uid=%s))",
        "attributes": ["uid", "cn", "mail"],
        "tls":        false,
        "startTLS":   false
    },
    // control config
    "control":{
        "ipAcl":{
            "deny":["127.0.0.1","192.168.0.0/24","192.168.1.0-192.168.1.255"], // deny ips
            "direct":[] // direct ips
        },
        "timeAcl":{
            "deny":["00:00-08:00","17:00-23:59"], // time deny
            "direct":[]  // time direct
        },
        "allowUser":["user1"],  // user allowed to login, should be [] if all ldap users are allowed
        "users":[{  // username and password，can use in none ldap env.
            "userName": "admin",
            "passWord": "123456",
            "tryNum": 3
        }]
    },
    "http": {
        "debug":false,
        "ips":["127.0.0.1"], // api trust ips
        "listen": "0.0.0.0:8080"
    }
}
```

### nginx demo for kibana

```
error_log logs/error.log debug;

events { }

http {
    proxy_cache_path cache/  keys_zone=auth_cache:10m;

    upstream auth_service {
        server 192.168.3.32:8080;
    }
    upstream kibana {
        server 192.168.3.32:5601;
    }

    server {
        listen 80;

        root html;
        index index.html;

        location / {
            auth_request /auth/auth-proxy;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header Host            $http_host;

	          error_page 401 =200 /auth/login;
            proxy_pass http://kibana;
        }

        location /auth/ {
              proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
              proxy_set_header Host            $http_host;
              proxy_set_header X-Target $request_uri;
              proxy_pass http://auth_service;
        }

        location = /auth/auth-proxy {
            internal;

            proxy_pass http://auth_service;

            proxy_pass_request_body off;
            proxy_set_header Content-Length "";
            proxy_cache auth_cache;
            proxy_cache_valid 200 10m;

            proxy_cache_key "$http_authorization$cookie_sessionID";

            proxy_set_header X-CookieName "sessionID";
            proxy_set_header Cookie sessionID=$cookie_sessionID;
        }
    }
}

```
