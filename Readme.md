Rate Limiting API in Golang for distributed sytems.
=========



Given
-----

You are given an external API endpoint which allows to check that given user-id does have the limit for the request on the given application or not.

Example HTTP calls

```
curl -X GET http://localhost:8888/rate
Request Headers:
1. user-id:<userid>
2. application-id:<application-id>

```
Task
----

Design a rate limiting application or API.
#### Requirements

- The API **must** support multiple users and multiple application support.
- If the HTTP status **200** means user is allowed to make call on respective application, http status should **429** when user exceeds the request for the given application.
