build:
    go build -o bin/csvutil main.go

run_convert:
    go run main.go convert -o out samples/csv/csv_202204.csv samples/schema/202205_header.txt
