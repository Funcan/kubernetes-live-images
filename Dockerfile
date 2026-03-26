FROM gcr.io/distroless/static-debian12:nonroot
COPY kubernetes-live-images /kubernetes-live-images
ENTRYPOINT ["/kubernetes-live-images"]
