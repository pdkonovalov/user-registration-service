# User registration service.

## Email verification.

This endpoint send verification code to email.

```bash
curl --request POST 'http://localhost/api/v1/users/email/verify?email=pdkonovalov@gmail.com'
```

## Create new user.

Verification code can be used to make account.

```bash
curl --request POST 'http://localhost/api/v1/users/new' \
--header 'Content-Type: application/json' \
--data '{
    "Name":"Petr",
    "Username":"pdkonovalov",
    "Password":"12345",
    "Email":"pdkonovalov@gmail.com",
    "EmailCode":97537
}'
```

The password must be between 5 and 20 characters long.

## Change password.

New verification code from first endpoint can be used to change password.

```bash
curl --request POST 'http://localhost/api/v1/users/password/new' \
--header 'Content-Type: application/json' \
--data '{
    "Password":"12346",
    "Email":"pdkonovalov@gmail.com",
    "EmailCode":501137
}'
```

## Get pair of jwt tokens.

```bash
curl 'http://localhost/api/v1/users/jwt/new?email=pdkonovalov@gmail.com&password=12346'
```

```json
{
    "AccessToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpcCI6IjE3Mi4yMC4wLjE6NDIzNDQiLCJzdWIiOiJwZGtvbm92YWxvdkBnbWFpbC5jb20iLCJleHAiOjE3MzQyMzgwOTh9.mo_7xZ_kHdzBi0_rXbipWs5FJwHqliQdcR4YmsChX44jJeG-CQxNZcTqkPuBEsoRZJdSqX0JH_LM13iNKaNDMA",
    "RefreshToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJwZGtvbm92YWxvdkBnbWFpbC5jb20iLCJleHAiOjE3MzQzMjQzNzh9.oQgePe0NS7lwpQ2ssDtRPsddA_WjahiQA5dxCabb2yFA80lxo9KpUMS54R0nhEu7Ub8dvPze6SklhPsow7GSrA"
}
```

## Refresh pair of jwt tokens.

```bash
curl --request POST 'http://localhost/api/v1/users/jwt/refresh' \
--header 'Content-Type: application/json' \
--data '{
    "AccessToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpcCI6IjE3Mi4yMC4wLjE6NDIzNDQiLCJzdWIiOiJwZGtvbm92YWxvdkBnbWFpbC5jb20iLCJleHAiOjE3MzQyMzgwOTh9.mo_7xZ_kHdzBi0_rXbipWs5FJwHqliQdcR4YmsChX44jJeG-CQxNZcTqkPuBEsoRZJdSqX0JH_LM13iNKaNDMA",
    "RefreshToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJwZGtvbm92YWxvdkBnbWFpbC5jb20iLCJleHAiOjE3MzQzMjQzNzh9.oQgePe0NS7lwpQ2ssDtRPsddA_WjahiQA5dxCabb2yFA80lxo9KpUMS54R0nhEu7Ub8dvPze6SklhPsow7GSrA"
}'
```

```json
{
    "AccessToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpcCI6IjE3Mi4yMC4wLjE6Mzc3ODQiLCJzdWIiOiJwZGtvbm92YWxvdkBnbWFpbC5jb20iLCJleHAiOjE3MzQyMzgyMTJ9.1FWg0mGl_-Mjww_5v0keulTM4WDg2Of6_wxHmk6sdjG7OtvYdZCupvWqwNsTW1MhQhwERqQFaF8kAUfzGBiemw",
    "RefreshToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJwZGtvbm92YWxvdkBnbWFpbC5jb20iLCJleHAiOjE3MzQzMjQ0OTJ9.oz_C1oeCxPLqhDcV1uHdiGMDLtmzYCALwuS0GVbPaEbNjEmxiN5yTuThLQydGofOe1rzdSrmkrAYoHgFCK0HQQ"
}
```

If refresh request make from new ip service send email allert.
