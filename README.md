####JumpCloud Interview Assignment #####


*When launched, program starts and listens on port 8080.

#Endpoints: 

POST /hash
GET /stats
GET /shutdown


####### POST /hash #######

#Example request
curl -d "password=angryMonkey" -X POST http://localhost:8080/shutdown

#Note: When password is not present in body of request with value, response will be: "Error - 'password' not found!"

#Password is tased with SHA512 as hashing algorithm and a Base64 encoded string is returned. 

#Server will respond with password in ~5 seconds .

####### GET /stats #######

#Example request
curl -X GET http://localhost:8080/stats

#JSON object is returned with two k-v pairs.
--“total” - total number of requests will server is running
--“average” - average time in microseconds for all POST requests to /hash endpoint. (Time taken to revive, hash, encode, password)

#Note: Object values are reset when server is shutdown

####### GET /shutdown #######

#Example request#
curl -X GET http://localhost:8080/shutdown

#Server is immediately shutdown.
