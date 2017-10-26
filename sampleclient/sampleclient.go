// A sample program for sending a DICOM file to a remote provider using C-STORE protocol.
package main

import (
	"flag"

	"github.com/grailbio/go-dicom"
	"github.com/grailbio/go-netdicom"
	"github.com/grailbio/go-netdicom/sopclass"
	"v.io/x/lib/vlog"
)

var (
	serverFlag        = flag.String("server", "localhost:10000", "host:port of the remote application entity")
	storeFlag         = flag.String("store", "", "If set, issue C-STORE to copy this file to the remote server")
	aeTitleFlag       = flag.String("ae-title", "testclient", "AE title of the client")
	remoteAETitleFlag = flag.String("remote-ae-title", "testserver", "AE title of the server")
	findFlag          = flag.String("find", "", "blah")
)

func newServiceUser(sopClasses []string) *netdicom.ServiceUser {
	su, err := netdicom.NewServiceUser(netdicom.ServiceUserParams{
		CalledAETitle:  *aeTitleFlag,
		CallingAETitle: *remoteAETitleFlag,
		SOPClasses:     sopClasses})
	if err != nil {
		vlog.Fatal(err)
	}
	defer su.Release()
	vlog.Infof("Connecting to %s", *serverFlag)
	su.Connect(*serverFlag)
	return su
}

func cStore(inPath string) {
	su := newServiceUser(sopclass.StorageClasses)
	defer su.Release()
	dataset, err := dicom.ReadDataSetFromFile(inPath, dicom.ReadOptions{})
	if err != nil {
		vlog.Fatalf("%s: %v", inPath, err)
	}
	err = su.CStore(dataset)
	if err != nil {
		vlog.Fatalf("%s: cstore failed: %v", inPath, err)
	}
	vlog.Infof("C-STORE finished successfully")
}

func cFind(argStr string) {
	su := newServiceUser(sopclass.StorageClasses)
	defer su.Release()
	args := []*dicom.Element{
		dicom.MustNewElement(dicom.TagSpecificCharacterSet, "ISO_IR 100"),
		dicom.MustNewElement(dicom.TagAccessionNumber, ""),
		dicom.MustNewElement(dicom.TagReferringPhysicianName, ""),
		dicom.MustNewElement(dicom.TagPatientName, ""),
		dicom.MustNewElement(dicom.TagPatientID, ""),
		dicom.MustNewElement(dicom.TagPatientBirthDate, ""),
		dicom.MustNewElement(dicom.TagPatientSex, ""),
		dicom.MustNewElement(dicom.TagStudyInstanceUID, ""),
		dicom.MustNewElement(dicom.TagRequestedProcedureDescription, ""),
		dicom.MustNewElement(dicom.TagScheduledProcedureStepSequence,
			dicom.MustNewElement(dicom.TagItem,
				dicom.MustNewElement(dicom.TagModality, ""),
				dicom.MustNewElement(dicom.TagScheduledProcedureStepStartDate, ""),
				dicom.MustNewElement(dicom.TagScheduledProcedureStepStartTime, ""),
				dicom.MustNewElement(dicom.TagScheduledPerformingPhysicianName, ""),
				dicom.MustNewElement(dicom.TagScheduledProcedureStepStatus, ""))),
	}
	for result := range su.CFind(netdicom.QRLevelStudy, args) {
		if result.Err != nil {
			vlog.Errorf("C-FIND error: %v", result.Err)
			continue
		}
		vlog.Errorf("Got response with %d elems", len(result.Elements))
		for _, elem := range result.Elements {
			vlog.Errorf("Got elem: %v", elem.String())
		}
	}
}

func main() {
	flag.Parse()
	vlog.ConfigureLibraryLoggerFromFlags()

	if *storeFlag != "" {
		cStore(*storeFlag)
	} else if *findFlag != "" {
		cFind(*findFlag)
	} else {
		vlog.Fatal("Either -store or -find must be set")
	}
}
