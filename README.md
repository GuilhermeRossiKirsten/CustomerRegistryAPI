# CustomerRegistryAPI

O projeto foi desenvolvido dentro de uma máquina virtual no ambiente do codespaces do github

```bash
Linux codespaces-ab61da 6.8.0-1044-azure #50~22.04.1-Ubuntu SMP Wed Dec  3 15:13:22 UTC 2025 x86_64 x86_64 x86_64 GNU/Linux
```


go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

go install github.com/go-swagger/go-swagger/cmd/swagger@latest

swagger generate spec -o internal/docs/swagger.json --scan-models