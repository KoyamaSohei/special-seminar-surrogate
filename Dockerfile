FROM golang
WORKDIR /go/src/github.com/KoyamaSohei/special-seminar-surrogate
ENV GO111MODULE=on
COPY . .
RUN go build
CMD ./special-seminar-surrogate
