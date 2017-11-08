[![GoDoc](https://godoc.org/github.com/grailbio/go-netdicom?status.svg)](https://godoc.org/github.com/grailbio/go-netdicom) [![Build Status](https://travis-ci.org/grailbio/go-netdicom.svg?branch=master)](https://travis-ci.org/grailbio/go-netdicom.svg?branch=master)

# Golang implementation of DICOM network protocol.

See doc.go for (incomplete) documentation.  See storeclient and storeserver for
examples.

Inspired by https://github.com/pydicom/pynetdicom3.

Status as of 2017-10-02:

- C-STORE, C-FIND, C-GET work, both for the client and the server. Look at
  sampleclient, sampleserver, or e2e_test.go for examples.  In general, the
  server (provider)-side code is better tested than the client-side code.

- Compatibility has been tested against pynetdicom and Osirix MD.

TODO:

- Documentation.

- Better SSL support.

- Implement the rest of DIMSE protocols, in particular C-MOVE on the client
  side, and N-* commands.

- Better message validation.

- Remove the "limit" param from the Decoder, and rely on io.EOF detection instead.
