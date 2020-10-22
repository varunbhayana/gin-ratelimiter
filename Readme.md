Rate Limiting API in Golang for distributed sytems.
=========



Given
-----

You are given an external API endpoint which allows to check that given user-id does have the limit for the request or not.

Example HTTP calls

```
curl -X GET http://localhost:8888/rate
Request Header - user-id
```
