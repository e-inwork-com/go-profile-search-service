# [e-inwork.com](https://e-inwork.com)

## Golang Profile Search Microservice
This microservice searches profile data and integrates with the following microservices:
- [Golang User Microservice](https://github.com/e-inwork-com/go-user-service)
- [Golang Profile Microservice](https://github.com/e-inwork-com/go-profile-service)
- [Golang Proifle Indexing Mircoservice](https://github.com/e-inwork-com/go-profile-indexing-service)

To run both of the microservices, follow the command below:
1. Install Docker
    - https://docs.docker.com/get-docker/
2. Git clone this repository to your localhost, and from the terminal run below command:
   ```
   git clone git@github.com:e-inwork-com/go-profile-search-service
   ```
3. Change the active folder to `go-user-search-service`:
   ```
   cd go-profile-search-service
   ```
4. Run Docker Compose:
   ```
   docker-compose -f docker-compose.local.yml up -d
   ```
5. Create a user in the User API with CURL command line:
    ```
    curl -d '{"email":"jon@doe.com", "password":"pa55word", "first_name": "Jon", "last_name": "Doe"}' -H "Content-Type: application/json" -X POST http://localhost:8000/service/users
    ```
6. Login to the User API:
   ```
   curl -d '{"email":"jon@doe.com", "password":"pa55word"}' -H "Content-Type: application/json" -X POST http://localhost:8000/service/users/authentication
   ```
7. You will get a token from the response login, and set it as a token variable for example like below:
   ```
   token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjhhY2NkNTUzLWIwZTgtNDYxNC1iOTY0LTA5MTYyODhkMmExOCIsImV4cCI6MTY3MjUyMTQ1M30.S-G5gGvetOrdQTLOw46SmEv-odQZ5cqqA1KtQm0XaL4
   ```
8. Create a profile for current user, you can use any image or use the image on the folder test:
   ```
   curl -F profile_name="Jon Doe" -F profile_picture=@api/test/images/profile.jpg -H "Authorization: Bearer $token"  -X POST http://localhost:8000/service/profiles
   ```
9. Get a list of profiles from the profile search endpoint:
   ```
   curl "http://localhost:8000/service/search/profiles?q=*:*"
   ```
10. Good luck!