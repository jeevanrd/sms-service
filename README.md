# sms-service
SMS Service

Assumptions
1. Postgres db server up and running
2. Redis server up and running

Before running tests.
1)clear redis cache

redis-cli -r 1 flushall

2) Restore the dump into your own postgresql server.
psql -U postgres -f testdatadump.sql

To run test suite
go test tests
