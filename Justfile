run:
	docker build \
		-t kite-dev \
		--progress none \
		--target dev . && \
	docker run \
		-it \
		--rm \
		-p 8000:8000 \
		-v $(pwd)/kite_workdir:/var/lib/kite \
		--name kite-dev \
		kite-dev
