# ArticleAPI
Prerequisite:
---

Docker setup:

https://docs.docker.com/docker-for-mac/install/

Go 1.12:

https://golang.org/dl/

# To start a standalone version of the app

make standalone
--
APP is accessable on port 8080 by default.
http://localhost:8080/api/healthcheck -Get
http://localhost:8080/api/articles - post
http://localhost:8080/api/articles/{id} - Get
http://localhost:8080/api/tags/{tagName}/{date} - Get- tagName is case sensitive

This creates an instance of mongodb from docker and waits for it to be up and running before starting the app

It builds the docker image of the app and start the app

It runs unit tests as part of douker build.

# To start App - locally with mongo docker
-----------------------------------------
make local
--

It runs a local docker instance in 'buildMongo' step. and starts the app using go run

App starts on post 4852 by default

# To see test coverage
----------------------

make coverage
--

To view the code coverage report
it also tests for race conditions

# To use an existing Mongo instance
-------------------------------

Modify the cmd\server\config\config.local.json to local mongo DB properties and run:

make start
--

This will connect the app to the existing local mongo instance.


# Features Implemented:
--------------------------

Dockerizied setup
-----------------
1. App can be started as a standalone app
2. multi- stage docker build is used, this ensures that the final docker image size, including all configurations is only 18.3MB
3. The application is tested and build before the image is generated. we create the executable and load just the executable and the configurations on a scratch image, to get a light weight image.
4. Kingpin is used to inject flags to the build. this helps in runningh the appliocation with different configurations.
5. Docker-compose is used to setup the environment for the app.
6. Mongo docker image is linked as a dependency
7. Docker healthchecks are in place
8. the app listens to the mongo docker healthcheck and waits this is ready to receive request. this ensures that the DB is always up and running before the app starts
9. Make files have simple steps to bring up the dockerizied standalone

MongoDB Integration:
------
1. Have used MongoDB, as it is very good with Nosql index based querying
2. It has an inbuilt caching mechanism which can cache close to 1/3rd of its size as a in-memory cache.   
3. The most recent data accessed is cached, not the query results itself. hence the queries are more reliable.
4. Complex aggregate queries are possible with mongo. the complete requirement complexity of the get /tagname/date/ is handled in mongo query.
5. This is very efficient as mongo allows indexing, even on array fields.

6. Indexing has been enabled.These are configurable, hence very reusable
7. DBproperties are derived from configurations - currently a local and a standlone configuration is in place.

Interface based Dependency Injection:
--------------
Concept of inheritance is used, to ensure extendability.
We can add another db interface say postgress or dynamo etc.
The new client has to implement the DBClient interface methods to qualify as a dbclient to be consumed by the application.

This interface based dependency injection is used throughout the app. this helps in injecting mocks for unit testing.


Graceful shutdown of server 
-------

1. a go routine is fired to listen for a termination event over a quit channel
2. After the server is started and is serving request the main thread is blocked trying to read data on a done channel.
3. when a os interruption is fed in the quit channel the go routine is unblocked and performs the shutdown related activites before shutting doen the server


Go Concepts:
---

1. Idiomatic error handling is used wherever possible.
2. Defer is used to execute items at the end of the function, like closing a file
3. slices and maps are used extensively throughout the app
4. File handling is used to read the app configurations
5. Json parser and bson parser used for decoding. 
6. Interface and type casting is used to generalize 
7. go routines used to achieve graceful shutdown
8. channels used to block and unblock the flow
9. go-chi router is used for apI routing
10. logrus is used for logging
11. middlewares are integrated with routers

REST API:
---------

A server uses a Delegate to delegate routes.

A Delegate interface is implemented by a delegate struct, which comprises of a Store.

delegate "Is a" Delegate

delegate "Has a" store

A store "Has a" DBClient.

any DBclient can be injected into the store as along as it implements the interface methods.

DBclient is initialized at startup. it establishes DB connectivity. 

MongoClient "Is a" DBClient

mongoClient has logic specific to mongodb connectivity and other mongo specific implementation syntax.

mongoClient is the only place where the connection,collection etc - items speciifc to mondo db can be accessed. these are encapsulated for the outside world.

router -> delegate -> store -> client -> db

1. Go-chi router and middlerware are implemented
2. Error handling is implemented
3. Go modules and vendor based dependencies are used.
4. App tested for race conditions

# Assumptions:
------------

1. Duplicate id insertions, would be treated as an error
2. date, tags and id fields are mandatory
3. An article can be created without title and/or body
4. Date should be of format YYYY-MM-DD
5. Date can be a future date.
6. Duplicate entries in tags will be persisted as is but filtered in view
7. Other additional unknow fields in the request will be ignored
8. There are no length restrictions on the fields
9. Tags are free text string – not predefined enums - they are case sensitive in search query





