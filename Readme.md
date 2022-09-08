# Project Details

This is just a learning project done to better learn Go while building a simple web application

The frontend is extremely simple because I had no interest in the scope of this POC on improving it.

## To run

- `run.sh` to build app
- `cd docker` and `docker-compose --f ./postgres-docker.yml up` to run postgres docker image
- `soda migrate up` to run migrations to populate database

## Built using

- Go version 1.18
- [chi router](github.com/go-chi/chi) for routing
- [scs](github.com/alexedwards/scs) for session management
- [nosurf](github.com/justinas/nosurf) for CSRF protection
- `POSTGREsql` as database
- [soda](https://gobuffalo.io/documentation/database/soda/) for migrations