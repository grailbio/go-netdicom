package sopclass

import (
	"github.com/apaladiychuk/go-dicom/dicomuid"
)

// DICOM SOP UID listing.
//
// https://www.dicomlibrary.com/dicom/sop/
//
// Translated from sop_class.py in pynetdicom3; https://github.com/pydicom/pynetdicom3

func standardUID(uid string) string {
	return dicomuid.MustLookup(uid).UID
}

// VerificationClasses is for issuing C-ECHO
var VerificationClasses = []string{
	standardUID("1.2.840.10008.1.1"),
}

// StorageClasses for issuing C-STORE requests.
var StorageClasses = []string{
	standardUID("1.2.840.10008.5.1.1.27"),
	standardUID("1.2.840.10008.5.1.1.29"),
	standardUID("1.2.840.10008.5.1.1.30"),
	standardUID("1.2.840.10008.5.1.4.1.1.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.1.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.1.1.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.1.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.1.2.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.1.3"),
	standardUID("1.2.840.10008.5.1.4.1.1.1.3.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.10"),
	standardUID("1.2.840.10008.5.1.4.1.1.104.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.104.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.11"),
	standardUID("1.2.840.10008.5.1.4.1.1.11.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.11.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.11.3"),
	standardUID("1.2.840.10008.5.1.4.1.1.11.4"),
	standardUID("1.2.840.10008.5.1.4.1.1.11.5"),
	standardUID("1.2.840.10008.5.1.4.1.1.12.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.12.1.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.12.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.12.2.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.12.3"),
	standardUID("1.2.840.10008.5.1.4.1.1.128"),
	standardUID("1.2.840.10008.5.1.4.1.1.129"),
	standardUID("1.2.840.10008.5.1.4.1.1.13.1.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.13.1.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.13.1.3"),
	standardUID("1.2.840.10008.5.1.4.1.1.130"),
	standardUID("1.2.840.10008.5.1.4.1.1.131"),
	standardUID("1.2.840.10008.5.1.4.1.1.14.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.14.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.2.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.20"),
	standardUID("1.2.840.10008.5.1.4.1.1.3"),
	standardUID("1.2.840.10008.5.1.4.1.1.3.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.4"),
	standardUID("1.2.840.10008.5.1.4.1.1.4.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.4.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.4.3"),
	standardUID("1.2.840.10008.5.1.4.1.1.481.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.481.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.481.3"),
	standardUID("1.2.840.10008.5.1.4.1.1.481.4"),
	standardUID("1.2.840.10008.5.1.4.1.1.481.5"),
	standardUID("1.2.840.10008.5.1.4.1.1.481.6"),
	standardUID("1.2.840.10008.5.1.4.1.1.481.7"),
	standardUID("1.2.840.10008.5.1.4.1.1.481.8"),
	standardUID("1.2.840.10008.5.1.4.1.1.481.9"),
	standardUID("1.2.840.10008.5.1.4.1.1.5"),
	standardUID("1.2.840.10008.5.1.4.1.1.6"),
	standardUID("1.2.840.10008.5.1.4.1.1.6.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.6.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.66"),
	standardUID("1.2.840.10008.5.1.4.1.1.66.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.66.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.66.3"),
	standardUID("1.2.840.10008.5.1.4.1.1.66.4"),
	standardUID("1.2.840.10008.5.1.4.1.1.66.5"),
	standardUID("1.2.840.10008.5.1.4.1.1.67"),
	standardUID("1.2.840.10008.5.1.4.1.1.68.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.68.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.7"),
	standardUID("1.2.840.10008.5.1.4.1.1.7.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.7.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.7.3"),
	standardUID("1.2.840.10008.5.1.4.1.1.7.4"),
	standardUID("1.2.840.10008.5.1.4.1.1.77.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.77.1.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.77.1.1.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.77.1.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.77.1.2.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.77.1.3"),
	standardUID("1.2.840.10008.5.1.4.1.1.77.1.4"),
	standardUID("1.2.840.10008.5.1.4.1.1.77.1.4.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.77.1.5.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.77.1.5.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.77.1.5.3"),
	standardUID("1.2.840.10008.5.1.4.1.1.77.1.5.4"),
	standardUID("1.2.840.10008.5.1.4.1.1.77.1.6"),
	standardUID("1.2.840.10008.5.1.4.1.1.77.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.78.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.78.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.78.3"),
	standardUID("1.2.840.10008.5.1.4.1.1.78.4"),
	standardUID("1.2.840.10008.5.1.4.1.1.78.5"),
	standardUID("1.2.840.10008.5.1.4.1.1.78.6"),
	standardUID("1.2.840.10008.5.1.4.1.1.78.7"),
	standardUID("1.2.840.10008.5.1.4.1.1.78.8"),
	standardUID("1.2.840.10008.5.1.4.1.1.79.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.8"),
	standardUID("1.2.840.10008.5.1.4.1.1.80.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.81.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.88.11"),
	standardUID("1.2.840.10008.5.1.4.1.1.88.22"),
	standardUID("1.2.840.10008.5.1.4.1.1.88.33"),
	standardUID("1.2.840.10008.5.1.4.1.1.88.34"),
	standardUID("1.2.840.10008.5.1.4.1.1.88.40"),
	standardUID("1.2.840.10008.5.1.4.1.1.88.50"),
	standardUID("1.2.840.10008.5.1.4.1.1.88.59"),
	standardUID("1.2.840.10008.5.1.4.1.1.88.65"),
	standardUID("1.2.840.10008.5.1.4.1.1.88.67"),
	standardUID("1.2.840.10008.5.1.4.1.1.88.69"),
	standardUID("1.2.840.10008.5.1.4.1.1.88.70"),
	standardUID("1.2.840.10008.5.1.4.1.1.9"),
	standardUID("1.2.840.10008.5.1.4.1.1.9.1.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.9.1.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.9.1.3"),
	standardUID("1.2.840.10008.5.1.4.1.1.9.2.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.9.3.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.9.4.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.9.4.2"),
	standardUID("1.2.840.10008.5.1.4.1.1.9.5.1"),
	standardUID("1.2.840.10008.5.1.4.1.1.9.6.1"),
	standardUID("1.2.840.10008.5.1.4.34.7"),
	standardUID("1.2.840.10008.5.1.4.43.1"),
	standardUID("1.2.840.10008.5.1.4.44.1"),
	standardUID("1.2.840.10008.5.1.4.45.1"),
}

// QRFindClasses is for issuing C-FIND requests.
var QRFindClasses = []string{
	standardUID("1.2.840.10008.5.1.4.1.2.1.1"),
	standardUID("1.2.840.10008.5.1.4.1.2.2.1"),
	standardUID("1.2.840.10008.5.1.4.1.2.3.1"),
	standardUID("1.2.840.10008.5.1.4.31")}

// QRMoveClasses is for issuing C-MOVE requests.
var QRMoveClasses = []string{
	standardUID("1.2.840.10008.5.1.4.1.2.1.2"),
	standardUID("1.2.840.10008.5.1.4.1.2.2.2"),
	standardUID("1.2.840.10008.5.1.4.1.2.3.2")}

// QRGetClasses is for issuing C-GET requests.
var QRGetClasses = append([]string{
	standardUID("1.2.840.10008.5.1.4.1.2.1.3"),
	standardUID("1.2.840.10008.5.1.4.1.2.2.3"),
	standardUID("1.2.840.10008.5.1.4.1.2.3.3")},
	StorageClasses...)
