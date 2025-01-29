# gingonic-tournament

This project build on gitlab environment, I just upload github for archive purpose. Not suppose to work as CICD automated

## Project Structure

gingonic-tournament/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── crud/
│   │   ├── leaderboard.go
│   │   ├── tournament.go
│   │   ├── users.go
│   ├── db/
│   │   ├── postgres.go
│   │   ├── redis.go
│   ├── router/
│   │   ├── leaderboard.go
│   │   ├── tournament.go
│   │   ├── users.go
├── model/
│   ├── leaderboard.go
│   ├── tournament.go
│   ├── user.go
├── service/
│   ├── leaderboard.go
│   ├── tournament.go
│   ├── users.go
├── validation/
│   ├── leaderboard.go
│   ├── tournament.go
│   ├── users.go
├── build/
│   ├── Dockerfile
├── deployments/
│   ├── docker-compose.yaml
├── scripts/
│   ├── build.sh
│   ├── test.sh
├── vendor/
├── test/
│   ├── integration/
│   │   ├── tournament_test.go
│   │   ├── user_test.go
│   ├── unit/
│       ├── leaderboard_test.go
│       ├── tournament_test.go
│       ├── user_test.go
├── go.mod
├── go.sum
├── .env
├── .gitlab-ci.yml

## jobs

done - 
Update endpoints related with users now need to give their level information, its not set by their money level anymore. Give an error when missed information given.

done - 
Add an enpoint for tournament.go. It contains ongoing tournament list and users that involved. User can join a tournament. When 3 user joined a tournament, leaderboard will be available via redis database. Tournament datas should stored on postgres. Add tournament model file and service file like before. Seperate them from endpoints. 

done - 
create also leaderboard service file /service/leaderboard.go and leaderboard model file /model/leaderboard.go with these changes. with these informaiton. Service files should include functions that uses validations when we need to update db. Model file should include struct and validate codes.

done - 
Change db structure with adding crud file. Carry all code to /internal/crud/tournament.go to related with tournament.go db sessions codes. /internal/db/postgres just include InitPostgres function to call postgres sessions. Give me the go files below:
/internal/crud/leaderboard.go
/internal/crud/users.go
/internal/crud/tournament.go
/internal/db/postgres.go
/internal/db/redis.go
and related other files

done - 
Update redis.go file with crud changes.Update redis.go file with crud changes.  There are still setLeaderboard here: 
done - 
This was wrong, crud file only includes with this structure: 
/internal/crud/leaderboard.go
/internal/crud/users.go
/internal/crud/tournament.go
/internal/db/postgres.go
/internal/db/redis.go
apply crud changes to redis again same with postgres.
Carry all code to /internal/crud/*.go to related with *.go db sessions codes. /internal/db/redis file just include Init functions to call db sessions.

done -
Carry all validation related files to from model files.
/validation/leaderboard.go 
/validation/users.go 
/validation/tournament.go 

done -
Update the leaderboards should set by users level degree. Add a function that when users joined with their money and level informaiton, joining a tournament decrease their money with 50. Finalizing a tournament gives rank one 200 money, rank two 100 money and rank three 50 money. Change the necessary files.

done -
structure has been changed please implement to new files:
/internal/db/postgres.go
/internal/db/redis.go
/internal/router/leaderboard.go
/internal/router/users.go
/internal/router/tournament.go

done -
model layer has been updated. It contains struct and validate codes.
/model/leaderboard.go
/model/users.go
/model/tournament.go

done -
service-layer has been added. It needs to contains functions like below from :
/service/leaderboard.go
/service/users.go
/service/tournament.go

done -
service files should include functions that uses validations when we need to update db for example:
func calculateLevel(money int) int {
	return money/100 + 1
}

done -
Also update following files according to changes:
model files
/cmd/app/main.go

## SQL Transactions:
### Knowledge resource:
https://stephenn.com/2023/08/mastering-database-transactions-in-go-strategies-and-best-practices/

Each transaction begins with a specific job and ends when all the tasks in the group successfully completed. If any of the tasks fail, the transaction fails. Therefore, a transaction has only two results: **success** or **failure**.

### Properties of Transaction

- **Atomicity:** The outcome of a transaction can either be completely successful or completely unsuccessful. The whole transaction must be rolled back if one part of it fails.
- **Consistency:** Transactions maintain integrity restrictions by moving the database from one valid state to another.
- **Isolation:** Concurrent transactions are isolated from one another, assuring the accuracy of the data.
- **Durability:** Once a transaction is committed, its modifications remain in effect even in the event of a system failure.