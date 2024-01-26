#!/bin/sh

export COCKROACHDB_URL='cockroachdb://cockroach:@localhost:26257?sslmode=disable&user=root'
migrate create -ext sql -dir database/migration/ -seq init_mg
migrate -path database/migration/ -database ${COCKROACHDB_URL} -verbose up