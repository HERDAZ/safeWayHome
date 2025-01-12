# Safe Way Home
A geolocalisation app to make sure you get home safe.

## Install

Go version ```go 1.23.2 linux/amd64``` is the one used on the server.
The libraries used are :

 - net/http
 - github.com/gin-gonic/gin
 - golang.org/x/crypto/bcrypt
 - math/big
 - github.com/go-sql-driver/mysql
 - database/sql
 - errors
 - slices
 - fmt
 - time
 - log

It will also need a database, the one used in the server is a mysql database version ```8.0.40-0ubuntu0.22.04.1```.
You can use the template in dbTemplate.sql to generate the db with the good structure. Make sure to create an account that can read/write to the db to be used by go, put it's credential in dbConnectTemplate.go, and rename dbConnectTemplate.go to dbConnect.go.

## Launching the server

A shell script is used to launch the server. It contains 2 lines, changing wich one is commented allows you to choose if there should be logs (stored in log.txt)
```
$./runApp.sh
```

## Testing the server

A shell script is used to test the API, and is just a collection of curl with parameters (use ./testApp.sh to see all options)
```
$./testApp.sh <test>
```

## Warnings

During developement, it was noticed that when using the app while connected to a restricive network would refrain you from connecting to the server. U can use roaming data to resolve the problem, or change the server port to something more common (80, 22, 666 (DOOM multiplayer), ... ). Just be aware that this change should also be made in the app.
