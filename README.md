# Rinha de Backend - Segunda Edição

Versão Go da [rinha de backend 2ª edição - 2024/Q1](https://github.com/zanfranceschi/rinha-de-backend-2024-q1). 

## Stack

- go
- nginx
- postgresql

## Rodando o projeto

```bash
docker compose up nginx -d --build
```

## Resultados

### Resultado do Gatling local

Todas requisições abaixo de 800ms. (Estes testes utilizaram um máximo de 250MB RAM distribuidos entre os recursos. 60% menos recurso de memória RAM do que o permitido pela rinha!

![Gatling](docs/screenshots/gatling-1.png)

![Gatling](docs/screenshots/gatling-2.png)

## Métricas dos testes

Métricas colhidas no Docker Desktop após a execução do teste. O teste foi executado em um Mac Mini M1 16GB RAM/512GB SSD.

- Banco de dados (Postgresql)

![Banco de dados](docs/screenshots/metrica-banco-de-dados.png)

- Endpoints (Go)

![Endpoint 1 da API](docs/screenshots/metrica-api-endpoint-1.png)

![Endpoint 1 da API](docs/screenshots/metrica-api-endpoint-2.png)

- Proxy reverso (Nginx)

![Proxy reverso](docs/screenshots/metrica-proxy-reverso.png)

## Versões alternativas

### Implementações que elaborei em outras linguagens

- [rinha2-back-end-dotnet](https://github.com/jonathanperis/rinha2-back-end-dotnet)
- [rinha2-back-end-rust](https://github.com/jonathanperis/rinha2-back-end-rust)
- [rinha2-back-end-python](https://github.com/jonathanperis/rinha2-back-end-python)
- [rinha2-back-end-postgrest](https://github.com/jonathanperis/rinha2-back-end-postgrest)