This directory contains fuzz tests. It uses go-fuzz. Visit
https://github.com/dvyukov/go-fuzz and install two packages, go-fuzz-build and go-fuzz.
Then run:

```
go-fuzz-build github.com/yasushi-saito/go-netdicom/fuzzpdu
mkdir -p /tmp/fuzzpdu
go-fuzz -bin fuzzpdu-fuzz.zip -workdir /tmp/fuzzpdu
```
