#!/usr/bin/env python

import fnmatch
import os
import subprocess
import sys

for root, dirnames, filenames in os.walk('.'):
    for filename in fnmatch.filter(filenames, '*.yaml'):
        status = filename.split(".")[1]
        deployment = os.path.join(root, filename)
        print(status)
        return_code = subprocess.call("kubectl apply -f " + deployment, shell=True)
        print (return_code)
        if status == "pass" and return_code == 1 :
            sys.exit("Expecting " + deployment + " to pass")
        if status != "pass" and return_code == 0 :
            sys.exit("Expecting " + deployment + " to fail")

