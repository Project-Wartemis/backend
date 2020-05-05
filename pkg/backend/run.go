package backend

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func RunBackend(addr *string) {
	logrus.Infof("Running server: %s", *addr)
	recep := newReception()
	startRestApi(recep)
	logrus.Fatal(http.ListenAndServe(*addr, nil))
}
