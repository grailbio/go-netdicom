This directory contains fuzz tests. It uses go-fuzz. Visit
https://github.com/dvyukov/go-fuzz and install two packages, go-fuzz-build and go-fuzz.
Then run:

```
go-fuzz-build github.com/yasushi-saito/go-netdicom/fuzze2e
mkdir -p /tmp/fuzze2e
go-fuzz -procs 64 -bin fuzze2e-fuzz.zip -workdir /tmp/fuzze2e
```
