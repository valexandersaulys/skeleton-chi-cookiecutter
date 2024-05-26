# Skeleton Chi Blog

Unlike Django & Rails, [`chi`](https://go-chi.io/#/) does not enforce
structure. This is an attempt to layout such a structure. 

## Build & Deploy

Naked Application
```sh
make install
make test
make run
```

Using Docker
```sh
make docker-build
make docker-run
```

Then Deploy!
```sh
git init .
kamal init
# ... modify .env
kamal setup
kamal deploy
```

Make sure you have populated `.env` -- not to be committed:
```sh
KAMAL_REGISTRY_PASSWORD=
CSRF_PROTECTION_KEY=1234
COOKIE_STORE_AUTH_KEY=COOKIE_STORE_AUTH_KEY 
COOKIE_STORE_ENCRYPT_KEY=123456781234567812345678 
RUNTIME_ENV=
ADD_DUMMIES=
```

## Some Subtleties on Usage

### Atlas Migrations
Support for doing database migrations supported through [Atlas](https://atlasgo.io):
```sh
atlas migrate diff --env=gorm 
atlas migrate apply --env=gorm --url="sqlite:///tmp/chiblog.db?_journal_mode=WAL"
```
Modify the `atlas.hcl` file to your needs. Note that you'll need to supply the necesary url such that it matches the paths in `models.go`. 

### Watching SQL Queries in Logs
If you run with `--log-level=trace` argument passed, GORM will push all SQL queries into the trace logs via [a custom logger](models/gorm_logger.go). 

### Sessions Middleware

If using encryption, encryption keys need to be exactly 16, 24, or 32. Authentication keys can be any length. These are stored under the following environmental variables:

+ `COOKIE_STORE_AUTH_KEY=COOKIE_STORE_AUTH_KEY`:  authentication key
+ `COOKIE_STORE_ENCRYPT_KEY=123456781234567812345678`:  _encryption key_
