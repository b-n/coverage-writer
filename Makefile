.PHONY: gomodgen deploy delete

gomodgen:
	GO113MODULE=on go mod init

deploy:
	gcloud functions deploy coverage --entry-point Handle --runtime go111 --trigger-http

delete:
	gcloud functions delete coverage --entry-point Handle --runtime go111 --trigger-http
