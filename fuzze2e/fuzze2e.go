package fuzze2e

import (
	"flag"
	"log"
	"net"

	"github.com/apaladiychuk/go-dicom"
	"github.com/grailbio/go-netdicom"
	"github.com/apaladiychuk/go-netdicom/dimse"
	"github.com/apaladiychuk/go-netdicom/sopclass"
)

func startServer(faults netdicom.FaultInjector) net.Listener {
	netdicom.SetProviderFaultInjector(faults)
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Panic(err)
	}
	go func() {
		// TODO(saito) test w/ small PDU.
		params := netdicom.ServiceProviderParams{
			CStore: func(
				connState netdicom.ConnectionState,
				transferSyntaxUID string,
				sopClassUID string,
				sopInstanceUID string,
				data []byte) dimse.Status {
				return dimse.Status{Status: dimse.StatusSuccess}
			},
		}

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Accept error: %v", err)
				break
			}
			log.Printf("Accepted connection %v", conn)
			netdicom.RunProviderForConn(conn, params)
		}
	}()
	return listener
}

func runClient(serverAddr string, faults netdicom.FaultInjector) {
	dataset, err := dicom.ReadDataSetFromFile(
		"../testdata/reportsi.dcm",
		dicom.ReadOptions{})
	if err != nil {
		log.Fatal(err)
	}
	netdicom.SetUserFaultInjector(faults)
	su, err := netdicom.NewServiceUser(netdicom.ServiceUserParams{SOPClasses: sopclass.StorageClasses})
	if err != nil {
		log.Fatal(err)
	}
	su.Connect(serverAddr)
	err = su.CStore(dataset)
	log.Printf("Store done with status: %v", err)
	su.Release()
}

func init() {
	flag.Parse()
}

func Fuzz(data []byte) int {
	listener := startServer(netdicom.NewFuzzFaultInjector(data))
	runClient(listener.Addr().String(), netdicom.NewFuzzFaultInjector(data))
	listener.Close()
	return 0
}
