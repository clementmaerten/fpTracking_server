build :
		go get github.com/gorilla/mux
		go get github.com/gorilla/sessions
		go get github.com/satori/go.uuid
		go build

clean :
		rm fpTracking_server