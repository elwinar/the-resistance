# The Resistance API

This API is a game server for _The Resistance_, a strategic card game.

## Documentation

### Authentification

There is a single idempotent handler handling authentification. If the user
doesn't exists yet, the handler will create it with the given password. If the
user already exists, the handler will check that the password is right. In both
case, the handler will return a JWT that can be used for authenticating
requests later.

```
POST /login

{
    "login": "user",
    "password": "password"
}
```

```
200 OK

{
    "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE0OTQ2MDI5NDgsInVzZXIiOiJlbHdpbmFyIn0=.xk_7Dz5wBhxNn_Eb08JVhSoXmIos74-A6bGBC5PK1B4="
}
```

*Then token must be passed as the `token` header for each subsequent request.*
It is omitted in these examples for the sake of brevity.

Additonally, a token can be verified using an additonal route than will return
if the token is valid (tokens have a lifetime of 5 minutes).

```
POST /authenticate
token: eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE0OTQ2MDI5NDgsInVzZXIiOiJlbHdpbmFyIn0=.xk_7Dz5wBhxNn_Eb08JVhSoXmIos74-A6bGBC5PK1B4=
```

```
200 OK

{
    "authenticated": true
}
```

### Game

Games are the primary entity of the API. A game is simply an ID, with some dates to indicated if the game started, finished, etc.

A game isn't started until the `started_at` attribute is non-nul, and once started it is finished when the `finished_at` attribute is non-nul.

#### List the games

```
GET /game
```

```
200 OK

[
    {
        "id": 1,
        "created_at": "2017-05-19T15:30:53.311021359+02:00",
        "started_at": null,
        "finished_at": null
    }
]
```

#### Show a game

```
GET /game/1
```

```
200 OK


{
    "id": 1,
    "created_at": "2017-05-19T15:30:53.311021359+02:00",
    "started_at": null,
    "finished_at": null
}
```
#### Create a new game.

```
POST /game
```

```
200 OK

{
    "game": 1
}
```
