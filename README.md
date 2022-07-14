# csvutil-go
Personal CSV utility commands.

Currently it has only one functionality to make it easy for you to translate CSV headers (columns) into user-defined headers using a simple schema file.

## Usage

```shell
csvutil template sample.csv
```
It generates a template schema file and prints to the console.

```shell
csvutil generate -o ./output/ sample.csv... schema.txt
```

`generate` command converts the specified CSV files using the schema and generates translated files into the `output` directory.