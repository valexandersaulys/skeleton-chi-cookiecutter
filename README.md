# Skeleton Chi Project

**Forever a work in Progress!**

Templatized, via [Cookiecutter](https://cookiecutter.readthedocs.io/en/2.5.0/index.html), version of [my Skeleton Chi Project](https://github.com/valexandersaulys/Skeleton-Golang-Webapp). 

Tested on MacOS & Linux. 
```
cookiecutter https://github.com/valexandersaulys/skeleton-chi-cookiecutter
cd chi-blog
make install
make run
```


## Notes on Usage

### Atlas Migrations
Support for doing database migrations supported through [Atlas](https://atlasgo.io):
```sh
atlas migrate diff --env=gorm 
atlas migrate apply --env=gorm --url="sqlite:///tmp/chiblog.db?_journal_mode=WAL"
```
Modify the `atlas.hcl` file to your needs. Note that you'll need to supply the necesary url such that it matches the paths in `models.go`. 

### Sessions Middleware

If using encryption, encryption keys need to be exactly 16, 24, or 32. Authentication keys can be any length. These are stored under the following environmental variables:

+ `COOKIE_STORE_AUTH_KEY=COOKIE_STORE_AUTH_KEY`:  authentication key
+ `COOKIE_STORE_ENCRYPT_KEY=123456781234567812345678`:  _encryption key_


## License

[MIT License](https://mit-license.org)
```
Copyright © 2024 Vincent Alexander Saulys

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
