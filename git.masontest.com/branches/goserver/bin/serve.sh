#!/bin/sh

quit() {
    echo "---- OOPS ----"
    echo "Sorry, GoServer Starting failed:"
    echo ""
    cat logs/link.log
    exit 0
}

rm -rf runtime/* 
go run bin/builder.go > logs/link.log 2>&1
go build  -o bin/GoServer -i GoServer.go >> logs/link.log 2>&1

if cat logs/link.log | grep -e "fail|error|mis|undefined|no|too" -P -i; then 
    quit;
fi

killall -9 GoServer
./bin/GoServer
#rm logs/link.log

