From golang
COPY . $GOPATH/src/github.com/georgethomas111/urlshortner
RUN go install github.com/georgethomas111/urlshortner/main
EXPOSE 8080
CMD ["bin/main"] 
