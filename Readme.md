Rate Limiting API in Golang for distributed sytems.
=========





Given
-----

You are given an external API endpoint which allows to check that given user-id does have the limit for the request on the given application or not.
Appliaction is also deployed on ec2 with redis-elastic cache.


Example HTTP calls

```
curl -X GET http://localhost:8888/rate
Request Headers:
1. user-id:<userid>
2. application-id:<application-id>

```
Run
----
- Clone the application
- Create conf.json file from copying the conf.sample.json in rate-limiting/conf/ folder
- Can specify the values of paramters mentioned in conf.json like port,redis credentials and number of request allowed in a minute and a hour for a user.
- Now run go build . in rate-limitng project directory.
- Once after build u will see binary file rate-limiting in the current directory
- Run ./rate-limiting and hit enter. Will se logs of running the server on the mentioned Port.

Design a rate limiting application or API.
#### Requirements

- The API **must** support multiple users and multiple application support.
- If the HTTP status **200** means user is allowed to make call on respective application, http status should **429** when user exceeds the request for the given application.
- The API should be stateless so that it can be horizontaly scaled up easily.
- API should support sliding window approach ex user should have limit on request in a minute as well as an hour. All the parameter like number of requests should configurable from env file.
- Deploy to AWS with using elastic-cache


