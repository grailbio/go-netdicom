module github.com/apaladiychuk/go-netdicom

require (
	github.com/apaladiychuk/go-dicom v0.0.3
	github.com/stretchr/testify v1.2.2
)

replace github.com/ceph/go-ceph => ../../ceph/go-ceph

replace github.com/dghubble/oauth1 => ../../dghubble/oauth1

replace github.com/grailbio/base => ../base

replace github.com/grailbio/bigmachine => ../bigmachine

replace github.com/grailbio/bigslice => ../bigslice

replace github.com/grailbio/bio => ../bio

//replace github.com/apaladiychuk/go-dicom => ../go-dicom

replace github.com/apaladiychuk/go-netdicom => ../go-netdicom

replace github.com/grailbio/hts => ../hts

replace github.com/grailbio/ml => ../ml

replace github.com/grailbio/reflow => ../reflow

replace github.com/grailbio/testutil => ../testutil

replace github.com/grailbio/v23/factories/grail => ../v23/factories/grail

replace github.com/mijia/modelq => ../../mijia/modelq

replace github.com/unidoc/unidoc => ../../unidoc/unidoc

replace github.com/youtube/vitess => ../../youtube/vitess

replace v.io/x/ref/lib/flags/sitedefaults => ../../../v.io/x/ref/lib/flags/sitedefaults

replace github.com/golang/lint => ../../golang/lint

go 1.13
