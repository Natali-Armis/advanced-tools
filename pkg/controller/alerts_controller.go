package controller

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/entity"
	"advanced-tools/pkg/vars"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type AlertsController struct {
	clients *client.Client
}

func GetAlertsController(clients *client.Client) *AlertsController {
	return &AlertsController{
		clients: clients,
	}
}

func (controller *AlertsController) SilenceAlerts() error {
	for _, alertName := range vars.AlertsToSilenceDuringUpgrade {
		silence := entity.Silence{
			Matchers: []entity.Matcher{
				{
					Name:    vars.ALERT_NAME_LABEL,
					Value:   alertName,
					IsRegex: false,
				},
				{
					Name:    vars.ALERT_ENV_LABEL,
					Value:   vars.Environment,
					IsRegex: false,
				},
			},
			StartsAt:  time.Now(),
			EndsAt:    time.Now().Add(30 * time.Minute),
			CreatedBy: "eks-upgrade-manager",
			Comment:   "silencing alerts during eks upgrade",
		}
		err := controller.clients.AlertManagerClient.CreateSilence(&silence)
		if err != nil {
			return err
		}
		log.Info().Msgf("silence for [%v] [%v] created successfully", alertName, vars.Environment)
	}
	excludedTeamsPattern := buildExcludedTeamsPattern(vars.AlertsTeamsToNotSilence)
	silence := entity.Silence{
		Matchers: []entity.Matcher{
			{
				Name:    "team",
				Value:   excludedTeamsPattern,
				IsRegex: true, 
			},
			{
				Name:    vars.ALERT_ENV_LABEL,
				Value:   vars.Environment,
				IsRegex: false,
			},
		},
		StartsAt:  time.Now(),
		EndsAt:    time.Now().Add(30 * time.Minute),
		CreatedBy: "eks-upgrade-manager",
		Comment:   "preventively silencing alerts for teams not in the exclusion list during eks upgrade",
	}
	err := controller.clients.AlertManagerClient.CreateSilence(&silence)
	if err != nil {
		return err
	}
	log.Info().Msgf("preventive silence for non-excluded teams in [%v] created successfully", vars.Environment)
	return nil
}

func buildExcludedTeamsPattern(excludedTeams []string) string {
	joinedTeams := strings.Join(excludedTeams, "|")
	return "^(?!(" + joinedTeams + ")$).*"
}
