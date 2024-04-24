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
