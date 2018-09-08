#!/usr/bin/env python

import fnmatch
import os
import subprocess

for root, dirnames, filenames in os.walk('.'):
    for filename in fnmatch.filter(filenames, '*.yaml'):
        deployment = os.path.join(root, filename)
        subprocess.call("kubectl delete -f " + deployment, shell=True)
