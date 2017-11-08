// A sample program for issuing C-STORE or C-FIND to a remote server.
package main

import (
	"flag"

	"github.com/grailbio/go-dicom"
	"github.com/grailbio/go-dicom/dicomtag"
	"github.com/grailbio/go-netdicom"
	"github.com/grailbio/go-netdicom/sopclass"
	"github.com/grailbio/go-netdicom/dimse"
	"v.io/x/lib/vlog"
)

var (
	serverFlag        = flag.String("server", "localhost:10000", "host:port of the remote application entity")
	storeFlag         = flag.String("store", "", "If set, issue C-STORE to copy this file to the remote server")
	aeTitleFlag       = flag.String("ae-title", "testclient", "AE title of the client")
	remoteAETitleFlag = flag.String("remote-ae-title", "testserver", "AE title of the server")
	findFlag          = flag.String("find", "", "If nonempty, issue a C-FIND.")
	getFlag           = flag.String("get", "", "If nonempty, issue a C-GET.")
)

func newServiceUser(sopClasses []string) *netdicom.ServiceUser {
	su, err := netdicom.NewServiceUser(netdicom.ServiceUserParams{
		CalledAETitle:  *aeTitleFlag,
		CallingAETitle: *remoteAETitleFlag,
		SOPClasses:     sopClasses})
	if err != nil {
		vlog.Fatal(err)
	}
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

func cGet(argStr string) {
	su := newServiceUser(sopclass.QRGetClasses)
	defer su.Release()
	args := []*dicom.Element{
		dicom.MustNewElement(dicomtag.PatientID, "PAT004"),
	}

	n := 0
	err := su.CGet(netdicom.QRLevelPatient,
		args,
		func(transferSyntaxUID, sopClassUID, sopInstanceUID string, data []byte) dimse.Status {
			vlog.Infof("%d: C-GET data; transfersyntax=%v, sopclass=%v, sopinstance=%v data %dB",
				n, transferSyntaxUID, sopClassUID, sopInstanceUID, len(data))
			n++
			return dimse.Success
		})
	vlog.Infof("C-GET finished: %v", err)
}

func cFind(argStr string) {
	su := newServiceUser(sopclass.QRFindClasses)
	defer su.Release()
	args := []*dicom.Element{
		dicom.MustNewElement(dicomtag.SpecificCharacterSet, "ISO_IR 100"),
		dicom.MustNewElement(dicomtag.AccessionNumber, ""),
		dicom.MustNewElement(dicomtag.ReferringPhysicianName, ""),
		dicom.MustNewElement(dicomtag.PatientName, ""),
		dicom.MustNewElement(dicomtag.PatientID, ""),
		dicom.MustNewElement(dicomtag.PatientBirthDate, ""),
		dicom.MustNewElement(dicomtag.PatientSex, ""),
		dicom.MustNewElement(dicomtag.StudyInstanceUID, ""),
		dicom.MustNewElement(dicomtag.RequestedProcedureDescription, ""),
		dicom.MustNewElement(dicomtag.ScheduledProcedureStepSequence,
			dicom.MustNewElement(dicomtag.Item,
				dicom.MustNewElement(dicomtag.Modality, ""),
				dicom.MustNewElement(dicomtag.ScheduledProcedureStepStartDate, ""),
				dicom.MustNewElement(dicomtag.ScheduledProcedureStepStartTime, ""),
				dicom.MustNewElement(dicomtag.ScheduledPerformingPhysicianName, ""),
				dicom.MustNewElement(dicomtag.ScheduledProcedureStepStatus, ""))),
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
	} else if *getFlag != "" {
		cGet(*getFlag)
	} else {
		vlog.Fatal("Either -store or -find must be set")
	}
}
