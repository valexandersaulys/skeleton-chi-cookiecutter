
+ [ ] add [custom support](https://pkg.go.dev/github.com/gorilla/sessions#FilesystemStore) for sqlite in [sessions](middleware/session.go#L14)
+ [ ] add redis support for [sessions](middleware/session.go#L14)
+ [ ] go back to vanilla Go for handling HTMX
+ [ ] figure out a simple file upload both locally and with s3
+ [ ] Static file gathering, minifying, and pushing to s3 -- probably as a script?

+ [ ] [yasnippets](http://joaotavora.github.io/yasnippet/) to create
      Golang things
  + [ ] generic file starting with package
  + [ ] generic test
  + [ ] generic file one-file with main/init
  + [ ] generic test w/stretchr
  + [ ] generic test for a route (no content)
  + [ ] template variable via go:embed
  + [ ] compiled template
  + [ ] generic route with template pull and compiling, executed against anonymous struct
  + [ ] `sort.Slice` (no idea what I meant here)
  + [ ] extracting a cookie
  + [ ] extract html via `goquery`
  + [ ] iterating over resp cookies and adding them to a further request
  + [ ] full request-response-checkbody 
  + [ ] generic handler via `func(w http.ResponseWriter, r *http.Request) { })`



### Down the Road

+ [ ] [Very Simple Binary deploy with Traefik](https://stackoverflow.com/questions/58496270/traefik-v2-as-a-reverse-proxy-without-docker) or leveraging [Nginx/certbot combo](https://www.nginx.com/blog/using-free-ssltls-certificates-from-lets-encrypt-with-nginx/)
+ [ ] Various JS setups
  + [ ] [Tailwind CSS setup](https://tailwindcss.com/docs/installation) (via NPM) -- possibly as a cookiecutter flag?
  + [ ] [AlpineJS](https://alpinejs.dev)
  + [ ] [HTMX](https://htmx.org)
+ [ ] Support for non-SQLite databases via CLI
+ [ ] add hCaptcha support =>  [Stop bots and human abuse.](https://www.hcaptcha.com/), possibly with cookiecutter flag?
+ [ ] add [Swagger support](https://github.com/swaggo/http-swagger)


#### Very very down the road

Much less interest on my part with this

+ [ ] Support for MongoDB via CLI
+ [ ] Add HTMX and Hyperscript (download as a post hook, shutil and urllib3 in Python) 
+ [ ] [add simple asynq task queues](https://github.com/hibiken/asynq) 
+ [ ] Add raw SQL query support -- getting stringified responses is kind of difficult

