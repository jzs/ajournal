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

For database migrations i use migrate [https://github.com/mattes/migrate](https://github.com/mattes/migrate).
To migrate the database, simply run `./migrate up` inside the db folder of the project.

For packaging and deploying i use a tool called Sup [https://github.com/pressly/sup](https://github.com/pressly/sup) which makes it easy to deploy to multiple hosts etc. It is not strictly necessary and one could run the build scripts located in the scripts folder manually. Remember to set the correct environment variables for this.

The frontend is using the riotjs [http://www.riotjs.com](http://www.riotjs.com) library to do the ui. The riot compiler needs to be installed in case you are making changes to any of the tag files in www/tags
The riot compiler can be installed from npm `npm install riot -g`. Omit -g if you do not want to pollute your global npm install.


## Architecture
TODO Here will be a description of how the code is organized and the methodology behind.



## Deployment

To deploy, one can use Sup. The binary will be copied to `/usr/local/bin/ajournal`, the web data will be copied to `/var/www/ajournal`. During first deployment, remember to copy ajournal.service to `/usr/lib/systemd/system/ajournal.service`. Environment variables must ten be set up in `/etc/ajournal.conf`
