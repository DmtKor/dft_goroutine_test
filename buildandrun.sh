cd ./calc_src
go build -o ../build/calc
cd ../gen_data_src
go build -o ../build/gen_data
cd ..
./build/gen_data 10000 > ./data/data.txt
./build/calc < ./data/data.txt