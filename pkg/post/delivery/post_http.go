package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"redditapp/pkg/comment"
	cuc "redditapp/pkg/comment/usecase"
	"redditapp/pkg/post"
	puc "redditapp/pkg/post/usecase"
	"redditapp/tools"
	"strconv"
)

type PostHandle struct {
	PostUsecase    *puc.PostUsecase
	CommentUsecase *cuc.CommentUsecase
}

type DeliveryPost struct {
	post.Post
	Comments []*comment.Comment `json:"comments"`
}

type CommentJson struct {
	Comment string `json:"comment"`
}

type GetterBy struct {
	Field   string
	GetFunc func(string) ([]post.Post, error)
}

func NewPostHandle(postUse *puc.PostUsecase, commUse *cuc.CommentUsecase) *PostHandle {
	p := &PostHandle{
		PostUsecase:    postUse,
		CommentUsecase: commUse,
	}
	return p
}

func (h *PostHandle) GetAll(w http.ResponseWriter, r *http.Request) {

	list, err := h.PostUsecase.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	deliveryList, err := h.NewDeliveryList(list)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	returnList(w, deliveryList)
}

func (h *PostHandle) GetBy(w http.ResponseWriter, r *http.Request, getter *GetterBy) {
	vars := mux.Vars(r)
	field, ok := vars[getter.Field]
	if !ok {
		http.Error(w, fmt.Sprint(`{"message":"parm not found"}`), http.StatusNotModified)
		return
	}

	list, err := getter.GetFunc(field)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	deliveryList, err := h.NewDeliveryList(list)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	returnList(w, deliveryList)
}

func (h *PostHandle) GetByCategory(w http.ResponseWriter, r *http.Request) {
	h.GetBy(w, r, &GetterBy{
		Field:   "category_name",
		GetFunc: h.PostUsecase.GetByCategory,
	})
}

func (h *PostHandle) GetByUser(w http.ResponseWriter, r *http.Request) {
	h.GetBy(w, r, &GetterBy{
		Field:   "login",
		GetFunc: h.PostUsecase.GetByUser,
	})
}

func (h *PostHandle) GetPost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pst, err := h.PostUsecase.ShowPost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deliveryPost, err := h.NewDeliveryPost(pst)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(deliveryPost)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(res)
	return
}

func (h *PostHandle) AddPost(w http.ResponseWriter, r *http.Request) {

	pst := &post.Post{}
	err := tools.RequestBody(r, pst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usr, ok := tools.CurrentUser(r)
	if !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	_, err = h.PostUsecase.PostRepo.AddPost(pst, usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(pst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
	//fmt.Fprint(w, res)
}

func (h *PostHandle) AddComment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, `{"message":"there isn't param id"}`, http.StatusNotFound)
		return
	}

	commData := &CommentJson{}
	err := tools.RequestBody(r, commData)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	pst, err := h.PostUsecase.GetPost(id)
	if err != nil {
		http.Error(w, `{"message":"post not found"}`, http.StatusNotFound)
		return
	}

	currUser, ok := tools.CurrentUser(r)
	if !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	_, err = h.CommentUsecase.AddComment(pst.Id, commData.Comment, currUser)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	deliveryPost, err := h.NewDeliveryPost(pst)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(deliveryPost)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)

}

func (h *PostHandle) SetVote(w http.ResponseWriter, r *http.Request, funcVote func(string, string) error) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, `{"message":"there isn't param id"}`, http.StatusNotFound)
		return
	}

	usr, ok := tools.CurrentUser(r)
	if !ok {
		http.Error(w, `{"message":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	err := funcVote(id, strconv.Itoa(usr.Id))
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	pst, err := h.PostUsecase.GetPost(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	deliveryPost, err := h.NewDeliveryPost(pst)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(deliveryPost)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func (h *PostHandle) UpVote(w http.ResponseWriter, r *http.Request) {
	h.SetVote(w, r, h.PostUsecase.UpVote)
}

func (h *PostHandle) DownVote(w http.ResponseWriter, r *http.Request) {
	h.SetVote(w, r, h.PostUsecase.DownVote)
}

func (h *PostHandle) UnVote(w http.ResponseWriter, r *http.Request) {
	h.SetVote(w, r, h.PostUsecase.UnVote)
}

func (h *PostHandle) Delete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, `{"message":"there isn't param id"}`, http.StatusNotFound)
		return
	}

	usr, ok := tools.CurrentUser(r)
	if !ok {
		http.Error(w, `{"message":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	err := h.PostUsecase.Delete(id, usr.Name)
	if err != nil {
		http.Error(w, `{"message":"error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (h *PostHandle) NewDeliveryPost(pst *post.Post) (*DeliveryPost, error) {

	comments, err := h.CommentUsecase.GetByPost(pst.Id)
	if err != nil {
		return nil, err
	}

	res := &DeliveryPost{
		Post:     *pst,
		Comments: comments,
	}

	return res, nil
}

func (h *PostHandle) NewDeliveryList(list []post.Post) ([]*DeliveryPost, error) {

	//TODO: n+1 problem
	var err error
	deliveryList := make([]*DeliveryPost, len(list))
	for i, v := range list {
		deliveryList[i], err = h.NewDeliveryPost(&v)
		if err != nil {
			return nil, err
		}
	}

	return deliveryList, nil
}

func returnList(w http.ResponseWriter, list []*DeliveryPost) {
	res, err := json.Marshal(list)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
