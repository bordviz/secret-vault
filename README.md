# Secret Vault

![vault](./docs/vault.png)

## Contents
- [About](#about-project)
- [Customization](#customization)
- [Installation](#installation)
- [Usage](#usage)

## About project

Secret Vault is a simple analog of HashiCorp Valut that performs the functions of storing secrets and providing access to them to other users. There are 4 endpoints implemented in this project:
- Create a new Secret Vault (admin side)
- Create an access token to access the vault (administrator side)
- Retrieving data from the vault (administrator side)
- Retrieving data from the vault (user side)

## Customization
If you want to make changes to the startup configuration, you can modify the config file located in config/config.yaml, or you can create your own config file. If you add your own config file, change the CONFIG_PATH field in the .env file.

### Config file
```yaml
env: local # Application startup mode, depending on the selection, will differ the level and appearance of logs
migrations_path: ./migrations # Path to the folder where migrations to the database are located (it is not desirable to change it)

database: # Database connection
  host: localhost # Database host
  port: 5432 # Database port
  user: postgres # Database user
  password: postgres-user # Database password
  name: vault # Name of database sever
  sslmode: disable # Database SSL mode
  attemps: 5 # Number of attempts to connect to the database 
  delay: 5s # Interval between attempts to connect to the database
  timeout: 5s # Abort after a failed attempt to connect to the database

http-server: # HTTP Server configuration
  port: 8200 # HTTP Server port
  timeout: 5s # Read and Write timeout
  idle_timeout: 60s # Wait time to receive a repeat request from a user with an open connection 
```

### .env file

- CONFIG_PATH - Path to the config file
- ROOT_TOKEN - Administrator token for creating storages and creating new access rights 
- SECRET - The secret to creating custom tokens

## Installation
- [Local Installation](#local-installation)
- [Docker](#docker)

### Local installation
To migrate to the database, run the following command from the root of the project:

```bash
go run cmd/migrations/main.go --action=up
```

To start the application, you need to run the following commands in the root of the project:

```bash
go build -o build/main cmd/vault/main.go
```
```bash
./build/main
```
Or run without creating a binary
```bash
go run cmd/vault/main.go
```

### Docker
> IMPORTANT If you have made changes to the configuration files, check them against the parameters specified in docker-compose.yml

To start a project using Docker, run the following command from the root of the project:
```bash
docker-compose up -d
```
To stop, use:
```bash
docker-compose down
```

## Usage
- [Create a new storage](#create-a-new-storage)
- [Create a new user token](#create-a-new-user-token)
- [Retrieve storage as administrator](#retrieve-storage-as-administrator)
- [Retrieve storage as user](#retrieve-storage-as-user)

## Create a new storage    
#### Request
`POST /root/create`

#### Header 
`Authorization: Bearer <root token>`

#### Body
```json
{
    "name": "New vault",
    "data": {
        "param1": "some value",
        "param2": "some value 2",
        ...
    }
}
```
#### Response
```json
{
    "message": "new vault successfully created",
	"id": <vault_id>
}
```

## Create a new user token
#### Request
`POST /root/create-token`

#### Header 
`Authorization: Bearer <root token>`

#### Body
```json
{
    "vault_id": <vault_id>,
    "expires": 3600
}
```
The `expires` field defines the lifetime of the token in seconds.
#### Response
```json
{
	"token": <token>
}
```

## Retrieve storage as administrator
#### Request
`GET /root/get/{vault_id}`

#### Header 
`Authorization: Bearer <root token>`

#### Response
```json
{
    "id": <vault_id>,
    "name": "New vault",
    "data": {
        "param1": "some value",
        "param2": "some value 2",
        ...
    }
}
```
## Retrieve storage as user
#### Request
`GET /user/get`

#### Header 
`Authorization: Bearer <user token>`

#### Response
```json
{
    "id": <vault_id>,
    "name": "New vault",
    "data": {
        "param1": "some value",
        "param2": "some value 2",
        ...
    }
}
```

### ⭐️ If you like my project, don't spare your stars 🙃