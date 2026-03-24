# SAT JADE

## Instructions to set up your local environment
 - have Postres 17.6 installed
 - Go 1.24.1 
 - set up your env variables in the .env file as shown in the sample.env file
 - run `go run ./cmd/migration/seed` in order to create an admin user in the database (set the password and email in the ADMIN_PASSWORD and ADMIN_EMAIL variables inside your .env file)
