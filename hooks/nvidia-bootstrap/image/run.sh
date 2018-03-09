# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/sh

# Simple early detection of nvidia card
grep 10de118a /proc/bus/pci/devices || exit 0

# g2.2xlarge
# 0018    10de118a        64              ec000000                e000000c                       0                ea00000c                       0                    c101                ee000000                 1000000                 8000000                       0                 2000000                       0                     80            80000        nvidia


# This is pretty annoying.... note this is installed onto the host
chroot /rootfs apt-get update
chroot /rootfs apt-get install --yes gcc

mkdir -p /rootfs/tmp
cd /rootfs/tmp
# TODO: We can't download over SSL - presents an akamai cert
wget http://us.download.nvidia.com/XFree86/Linux-x86_64/367.57/NVIDIA-Linux-x86_64-367.57.run
echo 'fc94f5df7eb2ef243db381bc4458f911a6d76bff949701bedb249a3ebf369ff3da8dc5a7d52ab6ae3f23e947c419923f303cd57429a266a0f8e96df1039b1f5d  NVIDIA-Linux-x86_64-367.57.run' | sha512sum -c - || exit 1
chmod +x NVIDIA-Linux-x86_64-367.57.run
chroot /rootfs /tmp/NVIDIA-Linux-x86_64-367.57.run -s --install-libglvnd

chroot /rootfs nvidia-xconfig --busid=PCI:0:3:0 --use-display-device=None --virtual=1280x720
