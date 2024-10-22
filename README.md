# Makerble Task

This project is created as a part of selection process for Makerble. It implements a program for doctors and receptionists and helps perform the following tasks:

- Login authentication for both doctors and receptionists.

- Using JWT Auth and Refresh Tokens to ensure other paths are protected from use by unauthorized entitites.

- Ensured stored passwords are hashed for added security.

- Registering a new patient.(can only be used by receptionists).

- Viewing and updating the data for registered patients (for both doctors and receptionists).

- Deleting a patient's data (for receptionists only)

- This exclusivity to usage is implemented via the JWT tokens.

- Implements Redis to increase the speed during data fetching of patients.

## Testing out the project

The project is deployed at [makerble-90sx.onrender.com](makerble-90sx.onrender.com) . 

[Here]() is the api documentation.

**Note** - Since the project uses the free tier of render it might spin down due to inactivity , so while using the deployed link the first response might take a while , please have patience and retry in a minute. Thank you for your understanding.

The project exposes the following endponts :

### Ping Route

- **`/ping`** 
  - **Method:** GET
  - **Description:** Checks whether the server is alive.

### Authentication Routes

- **`/auth/register`**
  - **Method:** POST
  - **Description:** Registers a new user (e.g., doctor, receptionist).

- **`/auth/login`**
  - **Method:** POST
  - **Description:** Logs in an existing user.

### Patient Routes

- **`/patient`**
  - **Method:** POST
  - **Description:** Creates a new patient record. Accessible only by receptionists.

- **`/patient/:id`**
  - **Method:** GET
  - **Description:** Retrieves the data of a specific patient by their ID. Accessible by both doctors and receptionists.

- **`/patients`**
  - **Method:** GET
  - **Description:** Fetches the data of all patients. Accessible by both doctors and receptionists.

- **`/patient/:id`**
  - **Method:** DELETE
  - **Description:** Deletes a specific patient by their ID. Accessible only by receptionists.

- **`/patient/:id`**
  - **Method:** PUT
  - **Description:** Updates the records of a specific patient by their ID. Accessible by both doctors and receptionists.

### Notes
- Patient Routes are only accessible via authenticated users.

## Setting up the project Locally

To set up the project paste the follwing commands in your terminal:

```bash
git clone https://github.com/Swetabh333/Makerble.git
cd Makerble
go mod tidy
```
This will install all the required dependencies for the project.

Next you have to set the environment for the project. your root directory you have to paste the following information:

```bash
export DSN_STRING="<your_postgres_connection_string>/<your_database_name>"
export JWT_SECRET="<your_jwt_secret>"
export JWT_REFRESH_SECRET="<your_jwt_refresh_secret>"
export REDIS_URL="<your_redis_connection_url>:<port no.>" 
export REDIS_PASSWORD="<your_redis_password>"
```
I have used it this way for deployment purpose as the godotenv library gives an error on not detecting an env file. 

If you want to use a .env file , create a .env file in the root directory

and put these there

```
DSN_STRING="<your_postgres_connection_string>/<your_database_name>"
JWT_SECRET="<your_jwt_secret>"
JWT_REFRESH_SECRET="<your_jwt_refresh_secret>"
REDIS_URL="<your_redis_connection_url>:<port no.>" 
REDIS_PASSWORD="<your_redis_password>"
```
put this piece of code in the very beginninig of init function after error declaration in app/main/main.go

```
var err error

if err = godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file")
    }
//Connect to the database and return db instance
	db, err = databases.ConnectToDatabase()
```
Do the above and then in your terminal run

```bash
go get github.com/joho/godotenv
```

Now you'll have to build the project with the following command in the root directory

```bash
go build -o bin/exe app/main/main.go
```

This will create an executable in your bin folder, which you can run using

**NOTE** : make sure no other process is running on port 8080

```bash
./bin/exe
```

**Your backend is now listening at port `8080`**.

## Setting up using docker
You can use docker compose to set this up as well
```bash
git clone https://github.com/Swetabh333/Makerble.git
cd Makerble
```

then run the docker build command

```bash
docker compose up --build
```

**Your backend is now listening at port `8080`**.