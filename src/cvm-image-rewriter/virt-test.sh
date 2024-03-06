#!/bin/bash
#
# Create TD VM from libvirt template
#
set -e
set -m 
set -x

CURR_DIR=$(readlink -f "$(dirname "$0")")

GUEST_IMG="tdx-guest-ubuntu22.04.qcow2"
GUEST_NAME="tdx-guest"
GUEST_KERNEL="/boot/vmlinuz"
OVMF_CODE="/usr/share/qemu/OVMF.fd"

GUEST_ROOTDIR=/tmp/libvirt-vms
TEMPLATE="${CURR_DIR}/tdx-libvirt-ubuntu-host.xml.template"
FORCE=false
VCPU_NUM=1
MEM_SIZE=4

# TDX socket
TDX_SOCKET="/var/local/run/libvirt/libvirt-sock"

# log file 
LOG="/tmp/vm_log_$(date +"%FT%H%M").log"

usage() {
    cat << EOM
Usage: $(basename "$0") [OPTION]...
  -i <guest image file>     Default is tdx-guest-ubuntu22.04.qcow2 under current directory
  -n <guest name>           Name of TD guest
  -k <guest kernel>         Name of kernel
  -t <template file>        Default is "${CURR_DIR}/tdx-libvirt-ubuntu-host.xml.template"
  -f                        Force recreate
  -v <vcpu number>          VM vCPU number
  -m <memory size in GB>    VM memory size in GB
  -h                        Show this help
EOM
}

pre-check() {
    # Check whether TD guest name used
    if [ -e "/tmp/libvirt-vms/$GUEST_NAME.xml" ]; then 
	    rm -rf /tmp/libvirt-vms/$GUEST_NAME.xml
	    echo "rm -rf /tmp/libvirt-vms/$GUEST_NAME.xml"
    fi

    # Check if the socket file exists
    if [ ! -e "$TDX_SOCKET" ]; then
	    echo "Socket file does not exists: $TDX_SOCKET"
    else
	    echo "Socket file exists: $TDX_SOCKET"
    fi

    # Check whether current user belong to libvirt
    if [[ ! $(id -nG "$USER") == *"libvirt"* ]]; then
    echo WARNING! Please add user "$USER" into group "libvirt" via \"sudo usermod -aG libvirt "$USER"\"
    return 1
    fi
}

process_args() {
    while getopts ":i:n:k:t:v:m:fh" option; do
        case "$option" in
            i) GUEST_IMG=$OPTARG;;
            n) GUEST_NAME=$OPTARG;;
            k) GUEST_KERNEL=$OPTARG;;
            t) TEMPLATE=$OPTARG;;
            v) VCPU_NUM=$OPTARG;;
            m) MEM_SIZE=$OPTARG;;
            f) FORCE=true;;
            h) usage
               exit 0
               ;;
            *)
               echo "Invalid option '-$OPTARG'"
               usage
               exit 1
               ;;
        esac
    done

    echo "====================================================================="
    echo " Use Template   : ${TEMPLATE}"
    echo " Guest Name     : ${GUEST_NAME}"
    echo " Guest XML      : ${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
    echo " Guest Image    : ${GUEST_ROOTDIR}/${GUEST_NAME}.qcow2"
    echo " Guest Kernel   : ${GUEST_KERNEL}"
    echo " OVMF    	  : ${OVMF_CODE}"
    echo " Force Recreate : ${FORCE}"
    echo "====================================================================="

    if [[ ! -f ${GUEST_IMG} ]]; then
        echo "Error: Guest image ${GUEST_IMG} does not exist"
        exit 1
    fi

    if [[ ${FORCE} == "true" ]]; then
        echo "> Clean up the old guest... "
        virsh destroy "${GUEST_NAME}" || true
        sleep 2
        virsh undefine "${GUEST_NAME}" || true
        sleep 2
        rm "${GUEST_ROOTDIR}/${GUEST_NAME}.xml" -fr || true
    fi

    if [[ -f ${GUEST_ROOTDIR}/${GUEST_NAME}.xml ]]; then
        echo "Error: Guest XML ${GUEST_ROOTDIR}/${GUEST_NAME}.xml already exist."
        echo "Error: you can delete the old one via 'rm ${GUEST_ROOTDIR}/${GUEST_NAME}.xml'"
        exit 1
    fi

    if [[ ! -f ${TEMPLATE} ]]; then
        echo "Template ${TEMPLATE} does not exist".
        echo "Please specify via -t"
        exit 1
    fi

    # Validate the number of vCPUs
    if ! [[ ${VCPU_NUM} =~ ^[0-9]+$ && ${VCPU_NUM} -gt 0 ]]; then
        echo "Error: Invalid number of vCPUs: ${VCPU_NUM}"
        usage
        exit 1
    fi

    # Validate the size of memory
    if ! [[ ${MEM_SIZE} =~ ^[0-9]+$ && ${MEM_SIZE} -gt 0 ]]; then
        echo "Error: Invalid memory size: ${MEM_SIZE}"
        usage
        exit 1
    fi
}

create-vm() {
    mkdir -p ${GUEST_ROOTDIR}/
    echo "> Create ${GUEST_ROOTDIR}/${GUEST_NAME}.qcow2..."
    cp "${GUEST_IMG}" "${GUEST_ROOTDIR}/${GUEST_NAME}.qcow2"
    echo "> Create ${GUEST_ROOTDIR}/${GUEST_NAME}.xml..."
    cp "${TEMPLATE}" "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"

    echo "> Modify configurations..."
    sed -i "s/.*<name>.*/<name>${GUEST_NAME}<\/name>/" "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
    sed -i "s#REPLACE_IMAGE#${GUEST_ROOTDIR}/${GUEST_NAME}.qcow2#" "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
    sed -i "s/REPLACE_VCPU_NUM/${VCPU_NUM}/g" "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
    sed -i "s/REPLACE_MEM_SIZE/${MEM_SIZE}/g" "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
    sed -i "s#REPLACE_LOG#${LOG}#" "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
    sed -i "s#REPLACE_OVMF_CODE#${OVMF_CODE}#" "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
    sed -i "s#REPLACE_KERNEL#${GUEST_KERNEL}#" "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
}

start-vm() {
    echo "> Create VM domain..."
    virsh define "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
    sleep 2
    echo "> Start VM..."
    sudo virsh start "${GUEST_NAME}"
    sleep 2
    echo "> Connect console..."
    sudo virsh console "${GUEST_NAME}"
}

pre-check
process_args "$@"
create-vm
start-vm
