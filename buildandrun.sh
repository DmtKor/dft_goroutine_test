if [ $# -lt 1 ]; then
    echo "Usage:" $0 "<data size>"
else
    if [ $1 -lt 100 ]; then 
        echo "Data size should be greater or equal to 100"
    else
        cd ./calc_src
        go build -o ../build/calc
        cd ../gen_data_src
        go build -o ../build/gen_data
        cd ..
        ./build/gen_data $1 > ./data/data.txt
        ./build/calc < ./data/data.txt
    fi
fi