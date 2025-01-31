#!/bin/bash

if [ -z "${TESTPATH}" ]
then
    TESTPATH=$(dirname $(readlink -e "$0"))
    export TESTPATH="${TESTPATH}/tests"
fi

declare -a TESTCASES

for test_file in $(ls "${TESTPATH}"); do
    if [[ $test_file == *.g4th ]]
    then
        TESTCASES+=("$test_file")
    fi
done

record() {
    # TODO complete this function
    echo "record function invoked"
    printf '%s\n' "${TESTCASES[@]}"
}

run() {
    # TODO complete this function
    echo "run function invoked"
    printf '%s\n' "${TESTCASES[@]}"
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
