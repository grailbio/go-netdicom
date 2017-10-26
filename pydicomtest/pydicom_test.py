#!/usr/bin/env python3.6

"""Test that go-netdicom and pynetdicom3 can interoperate."""

import glob
import logging
import random
import socket
import subprocess
import tempfile
import time
import unittest

import pydicom
import pynetdicom3

logging.basicConfig(level=logging.DEBUG)

class TestPydicom(unittest.TestCase):
    def setUp(self) -> None:
        self.server_port = random.randrange(10000, 20000)
        self.tempdir = tempfile.mkdtemp()
        logging.info("Start server at port %d output %s", self.server_port, self.tempdir)
        self.server = subprocess.Popen(['../storeserver/storeserver',
                                        '--vmodule', '*=0',
                                        '--port', f'{self.server_port}',
                                        '--output', self.tempdir])
        while True:
            with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
                result = sock.connect_ex(('localhost', self.server_port))
                if result == 0:
                    break
            time.sleep(0.1)

    def tearDown(self) -> None:
        self.server.kill()
        self.server.wait()
        # shutil.rmtree(self.tempdir)
        logging.info("Kill server at %d output %s", self.server_port, self.tempdir)

    # TODO(saito) Test the combo go-netdicom client, pynetdicom client too.

    def test_store(self) -> None:
        """Test that pynetdicom client can talk to gonetdicom."""
        ae = pynetdicom3.AE(ae_title="testclient",
                            port=0,
                            scu_sop_class=pynetdicom3.StorageSOPClassList,
                            scp_sop_class=[],
                            transfer_syntax=[pydicom.uid.ExplicitVRLittleEndian])
        assoc = ae.associate('localhost', self.server_port, 'testserver')
        with open('../testdata/reportsi.dcm', 'rb') as f:
            in_ds = pydicom.read_file(f, force=True)
        self.assertEqual(assoc.send_c_store(in_ds).code, 0)
        assoc.release()

        outputs = glob.glob(self.tempdir + "/*.dcm")
        self.assertEqual(len(outputs), 1)
        with open(outputs[0], 'rb') as f:
            out_ds = pydicom.read_file(f, force=True)
        self.assertEqual(str(in_ds), str(out_ds))

if __name__ == '__main__':
    unittest.main()
