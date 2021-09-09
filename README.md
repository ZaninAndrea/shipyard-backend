# A generic go server [![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

The server can be used as a baseline to develop your own or even as is. The server handles authentication with JWT tokens and supports the CRUD operations for users, the data stored can be an arbitrary JSON. 

## How to use

To authenticate the requests put the authentication token in the "Authorization" header like this: "Bearer your-authentication-token".

### Server setup

All data is stored on a MongoDB database, so you should have one running, an easy to use managed MongoDB is offered by [Mongo Atlas](https://www.mongodb.com/cloud/atlas) (it has a free forever tier).

The server configuration is provided through environment variables, you can either set the environment variables or create a `.env` file. The variables that should be set are:

-   `CONNECTION_URI`: The connection uri to the mongo db
-   `JWT_SECRET`: The secret used to generate the tokens (set it to a strong password of your choice)
-   `PORT`: Port on which the server will be listening (8080 by default)

To launch the server run the following command:

```
go run .
```

Or build and run the build server on the platform of you choosing

```
go build
```

### Routes

-   `POST /login` Pass email and password to receive an authentication token

-   `POST /user` Pass email and password in the url query to register a new user, an authentication token will be returned
-   `GET /user` Returns the data associated with the authenticated user
-   `PUT /user` Pass a JSON payload to update the data associated with the authenticated user
-   `DELETE /user` Delete the authenticated user
