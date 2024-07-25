package entity

import "time"

type MaestroTask struct {
	Id                int64     `json:"id"`
	SlidingWindowSize int       `json:"slidingWindowSize"`
	FailThreshold     int       `json:"failThreshold"`
	CommandId         int64     `json:"commandId"`
	Command           string    `json:"command"`
	CommandType       int64     `json:"commandType"`
	Status            string    `json:"status"`
	CreationDate      time.Time `json:"creationDate"`
	CountSucceeded    int       `json:"countSucceeded"`
	UserFullName      string    `json:"userFullName"`
	StartDate         time.Time `json:"startDate"`
	EndDate           time.Time `json:"endDate"`
	Duration          int64     `json:"duration"`
}

type MaestroTenant struct {
	TenantId        string `json:"tenantId"`
	TenantName      string `json:"tenantName"`
	EnvironmentName string `json:"environmentName"`
	Stage           string `json:"stage"`
	EnvironmentType string `json:"environmentType"`
	Owner           string `json:"owner"`
}

type MaestroTenantResponse struct {
	Items []*MaestroTenant `json:"items"`
	Count int              `json:"count"`
}

type MaestroTaskResponse struct {
	Items []*MaestroTask `json:"items"`
}
