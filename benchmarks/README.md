# **`BENCHMARK`**

## **`Install K6`**

```sh
> brew install k6
```

## **`Run benchmark`**

```sh
> k6 run --vus=100 --duration 1m --summary-export=ping.json ping.js
```
