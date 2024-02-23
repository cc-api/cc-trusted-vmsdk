#!/bin/bash

CURR_DIR=$(pwd)

# setup PYTHONPATH
python3 -m pip install pip -U
python3 -m pip install "cctrusted_base @ git+https://github.com/cc-api/cc-trusted-api.git#subdirectory=common/python"
export PYTHONPATH=$PYTHONPATH:$CURR_DIR/src/python
