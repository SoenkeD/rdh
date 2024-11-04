#############################
#############################
### Exec 
#############################
#############################

.PHONY: test
test: vet fmt
	~/go/bin/ginkgo -r

#############################
#############################
### Golang 
#############################
#############################
.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...