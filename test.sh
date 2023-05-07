#!/bin/bash 

output_file="${1:-test.log}"

function run_test(){
    # remove db if it already exists
    rm -f test_*.txt

    ./main create --vault test_1.txt --masterpass mp 
    ./main list --vault test_1.txt --masterpass mp 
    ./main add --vault test_1.txt --masterpass mp --url ana.com --username ana --password ana
    ./main add --vault test_1.txt --masterpass mp --url ana2.com --username ana2 --password ana2
    ./main add --vault test_1.txt --masterpass mp --url ana.com --username ana3 --password ana3
    ./main list --vault test_1.txt --masterpass mp 
    ./main delete --vault test_1.txt --masterpass mp --url ana.com --username ana 
    ./main delete --vault test_1.txt --masterpass mp --url ana2.com --username ana2 
    ./main delete --vault test_1.txt --masterpass mp --url ana.com --username ana3 
    ./main delete --vault test_1.txt --masterpass mp --url ana.com --username ana3 
    ./main list --vault test_1.txt --masterpass mp

}

cd main
set -xe;

run_test > "$output_file" 2>&1


set +x;
if diff "$output_file" expected.log ; then 
    echo "SUCCESS"
fi

