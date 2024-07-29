#!/bin/bash

CURR_DIR=$(pwd)

# Check if the script is being run as root
if [ "$EUID" -ne 0 ]; then
  echo "Please run the script as root"
  exit 1
fi

# setup virtualenv and PYTHONPATH
apt-get update
apt-get install -y python3-virtualenv pkg-config libtss-dev libtss2-dev

if [[ -d ${CURR_DIR}/venv ]]; then
    echo "===========> Remove ${CURR_DIR}/venv and create a new one"
    rm -rf {CURR_DIR}/venv
fi

python3 -m virtualenv -p python3 ${CURR_DIR}/venv
source ${CURR_DIR}/venv/bin/activate
python3 -m pip install "evidence_api @ git+https://github.com/cc-api/evidence-api.git#subdirectory=common/python"
python3 -m pip install -r $CURR_DIR/src/python/requirements.txt
if [ ! $? -eq 0 ]; then
    echo "Failed to install python PIP packages, please check your proxy (https_proxy) or setup PyPi mirror."
    deactivate
    rm ${CURR_DIR}/venv -fr
    return 1
fi

export PYTHONPATH=$PYTHONPATH:$CURR_DIR/src/python
