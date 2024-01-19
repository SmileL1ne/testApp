# Test Application

This is a test application which is a little CRUD service for retrieving user information, enriching with additional information from external API's and saving into PostgreSQL database. I've used clean architecture and RESTful endpoints in this project. 

## Usage

1. Download this project
2. Download all required dependecies
3. Create '.env' file and fill it by given '.env.example' file
4. Run app:
```console
    $ go run ./cmd/app
```

## Endpoints

- **GET: /users** - retrieving all users in json format
> For **pagination** inlcude **'page'** query parameter for page number and **'pageSize'** for number of users per page

> For **sorting** include **'sort'** query parameter with next available parameters - 'alphabeticallyByName', 'alphabeticallyBySurname', 'byAge', 'byNationality'

- **POST: /users** - adding new user info
- **PUT: /users/:id** - updating user info (only that user initially put) by id
- **DELETE: /users/:id** - deleting user by id