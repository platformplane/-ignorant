# Scanner

- [Trivy](https://github.com/aquasecurity/trivy)
- [Gitleaks](https://github.com/gitleaks/gitleaks)

## Scan Directory

```shell
docker run --rm --pull always -v scanner-cache:/cache -v ./:/src -w /src ghcr.io/platformplane/scanner
```