run:
    docker build -t kyte --progress none . && \
    docker run \
        --rm \
        -u $(id -u):$(id -g) \
        -p 8000:8000 \
        -v $(pwd)/kyte_cache:/opt/kyte/cache/packages \
        -w /opt/kyte/cache/packages \
        --name kyte_backend \
        kyte