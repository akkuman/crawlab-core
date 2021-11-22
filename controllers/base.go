package controllers

import "github.com/gin-gonic/gin"

const (
	ControllerIdNode = iota
	ControllerIdProject
	ControllerIdSpider
	ControllerIdTask
	ControllerIdJob
	ControllerIdSchedule
	ControllerIdUser
	ControllerIdSetting
	ControllerIdToken
	ControllerIdVariable
	ControllerIdTag
	ControllerIdLogin
	ControllerIdColor
	ControllerIdPlugin
	ControllerIdDataSource
	ControllerIdDataCollection
	ControllerIdResult
	ControllerIdStats
	ControllerIdFiler
	ControllerIdPluginDo
	ControllerIdGit
	ControllerIdVersion
)

type ControllerId int

type BasicController interface {
	Get(c *gin.Context)
	Post(c *gin.Context)
	Put(c *gin.Context)
	Delete(c *gin.Context)
}

type ListController interface {
	BasicController
	GetList(c *gin.Context)
	PutList(c *gin.Context)
	PostList(c *gin.Context)
	DeleteList(c *gin.Context)
}

type Action struct {
	Method      string
	Path        string
	HandlerFunc gin.HandlerFunc
}

type ActionController interface {
	Actions() (actions []Action)
}

type ListActionController interface {
	ListController
	ActionController
}
