# RedPlanet User Service

RUS is a Go application designed to run periodic tasks using the Cron library. It connects to Mikrotik routers to fetch user data and interacts with the Saeplus API to retrieve additional information. The gathered data is then transformed and stored in a MongoDB database.

### Configuration

The application relies on environment variables for configuration. Set the following variables:

- `MIKROTIK_API_USER`: Mikrotik API username.
- `MIKROTIK_API_PASS`: Mikrotik API password.
- `MONGO_USER`: MongoDB username.
- `MONGO_PASS`: MongoDB password.
- `MONGO_HOST`: MongoDB host address.
- `SAEPLUS_TOKEN_HEADER`: Saeplus token header.
- `SAEPLUS_TOKEN`: Saeplus token.
- `SAEPLUS_API_HEADER`: Saeplus api header.
- `SAEPLUS_API_CONNECT`: Saeplus api key.