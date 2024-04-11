#!/bin/bash

CURR_DIR=$(pwd)

# setup virtualenv and PYTHONPATH
apt-get update && apt-get install -y python3-virtualenv

if [[ ! -d ${CURR_DIR}/venv ]]; then
    python3 -m virtualenv -p python3 ${CURR_DIR}/venv
    source ${CURR_DIR}/venv/bin/activate
    python3 -m pip install "cctrusted_base @ git+https://github.com/cc-api/cc-trusted-api.git#subdirectory=common/python"
    python3 -m pip install -r $CURR_DIR/src/python/requirements.txt
    if [ ! $? -eq 0 ]; then
        echo "Failed to install python PIP packages, please check your proxy (https_proxy) or setup PyPi mirror."
        deactivate
        rm ${CURR_DIR}/venv -fr
        return 1
    fi
else
    source ${CURR_DIR}/venv/bin/activate
fi

export PYTHONPATH=$PYTHONPATH:$CURR_DIR/src/python
