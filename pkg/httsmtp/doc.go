package httsmtp

// Provides a pseudo-SMTP server over HTTP
//
// The body is read from an HTTP POST request.
// Headers are parsed from the HTTP Headers any header that starts with
// John-Mail- will have the prefix stripped and provided in the HTTP message.
//
// The envelope comes from John-Envelope-From and John-Envelope-To.
//
// For example:
//
// curl -d "Testing" \
//   -H "John-Envelope-From: testing@example.com" \
//   -H "John-Envelope-To: test@example.com" \
//   -H "John-Mail-From: testing@example.com" \
//   -H "John-Mail-To: test@example.com" \
//   -H "John-Mail-Test-Header: this is a test" \
//   http://localhost:8080/
