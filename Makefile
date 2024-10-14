run:
	CGO_ENABLED=0 go run cmd/main.go

build:
	CGO_ENABLED=0 go build cmd/main.go && mv main app.exe

tidy:
	go mod tidy

clean_cals:
	python3 clean_cal.py

clean:
	rm *.exe