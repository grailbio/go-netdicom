package main

// A simple PACS server. Supports C-STORE, C-FIND, C-MOVE.
//
// Usage: ./sampleserver -dir <directory> -port 11111
//
// It starts a DICOM server and serves files under <directory>.

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/grailbio/go-dicom"
	"github.com/grailbio/go-dicom/dicomio"
	"github.com/grailbio/go-dicom/dicomtag"
	"github.com/grailbio/go-dicom/dicomuid"
	"github.com/grailbio/go-netdicom"
	"github.com/grailbio/go-netdicom/dimse"
)

var (
	portFlag     = flag.String("port", "10000", "TCP port to listen to")
	aeFlag       = flag.String("ae", "bogusae", "AE title of this server")
	remoteAEFlag = flag.String("remote-ae", "GBMAC0261:localhost:11112", `
Comma-separated list of remote AEs, in form aetitle:host:port, For example -remote-ae testae:foo.example.com:12345,testae2:bar.example.com:23456.
In this example, a C-GET or C-MOVE request to application entity "testae" will resolve to foo.example.com:12345.`)
	dirFlag = flag.String("dir", ".", `
The directory to locate DICOM files to report in C-FIND, C-MOVE, etc.
Files are searched recursivsely under this directory.
Defaults to '.'.`)
	outputFlag = flag.String("output", "", `
The directory to store files received by C-STORE.
If empty, use <dir>/incoming, where <dir> is the value of the -dir flag.`)

	tlsKeyFlag  = flag.String("tls-key", "", "Sets the private key file. If empty, TLS is disabled.")
	tlsCertFlag = flag.String("tls-cert", "", "File containing TLS cert to be presented to the peer.")
	tlsCAFlag   = flag.String("tls-ca", "", "Optional file containing certs to match against what peers present.")
)

type server struct {
	mu *sync.Mutex

	// Set of dicom files the server manages. Keys are file paths.  Guarded
	// by mu.
	datasets map[string]*dicom.DataSet

	// For generating new unique path in C-STORE. Guarded by mu.
	pathSeq int32
}

func (ss *server) onCStore(
	transferSyntaxUID string,
	sopClassUID string,
	sopInstanceUID string,
	data []byte) dimse.Status {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.pathSeq++
	path := path.Join(*outputFlag, fmt.Sprintf("image%04d.dcm", ss.pathSeq))
	out, err := os.Create(path)
	if err != nil {
		dirPath := filepath.Dir(path)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return dimse.Status{Status: dimse.StatusNotAuthorized, ErrorComment: err.Error()}
		}
		out, err = os.Create(path)
		if err != nil {
			log.Printf("%s: create: %v", path, err)
			return dimse.Status{Status: dimse.StatusNotAuthorized, ErrorComment: err.Error()}
		}
	}
	defer func() {
		if out != nil {
			out.Close()
		}
	}()
	e := dicomio.NewEncoderWithTransferSyntax(out, transferSyntaxUID)
	dicom.WriteFileHeader(e,
		[]*dicom.Element{
			dicom.MustNewElement(dicomtag.TransferSyntaxUID, transferSyntaxUID),
			dicom.MustNewElement(dicomtag.MediaStorageSOPClassUID, sopClassUID),
			dicom.MustNewElement(dicomtag.MediaStorageSOPInstanceUID, sopInstanceUID),
		})
	e.WriteBytes(data)
	if err := e.Error(); err != nil {
		log.Printf("%s: write: %v", path, err)
		return dimse.Status{Status: dimse.StatusNotAuthorized, ErrorComment: err.Error()}
	}
	err = out.Close()
	out = nil
	if err != nil {
		log.Printf("%s: close %s", path, err)
		return dimse.Status{Status: dimse.StatusNotAuthorized, ErrorComment: err.Error()}
	}
	log.Printf("C-STORE: Created %v", path)
	// Register the new file in ss.datasets.
	ds, err := dicom.ReadDataSetFromFile(path, dicom.ReadOptions{DropPixelData: true})
	if err != nil {
		log.Printf("%s: failed to parse dicom file: %v", path, err)
	} else {
		ss.datasets[path] = ds
	}
	return dimse.Success
}

// Represents a match.
type filterMatch struct {
	path  string           // DICOM path name
	elems []*dicom.Element // Elements within "ds" that match the filter
}

// "filters" are matching conditions specified in C-{FIND,GET,MOVE}. This
// function returns the list of datasets and their elements that match filters.
func (ss *server) findMatchingFiles(filters []*dicom.Element) ([]filterMatch, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	var matches []filterMatch
	for path, ds := range ss.datasets {
		allMatched := true
		match := filterMatch{path: path}
		for _, filter := range filters {
			ok, elem, err := dicom.Query(ds, filter)
			if err != nil {
				return matches, err
			}
			if !ok {
				log.Printf("DS: %s: filter %v missed", path, filter)
				allMatched = false
				break
			}
			if elem != nil {
				match.elems = append(match.elems, elem)
			} else {
				elem, err := dicom.NewElement(filter.Tag)
				if err != nil {
					log.Println(err)
					return matches, err
				}
				match.elems = append(match.elems, elem)
			}
		}
		if allMatched {
			if len(match.elems) == 0 {
				panic(match)
			}
			matches = append(matches, match)
		}
	}
	return matches, nil
}

func (ss *server) onCFind(
	transferSyntaxUID string,
	sopClassUID string,
	filters []*dicom.Element,
	ch chan netdicom.CFindResult) {
	for _, filter := range filters {
		log.Printf("CFind: filter %v", filter)
	}
	log.Printf("CFind: transfersyntax: %v, classuid: %v",
		dicomuid.UIDString(transferSyntaxUID),
		dicomuid.UIDString(sopClassUID))
	// Match the filter against every file. This is just for demonstration
	matches, err := ss.findMatchingFiles(filters)
	log.Printf("C-FIND: found %d matches, err %v", len(matches), err)
	if err != nil {
		ch <- netdicom.CFindResult{Err: err}
	} else {
		for _, match := range matches {
			log.Printf("C-FIND resp %s: %v", match.path, match.elems)
			ch <- netdicom.CFindResult{Elements: match.elems}
		}
	}
	close(ch)
}

func (ss *server) onCMoveOrCGet(
	transferSyntaxUID string,
	sopClassUID string,
	filters []*dicom.Element,
	ch chan netdicom.CMoveResult) {
	log.Printf("C-MOVE: transfersyntax: %v, classuid: %v",
		dicomuid.UIDString(transferSyntaxUID),
		dicomuid.UIDString(sopClassUID))
	for _, filter := range filters {
		log.Printf("C-MOVE: filter %v", filter)
	}

	matches, err := ss.findMatchingFiles(filters)
	log.Printf("C-MOVE: found %d matches, err %v", len(matches), err)
	if err != nil {
		ch <- netdicom.CMoveResult{Err: err}
	} else {
		for i, match := range matches {
			log.Printf("C-MOVE resp %d %s: %v", i, match.path, match.elems)
			// Read the file; the one in ss.datasets lack the PixelData.
			ds, err := dicom.ReadDataSetFromFile(match.path, dicom.ReadOptions{})
			resp := netdicom.CMoveResult{
				Remaining: len(matches) - i - 1,
				Path:      match.path,
			}
			if err != nil {
				resp.Err = err
			} else {
				resp.DataSet = ds
			}
			ch <- resp
		}
	}
	close(ch)
}

// Find DICOM files in or under "dir" and read its attributes. The return value
// is a map from a pathname to dicom.Dataset (excluding PixelData).
func listDicomFiles(dir string) (map[string]*dicom.DataSet, error) {
	datasets := make(map[string]*dicom.DataSet)
	readFile := func(path string) {
		if _, ok := datasets[path]; ok {
			return
		}
		ds, err := dicom.ReadDataSetFromFile(path, dicom.ReadOptions{DropPixelData: true})
		if err != nil {
			log.Printf("%s: failed to parse dicom file: %v", path, err)
			return
		}
		log.Printf("%s: read dicom file", path)
		datasets[path] = ds
	}
	walkCallback := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("%v: skip file: %v", path, err)
			return nil
		}
		if (info.Mode() & os.ModeDir) != 0 {
			// If a directory contains file "DICOMDIR", all the files in the directory are DICOM files.
			if _, err := os.Stat(filepath.Join(path, "DICOMDIR")); err != nil {
				return nil
			}
			subpaths, err := filepath.Glob(path + "/*")
			if err != nil {
				log.Printf("%v: glob: %v", path, err)
				return nil
			}
			for _, subpath := range subpaths {
				if !strings.HasSuffix(subpath, "DICOMDIR") {
					readFile(subpath)
				}
			}
			return nil
		}
		if strings.HasSuffix(path, ".dcm") {
			readFile(path)
		}
		return nil
	}
	if err := filepath.Walk(dir, walkCallback); err != nil {
		return nil, err
	}
	return datasets, nil
}

func parseRemoteAEFlag(flag string) (map[string]string, error) {
	aeMap := make(map[string]string)
	re := regexp.MustCompile("^([^:]+):(.+)$")
	for _, str := range strings.Split(flag, ",") {
		if str == "" {
			continue
		}
		m := re.FindStringSubmatch(str)
		if m == nil {
			return aeMap, fmt.Errorf("Failed to parse AE spec '%v'", str)
		}
		log.Printf("Remote AE '%v' -> '%v'", m[1], m[2])
		aeMap[m[1]] = m[2]
	}
	return aeMap, nil
}

func canonicalizeHostPort(addr string) string {
	if !strings.Contains(addr, ":") {
		return ":" + addr
	}
	return addr
}

func main() {
	flag.Parse()
	port := canonicalizeHostPort(*portFlag)
	if *outputFlag == "" {
		*outputFlag = filepath.Join(*dirFlag, "incoming")
	}
	remoteAEs, err := parseRemoteAEFlag(*remoteAEFlag)
	if err != nil {
		log.Panicf("Failed to parse -remote-ae flag: %v", err)
	}
	datasets, err := listDicomFiles(*dirFlag)
	if err != nil {
		log.Panicf("Failed to list DICOM files in %s: %v", *dirFlag, err)
	}
	ss := server{
		mu:       &sync.Mutex{},
		datasets: datasets,
	}
	log.Printf("Listening on %s", port)

	var tlsConfig *tls.Config
	if *tlsKeyFlag != "" {
		cert, err := tls.LoadX509KeyPair(*tlsCertFlag, *tlsKeyFlag)
		if err != nil {
			log.Panic(err)
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		if *tlsCAFlag != "" {
			ca, err := ioutil.ReadFile(*tlsCAFlag)
			if err != nil {
				log.Panic(err)
			}
			tlsConfig.RootCAs = x509.NewCertPool()
			tlsConfig.RootCAs.AppendCertsFromPEM(ca)
			tlsConfig.BuildNameToCertificate()
		}
	}

	params := netdicom.ServiceProviderParams{
		AETitle:   *aeFlag,
		RemoteAEs: remoteAEs,
		CEcho: func(connState netdicom.ConnectionState) dimse.Status {
			log.Printf("Received C-ECHO")
			return dimse.Success
		},
		CFind: func(connState netdicom.ConnectionState, transferSyntaxUID string, sopClassUID string,
			filter []*dicom.Element, ch chan netdicom.CFindResult) {
			ss.onCFind(transferSyntaxUID, sopClassUID, filter, ch)
		},
		CMove: func(connState netdicom.ConnectionState, transferSyntaxUID string, sopClassUID string,
			filter []*dicom.Element, ch chan netdicom.CMoveResult) {
			ss.onCMoveOrCGet(transferSyntaxUID, sopClassUID, filter, ch)
		},
		CGet: func(connState netdicom.ConnectionState, transferSyntaxUID string, sopClassUID string,
			filter []*dicom.Element, ch chan netdicom.CMoveResult) {
			ss.onCMoveOrCGet(transferSyntaxUID, sopClassUID, filter, ch)
		},
		CStore: func(connState netdicom.ConnectionState, transferSyntaxUID string,
			sopClassUID string,
			sopInstanceUID string,
			data []byte) dimse.Status {
			return ss.onCStore(transferSyntaxUID, sopClassUID, sopInstanceUID, data)
		},
		TLSConfig: tlsConfig,
	}
	sp, err := netdicom.NewServiceProvider(params, port)
	if err != nil {
		panic(err)
	}
	sp.Run()
}
