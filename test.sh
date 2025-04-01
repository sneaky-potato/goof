#!/bin/bash

set -e

if [ -z "${TESTPATH}" ]
then
    TESTPATH=$(dirname $(readlink -e "$0"))
    export TESTPATH="${TESTPATH}/tests"
fi

declare -a TESTCASES
declare -a TESTCASESOUT

for test_file in $(ls "${TESTPATH}"); do
    if [[ $test_file == *.goof ]]
    then
        TESTCASES+=("$test_file")
    elif [[ $test_file == *.goof.txt ]]
    then
        TESTCASESOUT+=("$test_file")
    fi
done

record() {
    for test_file in "${TESTCASES[@]}"
    do
        echo "recording file $test_file"
        go run main.go "tests/$test_file"
        ./output > "tests/$test_file.txt"
    done
}

fail() {
    echo "test failed $1"
    rm -rf "tests/*.tmp"
    exit 2
}

run() {
    LEN=${#TESTCASESOUT[@]}
    if [ "${#TESTCASESOUT[@]}" -eq "0" ]; then
        echo "no text files in test directory, please use --record before using this"
        exit 2
    fi

    for test_file in "${TESTCASES[@]}"
    do
        echo "testing file $test_file"
        go run main.go "tests/$test_file"
        ./output > "tests/$test_file.tmp"
        cmp --silent "tests/$test_file.tmp" "tests/$test_file.txt" || fail "$test_file"
        rm -rf "tests/$test_file.tmp"
        echo "PASS $test_file"
    done
}

show_help() {
    echo "usage: bash test.sh [OPTIONS]"
    echo "OPTIONS:"
    echo "    --record: record test cases present in /tests directory"
    echo "    --run:    run test cases present in /tests directory"
}

case "$1" in
    --record)
        record
        ;;
    --run)
        run
        ;;
    *)
        show_help
        ;;
esac
