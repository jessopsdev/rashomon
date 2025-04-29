# rashomon

**This is NOT for Production usage. It has not been tested to the rigors of production and is intended only for educational purposes. Use at your own risk.**

rashomon is an implementation of "request-signing" authentication for HTTP APIs, inspired by [AWS Sig V4](https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-authenticating-requests.html) but with the key difference in that it uses asymmetric cryptography. You can read more [further below](#why-hmac-is-used-in-aws-sigv4-what-will-rashomon-use).

In request-signing, certain attributes of an HTTP request form a *canonical request*. Utilizing the end-users private key the canonical request is signed, and the signature is shared in the *Authorization* header.

The server, utilizing it's public key, re-creates the canonical request using the incoming HTTP request, confirms it with the signature, then does a few additional checks. This allows for:
- Identity verification
- Elimination of secrets in-transit
- Strong protection against re-play attacks
- Strong guarentee that major components of the canonical request have not been manipulated in transit

The key thing to note is this allows API authentication without requiring storage of secrets for the API implementor.

## Usage

```shell
go get github.com/jessopsdev/rashomon
```

All `rashomon` functions exposed in the public API take a "Canonical Request", defined as:

```go
type CanonicalRequest struct {
	Method      string
	URL         string
	KeyID       string
	Timestamp   string
	Nonce       string
	ContentType string
	Body        string
}
```

The struct has a method `String()` that generates the string version (delimited by`\n`) of the canonical request, which ultimately gets signed. This looks something like:

```
GET
localhost:9090/v1/a2?best=1
x-rmn-key-id: <identifier of the public key, e.g.  dad4ca4f-cd82-4d08-a948-0a5a1e5ad786>
x-rmn-timestamp: <ISO 8601, e.g. 	2025-04-20T23:20:40+09:00>
x-rmn-nonce: <random generated value, e.g.  s30a5a1e5ad78%6>
content-type: <string, empty on GET>
body: <string, empty on GET>
```

The API itself is fairly self-explanatory::

```go
Sign(privateKey *ed25519.PrivateKey, req *CanonicalRequest) []byte
VerifySignature(publicKey *ed25519.PublicKey, req *CanonicalRequest, signature *[]byte) bool 
```

That's it.

## Design Notes

### Yet another disclaimer; don't use this.

Great time to remind you I'm not a security or cryptography expert. Remember this is for educatioanl purposes- do not use this in production.

### Dependencies

I don't like having core/critical components with lots of random third-party dependencies. A goal for rashomon was to use only the Go standard library.

### Why do you use asymmetric cryptography instead of AWS SigV4's HMAC implementation?

Most of what I can find on the topic in way of blogs and tech talks speaks to AWS SigV4 using HMAC for [AWS's need for incredible scale](https://aws.amazon.com/awstv/watch/44e6f1abd16/). It is many factors cheaper on computation to take the request as-is and *just re-calculate* the HMAC & SHA-256 signature than to do asymmetric cryptography. You can also do derived keys or some other fancy stuff that's beyond me... It's actually all quite brilliant.

But it got me thinking... most companies don't need to scale their APIs to billion requests-per-second scale and if I'm a little more worried about the need to have a central repository of HMAC keys for each client than I am scaling. That sounds like a larger *risk* of what can go wrong...  In fact, I'd rather just generate cryptographically random API keys, and then use [bcrypt](https://en.wikipedia.org/wiki/Bcrypt) to store rather than store customer secrets.

So, I figured asymmetric cryptography will be better for the average API. Ed25519 is still *millisecond scale* generally and quite often used for JWT and other token implementations, so I think this is generally good enough for most use cases, and gets us the benefits up above without the concerns.

### Nonce & Timestamp

We want to ensure that requests can't be re-played; that is, hypothetically, if another party can see the request in transit, they can't just copy the request and repeat it.

To prevent this, all canonical requests must have a *nonce* and an *ISO 8601 timestamp* designated. After checking the signature, `rashomon` confirms that the timestamp is within +/-1 minute of our server check, and that the nonce value has not been seen for this key in the past 5 minutes.

Failure here should result in denial; this is built into the core library, so that it's not part of individual implementation.

## Real-Implementation Notes
 
In a real world situation, to distribute this, you'd likely want a service backed by cache that handles key registration, nonce registration, and key look ups. For nonce a utilizing a feature such as TTL in redit would be helpful, with a key like `key:nonce` is probably useful. DynamoDB could also serve this usecase quite well.

You could then create a middleware such as `net/http` that runs local to your APIs that handles request verification.

The correct logic would be:
1. Validate request
2. Perform key lookup; deny if key not found
3. Verify signature
4. Check timestamp is within +/-1 Minutes
5. Perform nonce lookup (recommended TTL of 5 minutes)

### Invalidating Keys

Invalidating keys means simply making sure look-ups on `x-rmn-key-id` fail. This would prevent any any signatures for this key, even *future* from working.
