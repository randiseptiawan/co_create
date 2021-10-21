package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rysmaadit/go-template/common/responder"
	"github.com/rysmaadit/go-template/external/mysql"
)

type Invite struct {
	Email string `json:"email"`
}

type InvitedUser struct {
	NamaLengkap string `json:"NamaLengkap"`
	Email       string `json:"string"`
}

type CollaboratedUser struct {
	NamaLengkap string `json:"NamaLengkap"`
	Email       string `json:"string"`
}

func CreateProjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var project mysql.Project
		payloads, err := ioutil.ReadAll(r.Body)
		if err != nil {
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		json.Unmarshal(payloads, &project)
		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		tknStr := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
				return
			}
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		if !tkn.Valid {
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}
		project.ProjectAdminId = claims.UserId
		err = mysql.CreateProject(&project)

		if err != nil {
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		projectAdmin, err := mysql.GetUserById(claims.UserId)
		if err != nil {
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		project.ProjectAdminName = projectAdmin.NamaLengkap
		project.ProjectAdminEmail = projectAdmin.Email
		responder.NewHttpResponse(r, w, http.StatusCreated, project, nil)
	}
}

func ReadProjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args := mux.Vars(r)
		i, _ := strconv.ParseUint(args["id"], 10, 64)

		project, err := mysql.ReadProject(i)
		if err != nil {
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		projectAdmin, err := mysql.GetUserById(project.ProjectAdminId)
		if err != nil {
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		project.ProjectAdminName = projectAdmin.NamaLengkap
		project.ProjectAdminEmail = projectAdmin.Email

		// collaboratorUser, err := mysql.GetCollaboratorUser(i)
		// if err != nil {
		// 	responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
		// 	return
		// }
		// fmt.Println(collaboratorUser)
		responder.NewHttpResponse(r, w, http.StatusOK, project, nil)
	}
}

func ReadAllProjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project, err := mysql.ReadAllProject()
		if err != nil {
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		for i := 0; i <= len(project)-1; i++ {
			projectAdmin, err := mysql.GetUserById(project[i].ProjectAdminId)
			if err != nil {
				responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
				return
			}
			project[i].ProjectAdminName = projectAdmin.NamaLengkap
			project[i].ProjectAdminEmail = projectAdmin.Email
		}
		responder.NewHttpResponse(r, w, http.StatusOK, project, nil)
	}
}

func DeleteProjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args := mux.Vars(r)
		i, _ := strconv.ParseUint(args["id"], 10, 64)

		err := mysql.DeleteProject(i)
		if err != nil {
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		responder.NewHttpResponse(r, w, http.StatusOK, "success", nil)
	}
}

func UpdateProjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var project mysql.Project
		args := mux.Vars(r)
		payloads, err := ioutil.ReadAll(r.Body)
		if err != nil {
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		json.Unmarshal(payloads, &project)

		i, _ := strconv.ParseUint(args["id"], 10, 64)

		err = mysql.UpdateProject(uint(i), project)
		if err != nil {
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		projectUpdated, err := mysql.ReadProject(i)
		if err != nil {
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		responder.NewHttpResponse(r, w, http.StatusOK, projectUpdated, nil)
	}
}
