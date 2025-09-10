# gopi

this project was made with GORM and postgres
postgres is used because of scalability and plugin capabilities which makes it awesome choice for projects that needs scalability
hashing strategy is argon2 because of speed efficiency
swagger is also available in `http://localhost:8000/swagger/index.html`

## env management
for this project `.env` is available for demo but you can configure env-vault to download the `.env` for you
but i recommend using env-vault for managing the environment variables of the project as the details are available in the [Link](https://www.dotenv.org/)

## How can I test the Endpoints?
there is a JSON file named [go Auth.postman_collection](https://github.com/jexroid/gopi/blob/main/go%20Auth.postman_collection.json) that is used for testing the endpoints which are written with small test cases

## Testing
tests are available and can be use with Makefile also project will be tested using CI/CD features.

## Docker image
the Docker image of this file is available at [dockerhub](https://hub.docker.com/repository/docker/jextoid/gauth-gopi/general)

## Security features:
 - HTTP headers
 - CORS
 - XSS protection
 - NoSniff (Don't let anyone sniff the credentials endpoints)
 - CSRF (Not implemented because of lack of frontend services)

   also remember to change the GO_ENV to production on production level. (you can make this automatic in env-vault)
