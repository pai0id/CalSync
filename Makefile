run:
	@go run  cmd/main.go

build:
	@go build -o app.exe cmd/main.go

tidy:
	@go mod tidy

clean_cals:
	python3 clean_cal.py