# John - mailets for Go

This is a simple SMTP Server that can have a chain of mail processors inserted, this is not a production tool, and was written for teaching purposes.

## Running John

```shell
$ go build ./cmd/john
$ ./john smtp --init-maildir
```

Mail will be delivered to `./tmp` in [Maildir](https://en.wikipedia.org/wiki/Maildir) format.

## Can't use SMTP?

It might be simpler to use HTTP to send emails.

```shell
$ ./john http --init-maildir
```

This defaults to serving on :8080

You can send emails via `curl`

```shell
#!/bin/sh
curl -d "Testing" \
  -H "John-Envelope-From: testing@example.com" \
  -H "John-Envelope-To: test@example.com" \
  -H "John-Mail-From: testing@example.com" \
  -H "John-Mail-To: test@example.com" \
  -H "John-Mail-Test-Header: this is a test" \
  http://localhost:8080/
```

Again mail will be delivered to `./tmp` in [Maildir](https://en.wikipedia.org/wiki/Maildir) format.

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
