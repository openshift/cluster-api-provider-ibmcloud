SECTORS		= 32768
SECTOR_SIZE	= 512
TARGET_MBR	= mbr.bin
TARGET_IMG	= linux.img
TARGET_FS	= fs.bin

mbr.bin:
	dd if=/dev/zero of=$(TARGET_MBR) seek=$(SECTORS) bs=$(SECTOR_SIZE) count=1
	fdisk -i -y $(TARGET_MBR)

# Linux only and require admin authority
#!!! Extract xfs from gpt using cmd/main.go.
linux.img: main
	$(eval DEVICE := $(shell losetup -f))
	dd of=$(TARGET_IMG) count=0 seek=1 bs=41943040
	losetup $(DEVICE) $(TARGET_IMG)
	parted $(DEVICE) -s mklabel gpt -s mkpart primary xfs 0 100%
	mkfs.xfs $(DEVICE)p1
	losetup -d $(DEVICE)

main: cmd/main.go
	go build -o main cmd/main.go

fs.bin: linux.img main
	./main linux.img
	mv primary0 $(TARGET_FS)