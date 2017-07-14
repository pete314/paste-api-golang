## Author: Peter Nagy <https://peternagy.ie>
## Since: 06, 2017
## Description: --

DEPEND=github.com/valyala/fasthttp github.com/buaazp/fasthttprouter \
        github.com/pquerna/ffjson/ffjson gopkg.in/redis.v5 \
		golang.org/x/crypto/bcrypt github.com/asaskevich/govalidator \
		github.com/gocql/gocql github.com/satori/go.uuid

paste-api-golang: clean
	go get -u $(DEPEND)
	go build -o paste-api-golang ./src/runner.go

test: paste-api-golang
	go test -v ./src/tests/modules/v0/unit/...

clean:
	rm -f paste-api-golang
