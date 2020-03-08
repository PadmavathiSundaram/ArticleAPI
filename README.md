# ArticleAPI
Prerequisite:
Docker setup
https://docs.docker.com/docker-for-mac/install/

# To start a standalone version of the app

make standalone

this creates an instance of mongodb from docker and waits for it to be up and running before starting the app

It builds the docker image of the app and start the app

It runs unit tests as part of douker build.

# To start App - locally with mongo docker
-----------------------------------------
make local

It runs a local docker instance in 'buildMongo' step. and starts the app using go run

# To see test coverage
----------------------
make coverage

to view the code coverage report
it also tests for race conditions

#To use existing Mongo instance
-------------------------------

Modify the cmd\server\config\config.local.json to local mongo DB properties and run:

make start

This will connect the app to the existing local mongo instance.

