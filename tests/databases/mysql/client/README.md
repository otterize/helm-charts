# MySQL client test app

## Build & push instructions

```bash
docker buildx build --platform=linux/amd64,linux/arm64 --push -f tests/databases/mysql/client/Dockerfile -t otterize/mysql-integration-test-client:latest .
```
