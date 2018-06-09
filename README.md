# Star Wars API

![travis](https://travis-ci.org/crowleyfelix/star-wars-api.svg?branch=master)

## Pre-requisitos

- [make]
- [go] 1.8.3
- [docker] 18.03.1
- [docker-compose] 1.21.2

## Rodando a aplicação

Para executar a stack (API + banco de dados) da aplicação via docker, execute o seguinte comando:

Iniciar

```bash
make run
```

Parar

```bash
make stop
```

Também é possível executar apenas a API manualmente, porém, será necessário fornecer as variaveis de ambiente descritas no arquivo [.env.example].

```bash
go get -d -u -v ./...
go run server/main.go
```

ou

```bash
go get -d -u -v ./...
go build -o run server/main.go && ./run
```

## Interagindo com a aplicação

Os endpoints estão documentados no formato do swagger e encontram-se em [docs/swagger.yml](docs/swagger.yml).

## Rodando os testes

Para executar os testes automzatizados, execute os seguintes comandos:

- Testes de unidade

```bash
make test
```

- Testes de integração

```bash
make test-integ
```

Também é possível medir o coverage

- Console report

```bash
make cov
```

- Html report

```bash
make cov-html
```

[.env.example]:.env.example
[make]:https://www.gnu.org/software/make/manual/make.html
[go]:https://golang.org/dl/
[docker]:https://www.docker.com/community-edition#/download
[docker-compose]:https://docs.docker.com/compose/install/