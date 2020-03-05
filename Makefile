PROJ := json-fmt

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

deploy: build
	docker build --tag gcr.io/darren-prd/$(PROJ) . 
	docker push gcr.io/darren-prd/$(PROJ)
	@rm -rf server
	gcloud run deploy $(PROJ) --image gcr.io/darren-prd/$(PROJ) --platform managed --max-instances 1
	
# to use in a test curl like so
# curl -v -H "Authorization: Bearer $(gcloud auth print-identity-token)" https://blah-yiohemudcq-uc.a.run.app/start
auth:
	@gcloud auth print-identity-token