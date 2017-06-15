# j2xrp

A reverse proxy that converts a JSON request to an XML request

We host an instance of this application on https://j2xrp.herokuapp.com/.

## Usage

### Example

With the hosted Heroku application:

```
curl -vvv \
  -H "Content-Type: application/json" \
  -X POST -d '{"username":"xyz","password":"xyz"}' \
    https://j2xrp.herokuapp.com/
```

### Build

#### Static Binary

Build a fully static binary:

    $ CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/j2xrp .

#### Docker

    $ docker build -t flaccid/j2xrp .

### Run

You can run from the main entrypoint locally without building:

    $ go run main.go


License and Authors
-------------------
- Author: Chris Fordham (<chris@fordham-nagy.id.au>)

```text
Copyright 2017, Chris Fordham

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
