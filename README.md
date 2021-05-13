# tiniapp-backend-oauth-sample

## [API Usage](docs/api_usage.md)

## TODO

- Fill `<your-client-id>`, `<your-client-secret>` to `.env` file.

  ```sh
  API_SERVICE_CLIENT_ID=<your-client-id>
  API_SERVICE_CLIENT_SECRET=<your-client-secret>
  ```

- Add your code.

- Add your tests.

- Add your benchmarks.

## Config

- Your config SHOULD be added to `config.go` and `config.yaml` for automatic binding.

- You can use `.env` to set config for local development.

## Local Development

- Run

```sh
$> ./scripts/bin start
```

- Test

```sh
$> ./scripts/bin start
```

## Local Deployment

```sh
> docker compose -f docker-compose.yml up -d
```

To open shell:

```sh
> docker compose -f docker-compose.yml run app sh
```

## Lint Code

```sh
> ./scripts/bin code_lint
```

## Format Code

```sh
> ./scripts/bin code_format
```
