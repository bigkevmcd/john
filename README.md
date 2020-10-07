# John mailets for Go

This is a simple SMTP Server that can have a chain of mail processors inserted.

```
$ echo "This is the message body and contains the message" | mailx -v -r \
  "someone@example.com" -s "This is the subject" -S smtp="localhost:2525" \
  testing@example.com
```
