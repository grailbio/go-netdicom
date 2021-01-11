// A sample program for issuing C-STORE or C-FIND to a remote server.
package main

import (
	"crypto/tls"
	"flag"
	"log"

	"github.com/BTsykaniuk/go-netdicom"
	"github.com/BTsykaniuk/go-netdicom/dimse"
	"github.com/BTsykaniuk/go-netdicom/sopclass"
	"github.com/apaladiychuk/go-dicom"
	"github.com/apaladiychuk/go-dicom/dicomtag"
)

var (
	serverFlag        = flag.String("server", "178.208.149.75:21113", "host:port of the remote application entity")
	storeFlag         = flag.String("store", "", "If set, issue C-STORE to copy this file to the remote server")
	aeTitleFlag       = flag.String("ae-title", "testclient", "AE title of the client")
	remoteAETitleFlag = flag.String("remote-ae-title", "SSLAIBUSSCP", "AE title of the server")
	findFlag          = flag.Bool("find", false, "Issue a C-FIND.")
	getFlag           = flag.Bool("get", false, "Issue a C-GET.")
	seriesFlag        = flag.String("series", "", "Study series UID to retrieve in C-{FIND,GET}.")
	studyFlag         = flag.String("study", "", "Study instance UID to retrieve in C-{FIND,GET}.")
)

func newServiceUser(sopClasses []string) *netdicom.ServiceUser {
	su, err := netdicom.NewServiceUser(netdicom.ServiceUserParams{
		CalledAETitle:  *remoteAETitleFlag,
		CallingAETitle: *aeTitleFlag,
		SOPClasses:     sopClasses})
	if err != nil {
		log.Panic(err)
	}

	cert, err := tls.LoadX509KeyPair("client_keystore.crt", "client_keystore.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	log.Printf("Connecting to %s", *serverFlag)
	su.Connect(*serverFlag, &config)
	return su
}

func cStore(inPath string) {
	su := newServiceUser(sopclass.StorageClasses)
	defer su.Release()
	dataset, err := dicom.ReadDataSetFromFile(inPath, dicom.ReadOptions{})
	if err != nil {
		log.Panicf("%s: %v", inPath, err)
	}
	err = su.CStore(dataset)
	if err != nil {
		log.Panicf("%s: cstore failed: %v", inPath, err)
	}
	log.Printf("C-STORE finished successfully")
}

func generateCFindElements() (netdicom.QRLevel, []*dicom.Element) {
	if *seriesFlag != "" {
		return netdicom.QRLevelStudy, []*dicom.Element{dicom.MustNewElement(dicomtag.SeriesInstanceUID, *seriesFlag)}
	}
	if *studyFlag != "" {
		return netdicom.QRLevelStudy, []*dicom.Element{dicom.MustNewElement(dicomtag.StudyInstanceUID, *studyFlag)}
	}
	args := []*dicom.Element{
		dicom.MustNewElement(dicomtag.AccessionNumber, "*"),
		dicom.MustNewElement(dicomtag.ReferringPhysicianName, "*"),
		dicom.MustNewElement(dicomtag.PatientName, "*"),
		dicom.MustNewElement(dicomtag.PatientID, "*"),
		dicom.MustNewElement(dicomtag.PatientBirthDate, "*"),
		dicom.MustNewElement(dicomtag.PatientSex, "*"),
		dicom.MustNewElement(dicomtag.StudyID, "1.2.40.0.13.1.251004474544013922953057532187552421649"),
		dicom.MustNewElement(dicomtag.RequestedProcedureDescription, "*"),
	}
	return netdicom.QRLevelPatient, args
}

func cGet() {
	su := newServiceUser(sopclass.QRGetClasses)
	defer su.Release()
	qrLevel, args := generateCFindElements()
	n := 0
	err := su.CGet(qrLevel, args,
		func(transferSyntaxUID, sopClassUID, sopInstanceUID string, data []byte) dimse.Status {
			log.Printf("%d: C-GET data; transfersyntax=%v, sopclass=%v, sopinstance=%v data %dB",
				n, transferSyntaxUID, sopClassUID, sopInstanceUID, len(data))
			n++
			return dimse.Success
		})
	log.Printf("C-GET finished: %v", err)
}

func cFind() {
	su := newServiceUser(sopclass.QRFindClasses)
	defer su.Release()
	qrLevel, args := generateCFindElements()
	for result := range su.CFind(qrLevel, args) {
		if result.Err != nil {
			log.Printf("C-FIND error: %v", result.Err)
			continue
		}
		log.Printf("Got response with %d elems", len(result.Elements))
		for _, elem := range result.Elements {
			log.Printf("Got elem: %v", elem.String())
		}
	}
}

func main() {
	flag.Parse()
	if *storeFlag != "" {
		cStore(*storeFlag)
	} else if *findFlag {
		cFind()
	} else if *getFlag {
		cGet()
	} else {
		log.Panic("Either -store, -get, or -find must be set")
	}
}
