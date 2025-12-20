run:
	docker build -t kite --progress none . && \
	docker run \
	  --rm \
		-u $(id -u):$(id -g) \
		-p 8000:8000 \
		-v $(pwd)/kite_cache:/opt/kite/cache/packages \
		-w /opt/kite/cache/packages \
		--name kite_backend \
		kite