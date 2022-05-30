package delivery

import (
	"github.com/gorilla/mux"
	"net/http"
	"redditapp/pkg/comment/usecase"
	cuc "redditapp/pkg/comment/usecase"
	"redditapp/tools"
)

type CommentHandle struct {
	CommentUsecase *cuc.CommentUsecase
}

type CommentJson struct {
	Comment string `json:"comment"`
}

func NewCommentHandle(uc *usecase.CommentUsecase) *CommentHandle {
	return &CommentHandle{
		CommentUsecase: uc,
	}
}

func (h *CommentHandle) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, `{"message":"there isn't param id"}`, http.StatusNotFound)
		return
	}

	commentId, ok := vars["comment_id"]
	if !ok {
		http.Error(w, `{"message":"there isn't param comment_id"}`, http.StatusNotFound)
		return
	}

	usr, ok := tools.CurrentUser(r)
	if !ok {
		http.Error(w, `{"message":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	err := h.CommentUsecase.Delete(id, commentId, usr.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
