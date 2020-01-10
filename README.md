# Echo + Storm barebones API

A simple API project template using Golang, Echo & Storm.

To create run:

```bash
git clone https://github.com/MBeliou/go-api-template.git
```
or manually through the github website.

To run:
```bash
    # This will run on the default port 1323 with mydb.db 
    # as the database file
    go run .

    # Change the port and the databse file
    go run --port=8080 --database=otherdb.db
```

## Echo

Find out more about [Echo](https://echo.labstack.com/), the web framework used here.

## Storm

Find out more about  [Storm](https://github.com/asdine/storm#options), a toolkit for [BoltDB](https://github.com/etcd-io/bbolt) which is a key/value store used as our database here.

___
## TODO
* [ ] Add Tests
* [ ] Better handling of the JWT secret signing key
* [ ] Use Transactions where needed
* [ ] Make a frontend