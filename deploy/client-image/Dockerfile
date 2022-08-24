FROM localhost:5000/air:latest AS build

WORKDIR /src/

ADD ./integration/puretcp /src/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod graph | awk '$1 !~ /@/ { print $2 }' | xargs -r go get -x && \
    go list -m -f '{{ if not .Main }}{{ .Path }}/...@{{ .Version }}{{ end }}' all | tail -n +2 | \
    CGO_ENABLED=1 GOOS=linux xargs go build -v -installsuffix cgo -i; echo done

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    mkdir bin && \
    GOOS=linux CGO_ENABLED=0 go build -o ./bin/ -ldflags "-extldflags -static -s -w" ./cmd/client/

FROM scratch

# Import the user and group files from the build stage.
COPY --from=build /user/group /user/passwd /etc/
# Import the Certificate-Authority certificates for enabling HTTPS.
#COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV APP_ROOT /opt/user
ENV PATH /opt/user

COPY --from=build /src/bin $APP_ROOT
#COPY --from=build /src/build/certs/*.crt /usr/local/share/ca-certificates/
#EXPOSE 8000
USER user
CMD ["/opt/user/client"]