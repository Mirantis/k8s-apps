build:
        # Create vendor dir
	glide i
	# Hack for creation right vendor
	mv vendor/github.com/kubernetes-incubator/ ./
	rm -rf vendor/
	mv kubernetes-incubator/service-catalog/vendor/ ./
	mv kubernetes-incubator/ vendor/github.com/
	# Build
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo
