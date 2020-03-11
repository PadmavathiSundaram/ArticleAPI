# ArticleAPI
Prerequisite:
Docker setup
https://docs.docker.com/docker-for-mac/install/
Go 1.12
https://golang.org/dl/

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


#####Features Implemented:
--------------------------

Dockerizied setup
    App can be started as a standalone app
    Docker-compose is used to setup the environment for the app to run using mongo docker image as a dependency
    docker healthchecks in place
    make files have simple steps to bring up the dockerizied standalone

MongoDB Integration:
    have used MongoDB as it is very good with Nosql index based querying
    It has an inbuilt caching mechanism which can cache close to 1/3rd of its size in memory.   
    The most recent data accessed is cached, not the query results itself. hence the queries are more reliable.
    Complex aggregate queries are possible with mongo. the complete requirement complexity of th eget tagname/date/ query is handled in mongo query.
    this is very efficient as mongo allows indexing, even on array fields.

    Indexing has been enabled.
    these are configurable, hence very reusable

    DBproperties are derived from configurations - currently a local and a standlone configuration is in place.
Interface based Dependency injection:
    Concept of inheritance is used, to ensure extendability.
    We can add another db interface say postgress or dynamo etc.
    The new client has to implement the DBClient interface methods to qualify as a dbclient to be consumed by the application.

    this interface based dependency injection is used throughout the app. this helps in injecting mocks for unit testing 



Graceful shutdown of server 
