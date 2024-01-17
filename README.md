# Skeleton Chi Project

**Work-In-Progress**

Templatized, via [Cookiecutter](#), version of [my Skeleton Chi Project](https://github.com/valexandersaulys/Skeleton-Golang-Webapp). 


## To Replace with Cookiecutter
```
vincent@Vincents-MacBook-Pro ~/w/s/{{cookiecutter.directory_name}}> ag "example/skeleton"
middleware/auth_test.go
4:	"example/skeleton/config"
5:	"example/skeleton/models"
6:	"example/skeleton/services"

middleware/session.go
5:	"example/skeleton/config"
6:	"example/skeleton/models"

middleware/auth.go
5:	// "example/skeleton/config"

go.mod
1:module example/skeleton   #  COOKIECUTTER

app/app_test.go
6:	"example/skeleton/config"
7:	"example/skeleton/middleware"

app/app.go
4:	"example/skeleton/config"
5:	"example/skeleton/views"

models/models.go
4:	"example/skeleton/config"

Makefile
29:	RUNTIME_ENV=TESTING ADD_DUMMIES=1 dlv test example/skeleton/views -test.v

models/models_test.go
4:	"example/skeleton/config"

views/hello_world.go
4:	"example/skeleton/middleware"

views/auth.go
4:	"example/skeleton/middleware"
5:	"example/skeleton/services"
6:	tmpl "example/skeleton/templates"

views/views.go
5:	"example/skeleton/middleware"
6:	tmpl "example/skeleton/templates"

views/auth_test.go
4:	"example/skeleton/middleware"
5:	"example/skeleton/models"

views/views_test.go
4:	"example/skeleton/config"
5:	"example/skeleton/middleware"

main.go
4:	"example/skeleton/app"
5:	"example/skeleton/config"
6:	"example/skeleton/middleware"
7:	"example/skeleton/models"

views/posts.go
4:	"example/skeleton/middleware"
5:	"example/skeleton/models"
6:	"example/skeleton/services"
7:	tmpl "example/skeleton/templates"

services/auth.go
5:	"example/skeleton/models"

views/posts_test.go
5:	"example/skeleton/middleware"
6:	"example/skeleton/models"
7:	"example/skeleton/services"

services/posts.go
5:	"example/skeleton/models"

services/posts_test.go
4:	"example/skeleton/models"

services/services_test.go
4:	"example/skeleton/config"
```

### Replace my email
```
vincent@Vincents-MacBook-Pro ~/w/s/{{cookiecutter.directory_name}}> ag "vincent@saulys.me"
models/post_test.go
11:		Email:    "vincent@saulys.me",

models/models.go
49:		Email:    "vincent@saulys.me",

models/user_test.go
11:		Email:    "vincent@saulys.me",
22:	assert.Equal(t, user.Email, "vincent@saulys.me")

views/auth_test.go
37:			"email":    "vincent@saulys.me",

views/posts_test.go
40:	cookie := start_session("vincent@saulys.me", "password", app)
106:	cookie := start_session("vincent@saulys.me", "password", app)

services/posts_test.go
73:		Email:    "vincent@saulys.me",
```
