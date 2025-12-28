#!/bin/bash
#
# create-test-vm.sh
#
# Creates a lightweight test VM inside the Lima libvirt environment
# for testing the virsh-sandbox API control plane.
#
# Usage: ./create-test-vm.sh [vm-name]

set -euo pipefail

VM_NAME="${1:-test-vm}"
VM_MEMORY=1024
VM_VCPUS=1
VM_DISK_SIZE="5G"
BASE_IMAGE_DIR="/var/lib/libvirt/images/base"
CLOUD_IMAGE_URL="https://cloud-images.ubuntu.com/minimal/releases/jammy/release/ubuntu-22.04-minimal-cloudimg-amd64.img"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}[INFO]${NC} Creating test VM: ${VM_NAME}"

# Download cloud image if not present
CLOUD_IMAGE="${BASE_IMAGE_DIR}/ubuntu-22.04-minimal-cloudimg-amd64.img"
if [ ! -f "${CLOUD_IMAGE}" ]; then
    echo -e "${BLUE}[INFO]${NC} Downloading Ubuntu cloud image..."
    sudo mkdir -p "${BASE_IMAGE_DIR}"
    sudo wget -q --show-progress -O "${CLOUD_IMAGE}" "${CLOUD_IMAGE_URL}"
    sudo chmod 644 "${CLOUD_IMAGE}"
fi

# Create a copy for this VM
VM_DISK="${BASE_IMAGE_DIR}/${VM_NAME}.qcow2"
if [ -f "${VM_DISK}" ]; then
    echo -e "${BLUE}[INFO]${NC} VM disk already exists, skipping creation"
else
    echo -e "${BLUE}[INFO]${NC} Creating VM disk from cloud image..."
    sudo qemu-img create -f qcow2 -b "${CLOUD_IMAGE}" -F qcow2 "${VM_DISK}" "${VM_DISK_SIZE}"
    sudo chmod 644 "${VM_DISK}"
fi

# Create cloud-init configuration
CLOUD_INIT_DIR="/tmp/cloud-init-${VM_NAME}"
mkdir -p "${CLOUD_INIT_DIR}"

# User data - configure the VM
cat > "${CLOUD_INIT_DIR}/user-data" << 'USERDATA'
#cloud-config
hostname: test-vm
manage_etc_hosts: true

users:
  - name: testuser
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    lock_passwd: false

chpasswd:
  list: |
    testuser:testpassword
    root:rootpassword
  expire: false

ssh_pwauth: true

packages:
  - curl
  - wget
  - vim
  - htop

runcmd:
  - echo "Test VM is ready for virsh-sandbox testing" > /etc/motd
  - systemctl enable ssh
  - systemctl start ssh

final_message: "Test VM boot completed in $UPTIME seconds"
USERDATA

# Meta data
cat > "${CLOUD_INIT_DIR}/meta-data" << METADATA
instance-id: ${VM_NAME}
local-hostname: ${VM_NAME}
METADATA

# Create cloud-init ISO
CLOUD_INIT_ISO="${BASE_IMAGE_DIR}/${VM_NAME}-cloud-init.iso"
echo -e "${BLUE}[INFO]${NC} Creating cloud-init ISO..."
sudo genisoimage -output "${CLOUD_INIT_ISO}" -volid cidata -joliet -rock \
    "${CLOUD_INIT_DIR}/user-data" \
    "${CLOUD_INIT_DIR}/meta-data" 2>/dev/null
sudo chmod 644 "${CLOUD_INIT_ISO}"

# Clean up temp directory
rm -rf "${CLOUD_INIT_DIR}"

# Check if VM already exists
if virsh dominfo "${VM_NAME}" &>/dev/null; then
    echo -e "${BLUE}[INFO]${NC} VM ${VM_NAME} already exists"
    virsh dominfo "${VM_NAME}"
    exit 0
fi

# Create the VM using virt-install
echo -e "${BLUE}[INFO]${NC} Creating VM with virt-install..."
sudo virt-install \
    --name "${VM_NAME}" \
    --memory "${VM_MEMORY}" \
    --vcpus "${VM_VCPUS}" \
    --disk "path=${VM_DISK},format=qcow2" \
    --disk "path=${CLOUD_INIT_ISO},device=cdrom" \
    --os-variant ubuntu22.04 \
    --network network=default \
    --graphics none \
    --console pty,target_type=serial \
    --import \
    --noautoconsole \
    --wait 0

echo -e "${GREEN}[SUCCESS]${NC} Test VM '${VM_NAME}' created successfully!"
echo ""
echo "VM Details:"
virsh dominfo "${VM_NAME}"
echo ""
echo "To connect to the VM console:"
echo "  virsh console ${VM_NAME}"
echo ""
echo "To get the VM's IP address (after it boots):"
echo "  virsh domifaddr ${VM_NAME}"
echo ""
echo "Default credentials:"
echo "  Username: testuser"
echo "  Password: testpassword"

