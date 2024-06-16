## Build & push instructions
```bash
docker buildx build --platform=linux/amd64,linux/arm64 --push -f tests/databases/postgres/client/Dockerfile -t otterize/postgres-integration-test-client:latest .
```