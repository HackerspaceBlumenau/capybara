FROM golang as builder
ADD . /go/src/github.com/hackerspaceblumenau/capybara
RUN cd /go/src/github.com/hackerspaceblumenau/capybara \
    && go build -o /capybara

FROM debian:9-slim
COPY --from=builder /capybara /
ENTRYPOINT ["/capybara"]
