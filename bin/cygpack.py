#!/usr/bin/python3

import os
import sys
from getopt import getopt
import zipfile

# z: zip u: unzip d: directory
opts, args = getopt(sys.argv[1:], "z:u:d:", ["zip=", "unzip=", "dir="])

mode = "null"
file_name = ""
directory = ""

for (opt, arg) in opts:
	if opt in ("-z", "--zip"):
		if mode != "null":
			print("invalid arguments")
			sys.exit(2)
		mode = "zip"
		file_name = arg
	elif opt in ("-u", "--unzip"):
		if mode != "null":
			print("invalid arguments")
			sys.exit(2)
		mode = "unzip"
		file_name = arg
	elif opt in ("-d", "--dir"):
		directory = arg
	else:
		print("invalid arguments")
		sys.exit(1)

if mode == "zip":
	zip_file = zipfile.ZipFile(file_name, 'w')
	compressed_files = []
	for path, dir_lst, file_lst in os.walk(directory):
		for fname in file_lst:
			zip_file.write(os.path.join(path, fname), arcname=fname, compress_type=zipfile.ZIP_DEFLATED)
	zip_file.close()
elif mode == "unzip":
	zip_file = zipfile.ZipFile(file_name)
	zip_file.extractall(directory)
	zip_file.close()
else:
	print("invalid arguments")
	sys.exit(2)
