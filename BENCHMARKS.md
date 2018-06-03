# Benchmarks

In this document you will find benchmark results for different endpoints of ORY Hydra. All benchmarks are executed
using [rakyll/hey](https://github.com/rakyll/hey). Please note that these benchmarks run against the in-memory storage
adapter of ORY Hydra. These benchmarks represent what performance you would get with a zero-overhead database implementation.

We do not include benchmarks against databases (e.g. MySQL or PostgreSQL) as the performance greatly differs between
deployments (e.g. request latency, database configuration) and tweaking individual things may greatly improve performance.
We believe, for that reason, that benchmark results for these database adapters are difficult to generalize and potentially
deceiving. They are thus not included.

This file is updated on every push to master. It thus represents the benchmark data for the latest version.

All benchmarks run 10.000 requests in total, with 100 concurrent requests.

## OAuth 2.0

This section contains various benchmarks against

### Token Introspection

```

Summary:
  Total:	0.4724 secs
  Slowest:	0.1059 secs
  Fastest:	0.0002 secs
  Average:	0.0046 secs
  Requests/sec:	21169.1376
  
  Total data:	1550000 bytes
  Size/request:	155 bytes

Response time histogram:
  0.000 [1]	|
  0.011 [9694]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.021 [204]	|■
  0.032 [1]	|
  0.042 [0]	|
  0.053 [0]	|
  0.064 [0]	|
  0.074 [0]	|
  0.085 [14]	|
  0.095 [77]	|
  0.106 [9]	|


Latency distribution:
  10% in 0.0008 secs
  25% in 0.0018 secs
  50% in 0.0032 secs
  75% in 0.0051 secs
  90% in 0.0076 secs
  95% in 0.0094 secs
  99% in 0.0797 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0007 secs, 0.0002 secs, 0.1059 secs
  DNS-lookup:	0.0002 secs, 0.0000 secs, 0.0531 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0044 secs
  resp wait:	0.0035 secs, 0.0001 secs, 0.0482 secs
  resp read:	0.0002 secs, 0.0000 secs, 0.0057 secs

Status code distribution:
  [200]	10000 responses



```

### Client Credentials Grant

ORY Hydra uses BCrypt to obfuscate secrets of OAuth 2.0 Clients. When using flows such as the OAuth 2.0 Client Credentials
Grant, ORY Hydra validates the client credentials using BCrypt which causes (by design) CPU load. CPU load and performance
depend on the BCrypt cost which can be set using the environment variable `BCRYPT_COST`. For these benchmarks,
we have set `BCRYPT_COST=8`

```
```

## OAuth 2.0 Client Management

### Creating OAuth 2.0 Clients

```

Summary:
  Total:	166.4903 secs
  Slowest:	1.9983 secs
  Fastest:	0.0286 secs
  Average:	1.6563 secs
  Requests/sec:	60.0636
  
  Total data:	2960000 bytes
  Size/request:	296 bytes

Response time histogram:
  0.029 [1]	|
  0.226 [7]	|
  0.423 [9]	|
  0.620 [11]	|
  0.816 [10]	|
  1.013 [11]	|
  1.210 [14]	|
  1.407 [11]	|
  1.604 [2421]	|■■■■■■■■■■■■■■
  1.801 [6903]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  1.998 [602]	|■■■


Latency distribution:
  10% in 1.5426 secs
  25% in 1.6045 secs
  50% in 1.6648 secs
  75% in 1.7204 secs
  90% in 1.7713 secs
  95% in 1.8132 secs
  99% in 1.9173 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0009 secs, 0.0286 secs, 1.9983 secs
  DNS-lookup:	0.0003 secs, 0.0000 secs, 0.0661 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0011 secs
  resp wait:	1.6553 secs, 0.0284 secs, 1.9983 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0006 secs

Status code distribution:
  [201]	10000 responses



```

### Listing OAuth 2.0 Clients

```

Summary:
  Total:	1.3342 secs
  Slowest:	0.1486 secs
  Fastest:	0.0004 secs
  Average:	0.0131 secs
  Requests/sec:	7495.0405
  

Response time histogram:
  0.000 [1]	|
  0.015 [8892]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.030 [961]	|■■■■
  0.045 [46]	|
  0.060 [0]	|
  0.075 [0]	|
  0.089 [36]	|
  0.104 [52]	|
  0.119 [7]	|
  0.134 [4]	|
  0.149 [1]	|


Latency distribution:
  10% in 0.0091 secs
  25% in 0.0111 secs
  50% in 0.0122 secs
  75% in 0.0135 secs
  90% in 0.0155 secs
  95% in 0.0184 secs
  99% in 0.0820 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0007 secs, 0.0004 secs, 0.1486 secs
  DNS-lookup:	0.0002 secs, 0.0000 secs, 0.0496 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0086 secs
  resp wait:	0.0117 secs, 0.0003 secs, 0.0582 secs
  resp read:	0.0007 secs, 0.0000 secs, 0.0269 secs

Status code distribution:
  [200]	10000 responses



```

### Fetching a specific OAuth 2.0 Client

```

Summary:
  Total:	0.3439 secs
  Slowest:	0.1068 secs
  Fastest:	0.0001 secs
  Average:	0.0034 secs
  Requests/sec:	29080.4190
  
  Total data:	2650000 bytes
  Size/request:	265 bytes

Response time histogram:
  0.000 [1]	|
  0.011 [9892]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.021 [7]	|
  0.032 [0]	|
  0.043 [0]	|
  0.053 [0]	|
  0.064 [0]	|
  0.075 [0]	|
  0.085 [0]	|
  0.096 [11]	|
  0.107 [89]	|


Latency distribution:
  10% in 0.0009 secs
  25% in 0.0014 secs
  50% in 0.0021 secs
  75% in 0.0031 secs
  90% in 0.0044 secs
  95% in 0.0056 secs
  99% in 0.0927 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0008 secs, 0.0001 secs, 0.1068 secs
  DNS-lookup:	0.0003 secs, 0.0000 secs, 0.0663 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0023 secs
  resp wait:	0.0019 secs, 0.0001 secs, 0.0488 secs
  resp read:	0.0004 secs, 0.0000 secs, 0.0077 secs

Status code distribution:
  [200]	10000 responses



```
