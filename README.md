# John - mailets for Go

This is a simple SMTP Server that can have a chain of mail processors inserted, this is not a production tool, and was written for teaching purposes.

## Running John

```shell
$ go build ./cmd/john
$ ./john smtp
```

## Testing mail delivery

```shell
$ echo "This is the message body and contains the message" | mailx -v -r \
  "someone@example.com" -s "This is the subject" -S smtp="localhost:2525" \
  testing@example.com
```

## Running the tests

```shell
$ go test ./...
```
