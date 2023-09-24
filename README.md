# Inventory-Go

This is an API service written in Go for Inventory Management.

The purpose of this project is to demonstrate how to develop a complete API service with Database (MySQL) integration and configuration management in Go. The primary packages used in this project are:
- github.com/gorilla/mux
- github.com/go-sql-driver/mysql
- github.com/spf13/viper

## Hexagonal Architecture
This implementation is inspired by the blog, [The Practical Hexagonal Architecture for Golang](https://medium.com/@briannqc/the-practical-hexagonal-architecture-for-golang-742a49bc8d89). This enables us to decouple the core application logic with the external factors eg, database, api etc, hence, making it more scalable and expandable. This way we can add multiple input or output ports & adapters while keeping the application logic same as it is.

## Getting started

### Pre-requisites

The following needs to be installed before starting:
- go (1.19)
- docker

### Download Packages

To get started with the project locally, first download the required go packages

```sh
go mod download
```

### Build and Run

To build the service run the following command:

```sh
go build -o inventory ./cmd/main.go
```

The above command would generate a binary, `inventory`, in the project's root. Execute the binary to start the service:

```sh
./inventory
```

> Note: The service would expect a MySQL server running on `localhost:3306` by default. If the database server is not running, the inventory service would exit with a `connection refused` error. Make sure you have a MySQL server running prior to running this service.

### Configuration

Following environment variables needs to be set:

|Variable|Description|Default Value|
|-|-|-|
|MYSQL_HOST|MySQL Server Host|"127.0.0.1"|
|MYSQL_PORT|MySQL Server Port|3306|
|MYSQL_USER|MySQL Server Database User|"inventory_admin"|
|MYSQL_PASSWORD|MySQL Server User Password|"password"|
|MYSQL_DATABASE|MySQL Server Database Name|"inventory"|
|SERVER_PORT|Port on which inventory server would be exposed|"8080"|

Alternatively, a `config.toml` file in the project's root may be used to set these variables.

### Docker
As an alternate, docker compose can also be used to build and run the server.

To build the docker image of the inventory server, run the following command:

```sh
docker compose build
```

After the build is complete, run the following command to start the server:

```sh
docker compose up
```

> Note: You don't need to spin up a MySQL server before hand in this case, since docker compose takes care of that and starts a MySQL server for you.

## Contributing

Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request