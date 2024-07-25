package main

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/config"
	"advanced-tools/pkg/entity"
	"advanced-tools/pkg/vars"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"html"

	"github.com/rs/zerolog/log"
)

var (
	clients *client.Client
)

const (
	RESTART    = "restart"
	INCREASE   = "increase"
	DECREASE   = "decrease"
	DEPLOYMENT = "deployment"

	APP_DEPLOY_COMMAND_ID = "304"
)

func init() {
	config.Configure()
	clients = client.GetClient()
}

type InstanceChange struct {
	Operation           string
	TenantName          string
	TenantId            string
	Environment         string
	InitialInstanceType string
	FinalInstanceType   string
	EventDate           string
}

func main() {
	oldestPostNeeded := time.Now().AddDate(0, -3, 0)
	messages, err := clients.SlackClient.FetchHistoryUpToDate(vars.SlackSTResizeNotificationsChannel, oldestPostNeeded)
	if err != nil {
		log.Fatal().Msgf("could not fetch messages from slack channel [%v] %v", vars.SlackSTResizeNotificationsChannel, err.Error())
	}

	var instanceChanges []*InstanceChange

	for _, msg := range messages {
		txt := strings.ToLower(msg.Msg.Text)
		txt = html.UnescapeString(txt)
		if len(strings.TrimSpace(txt)) == 0 {
			continue
		}

		unixTime, err := strconv.ParseFloat(msg.Timestamp, 64)
		if err != nil {
			log.Error().Msgf("could not parse timestamp [%v] %v", msg.Timestamp, err.Error())
			continue
		}
		timestamp := time.Unix(int64(unixTime), 0).Format(time.RFC3339)

		if strings.Contains(txt, DECREASE) {
			instanceChanges = append(instanceChanges, parseResizeMessage(txt, timestamp, DECREASE)...)
		} else if strings.Contains(txt, INCREASE) {
			instanceChanges = append(instanceChanges, parseResizeMessage(txt, timestamp, INCREASE)...)
		} else if strings.Contains(txt, RESTART) {
			instanceChanges = append(instanceChanges, parseRestartMessage(txt, timestamp)...)
		}
	}

	distinctTenantIds := map[string]*entity.MaestroTenant{}
	for _, instanceChange := range instanceChanges {
		tenantsResult, err := clients.MaestroClient.GetTenants(map[string]string{
			"tenantName": instanceChange.TenantName,
			"length":     "1",
		})
		if err != nil {
			log.Error().Msgf("could not fill in extra tenant data for tenant [%v] on instance change [%v] on date [%v] %v",
				instanceChange.TenantName,
				instanceChange.Operation,
				instanceChange.EventDate,
				err.Error(),
			)
			continue
		}
		if len(tenantsResult.Items) == 0 {
			log.Error().Msgf("could not fill in extra tenant data for tenant [%v] on instance change [%v] on date [%v] %v",
				instanceChange.TenantName,
				instanceChange.Operation,
				instanceChange.EventDate,
				"no such tenant by maestro results",
			)
			continue
		}
		tenant := tenantsResult.Items[0]
		instanceChange.TenantId = fmt.Sprint(tenant.TenantId)
		if len(instanceChange.Environment) == 0 {
			instanceChange.Environment = tenant.EnvironmentName
		}
		distinctTenantIds[fmt.Sprint(tenant.TenantId)] = tenant
	}

	for tenantId, tenant := range distinctTenantIds {
		fmt.Printf("fetching tasks for tenant [%v]\n", tenant.TenantName)
		tasksResult, err := clients.MaestroClient.GetTasks(map[string]string{
			"host_ids[overlap]": tenantId,
			"commandId":         APP_DEPLOY_COMMAND_ID,
			"length":            "50",
		})
		if err != nil {
			log.Error().Msgf("could not fetch deployment case tasks for tenant [%v] %v",
				tenant.TenantName,
				err.Error(),
			)
			continue
		}
		if len(tasksResult.Items) == 0 {
			continue
		}
		deployments := tasksResult.Items
		for _, deployment := range deployments {
			if deployment.CreationDate.After(oldestPostNeeded) {
				instanceChanges = append(instanceChanges, &InstanceChange{
					Operation:   DEPLOYMENT,
					TenantName:  tenant.TenantName,
					TenantId:    fmt.Sprint(tenant.TenantId),
					Environment: tenant.EnvironmentName,
					EventDate:   deployment.CreationDate.Format(time.RFC3339),
				})
			}
		}
	}

	err = writeOutAsCsv(instanceChanges)
	if err != nil {
		log.Error().Msgf("could not write out result to csv %v", err.Error())
		return
	}
}

func writeOutAsCsv(instanceChanges []*InstanceChange) error {
	file, err := os.Create("output/instance_changes.csv")
	if err != nil {
		return fmt.Errorf("could not create csv file: %w", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	header := []string{"Operation", "TenantName", "TenantId", "Environment", "InitialInstanceType", "FinalInstanceType", "EventDate"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("could not write header to csv file: %w", err)
	}
	for _, change := range instanceChanges {
		record := []string{
			change.Operation,
			change.TenantName,
			change.TenantId,
			change.Environment,
			change.InitialInstanceType,
			change.FinalInstanceType,
			change.EventDate,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("could not write record to csv file: %w", err)
		}
	}
	return nil
}

func parseResizeMessage(msg string, timestamp string, operation string) []*InstanceChange {
	var instanceChanges []*InstanceChange
	re := regexp.MustCompile(`([a-z0-9-]+)\s*\[(m[0-9a-z.]+)\s*->\s*(m[0-9a-z.]+)\]`)
	lines := strings.Split(msg, "\n")
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			instanceChanges = append(instanceChanges, &InstanceChange{
				Operation:           operation,
				TenantName:          matches[1],
				InitialInstanceType: matches[2],
				FinalInstanceType:   matches[3],
				EventDate:           timestamp,
			})
		}
	}
	return instanceChanges
}

func parseRestartMessage(msg string, timestamp string) []*InstanceChange {
	var restartCases []*InstanceChange
	re := regexp.MustCompile(`(prod[0-9]*)\s*\|\s*([a-z0-9-]+)`)
	lines := strings.Split(msg, "\n")

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			restartCases = append(restartCases, &InstanceChange{
				Operation:   RESTART,
				Environment: matches[1],
				TenantName:  matches[2],
				EventDate:   timestamp,
			})
		}
	}

	return restartCases
}
