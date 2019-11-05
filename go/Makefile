linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ishouru_job main.go 
windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ishouru_job.exe main.go 
mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ishouru_job main.go 

gqlgen:
	cd app/graphql && gqlgen