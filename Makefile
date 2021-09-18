default: yarn fmt publiccopy 

	cd src/backend; cargo build

	$(MAKE) publicdelete

	mv src/backend/target/debug/dietpi_dashboard ./dietpi-dashboard

rust: fmt publiccopy

	cd src/backend; cargo build
	
	$(MAKE) publicdelete

	mv src/backend/target/debug/dietpi_dashboard ./dietpi-dashboard

yarn:
	cd src/frontend; yarn build

publiccopy:
	cp -r src/frontend/public src/backend/src

publicdelete:
	rm -r src/backend/src/public

fmt:
	rustfmt --edition 2018 src/backend/src/*

rustbuild: fmt publiccopy
	mkdir -p build/

	cd src/backend; cargo build --release --target x86_64-unknown-linux-gnu
	x86_64-linux-gnu-strip src/backend/target/x86_64-unknown-linux-gnu/release/dietpi_dashboard
	upx --lzma src/backend/target/x86_64-unknown-linux-gnu/release/dietpi_dashboard
	mv src/backend/target/x86_64-unknown-linux-gnu/release/dietpi_dashboard build/dietpi-dashboard-amd64

	cd src/backend; cargo build --release --target aarch64-unknown-linux-gnu
	aarch64-linux-gnu-strip src/backend/target/aarch64-unknown-linux-gnu/release/dietpi_dashboard
	upx --lzma src/backend/target/aarch64-unknown-linux-gnu/release/dietpi_dashboard
	mv src/backend/target/aarch64-unknown-linux-gnu/release/dietpi_dashboard build/dietpi-dashboard-arm64

	cd src/backend; cargo build --release --target arm-unknown-linux-gnueabihf
	/opt/rpi/arm-bcm2708/arm-linux-gnueabihf/bin/arm-linux-gnueabihf-strip src/backend/target/arm-unknown-linux-gnueabihf/release/dietpi_dashboard
	upx --lzma src/backend/target/arm-unknown-linux-gnueabihf/release/dietpi_dashboard
	mv src/backend/target/arm-unknown-linux-gnueabihf/release/dietpi_dashboard build/dietpi-dashboard-armhf

	$(MAKE) publicdelete

build: yarn rustbuild