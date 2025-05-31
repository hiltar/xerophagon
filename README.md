# xerophagon
Xerophagon is a fasting web application made with `golang 1.24`

# Installation guide

## SD card operations
```
sudo mount /dev/sdc1 /mnt/sdcard

sudo mkdir /mnt/sdcard/xerophagon
sudo mkdir /mnt/sdcard/xerophagon/static
sudo mkdir /mnt/sdcard/xerophagon/templates
# Copy xerophagon-arm executable into /mnt/sdcard/xerophagon
# Copy static files into /mnt/sdcard/xerophagon/static
# Copy template files into /mnt/sdcard/xerophagon/templates

sync
sudo umount mnt/sdcard
```

## Setup xerophagon
```
# SSH into rpi
ssh root@<rpi-ip>

# Copy xerophagon file contents from init.d folder into /etc/init.d/xerophagon

# Update RC configuration
chmod +x /etc/init.d/xerophagon
rc-update add xerophagon default

# data folder
mkdir /opt/xerophagon/

lbu add /opt/xerophagon/
lbu commit
reboot # after reboot xerophagon should be available from <rpi-ip>:5555
rc-service xerophagon status

```
