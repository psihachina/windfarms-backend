# windfarms-backend
windfarms-backend is http api server for service of imitation modeling windfarm works
## Installation
Requirement Golang 1.16 and highest

Prepare the database using export or [migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md) in advance

Use the make to install modeules
```bash
make install
```

Setup env
```bash
export GOOGLE_MAPS_API_KEY=<API_KEY>
export DB_PASSWORD=<PASSWORD>
export DB_USERNAME=<USERNAME>
```

## Usage
```bash
make run
```
## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
## License