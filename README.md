# The Resistance API

This API is a game server for _The Resistance_, a strategic card game.

## Documentation

### Authentification

There is a single idempotent handler handling authentification. If the user
doesn't exists yet, the handler will create it with the given password. If the
user already exists, the handler will check that the password is right. In both
case, the handler will return a JWT that can be used for authenticating
requests later.

*Payload*

```
POST /login

{
	"login": "user",
	"password": "password"
}
```

*Response*

```
200 OK

{
	"token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE0OTQ2MDI5NDgsInVzZXIiOiJlbHdpbmFyIn0=.xk_7Dz5wBhxNn_Eb08JVhSoXmIos74-A6bGBC5PK1B4="
}
```
