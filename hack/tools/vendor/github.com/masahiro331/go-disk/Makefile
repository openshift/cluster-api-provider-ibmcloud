SECTORS		= 32768
SECTOR_SIZE	= 512
TARGET_MBR	= mbr.bin

mbr.bin:
	dd if=/dev/zero of=$(TARGET_MBR) seek=$(SECTORS) bs=$(SECTOR_SIZE) count=1
	fdisk -i -y $(TARGET_MBR)