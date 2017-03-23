# a-Journal

a-journal is an example project on how one can structure a web project in Go.

It consists of a bunch of REST services which delivers data formatted in JSON as well as a frontend app that is written with the javascript library RiotJS (http://riotjs.com). Bulma (http://bulma.io) css framework is used for layout and styling.

The goal of this project is to serve as an example to build a SaaS solution using mostly stdlib. I do however use third party libraries in places where i see no need to roll my own.

The future of this project could end up in an actually hosted paid version, however the code base will remain open for either self hosting or for verifying that the project "doesn't do evil".

It is still in heavy development so things are not stable and many features are still missing.

Planned Coming features:

*   Support for multiple membership tiers
*   Stripe integration
*   Solid backend logging
*   Frontend logging
*   Better documentation of everything
*   Add Integration tests to services

## Dev setup

The code is go gettable and can be fetched by

```
go get -u bitbucket.org/sketchground/ajournal
```

For database migrations i use goose [https://bitbucket.org/liamstask/goose](https://bitbucket.org/liamstask/goose/).
To migrate the database, simply run `goose up` inside the root folder of the project.

For packaging and deploying i use a tool called Sup [https://github.com/pressly/sup](https://github.com/pressly/sup) which makes it easy to deploy to multiple hosts etc. It is not strictly necessary and one could run the build scripts located in the scripts folder manually. Remember to set the correct environment variables for this.