package controller

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/entity"
	"advanced-tools/pkg/vars"
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
			EndsAt:    time.Now().Add(20 * time.Minute),
			CreatedBy: "eks-upgrade-manager",
			Comment:   "silencing alerts during eks upgrade",
		}
		err := controller.clients.AlertManagerClient.CreateSilence(&silence)
		if err != nil {
			return err
		}
		log.Info().Msgf("silence for [%v] [%v] created successfully", alertName, vars.Environment)
	}
	return nil
}
