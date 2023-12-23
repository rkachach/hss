package hss

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rkachach/hss/internal/objectStore"
)

func PutBucket(w http.ResponseWriter, r *http.Request) {

	directoryName := mux.Vars(r)["directory"]

	err := objectStore.CreateDirectory(directoryName)
	if err != nil {
		// writeErrorResponse(w, errorCodes.ToAPIErr(ErrBucketAlreadyExists), r.URL)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}
