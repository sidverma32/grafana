package channels

import (
	"context"
	"path"
	"time"

	gokit_log "github.com/go-kit/kit/log"
	"github.com/prometheus/alertmanager/notify"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/alertmanager/types"
	"github.com/prometheus/common/model"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/alerting"
	old_notifiers "github.com/grafana/grafana/pkg/services/alerting/notifiers"
	"github.com/grafana/grafana/pkg/setting"
)

const (
	// victoropsAlertStateCritical - Victorops uses "CRITICAL" string to indicate "Alerting" state
	victoropsAlertStateCritical = "CRITICAL"

	// victoropsAlertStateWarning - VictorOps "WARNING" message type
	// victoropsAlertStateWarning = "WARNING"

	// victoropsAlertStateRecovery - VictorOps "RECOVERY" message type
	victoropsAlertStateRecovery = "RECOVERY"
)

// NewVictoropsNotifier creates an instance of VictoropsNotifier that
// handles posting notifications to Victorops REST API
func NewVictoropsNotifier(model *models.AlertNotification, t *template.Template) (*VictoropsNotifier, error) {
	url := model.Settings.Get("url").MustString()
	if url == "" {
		return nil, alerting.ValidationError{Reason: "Could not find victorops url property in settings"}
	}

	return &VictoropsNotifier{
		NotifierBase: old_notifiers.NewNotifierBase(model),
		URL:          url,
		log:          log.New("alerting.notifier.victorops"),
		tmpl:         t,
	}, nil
}

// VictoropsNotifier defines URL property for Victorops REST API
// and handles notification process by formatting POST body according to
// Victorops specifications (http://victorops.force.com/knowledgebase/articles/Integration/Alert-Ingestion-API-Documentation/)
type VictoropsNotifier struct {
	old_notifiers.NotifierBase
	URL  string
	log  log.Logger
	tmpl *template.Template
}

// Notify sends notification to Victorops via POST to URL endpoint
func (vn *VictoropsNotifier) Notify(ctx context.Context, as ...*types.Alert) (bool, error) {
	vn.log.Debug("Executing victorops notification", "notification", vn.Name)

	alerts := types.Alerts(as...)
	// Default to alerting and change based on state checks (Ensures string type)
	// TODO: how to do warnings? Should the default state be a configuration?
	messageType := victoropsAlertStateCritical
	if alerts.Status() == model.AlertResolved {
		messageType = victoropsAlertStateRecovery
	}

	// TODO: to be removed after figuring out WARNING. This is from 7.x.
	//for _, tag := range evalContext.Rule.AlertRuleTags {
	//	if strings.ToLower(tag.Key) == "severity" {
	//		// Only set severity if it's one of the PD supported enum values
	//		// Info, Warning, Error, or Critical (case insensitive)
	//		switch sev := strings.ToUpper(tag.Value); sev {
	//		case "INFO":
	//			fallthrough
	//		case "WARNING":
	//			fallthrough
	//		case "CRITICAL":
	//			messageType = sev
	//		default:
	//			vn.log.Warn("Ignoring invalid severity tag", "severity", sev)
	//		}
	//	}
	//}

	data := notify.GetTemplateData(ctx, vn.tmpl, as, gokit_log.NewNopLogger())
	var tmplErr error
	tmpl := notify.TmplText(vn.tmpl, data, &tmplErr)

	bodyJSON := simplejson.New()
	bodyJSON.Set("message_type", messageType)
	bodyJSON.Set("entity_id", "TODO") // TODO: not sure what ID to give. It was the rule name before.
	bodyJSON.Set("entity_display_name", tmpl(`{{ template "default.title" . }}`))
	bodyJSON.Set("timestamp", time.Now().Unix())
	bodyJSON.Set("state_message", tmpl(`{{ template "default.message" . }}`))
	bodyJSON.Set("monitoring_tool", "Grafana v"+setting.BuildVersion)
	bodyJSON.Set("alert_url", path.Join(vn.tmpl.ExternalURL.String(), "/alerting/list"))

	// Removed in 8.x.
	//bodyJSON.Set("metrics", fields)
	//bodyJSON.Set("state_start_time", evalContext.StartTime.Unix())

	b, err := bodyJSON.MarshalJSON()
	if err != nil {
		return false, err
	}
	cmd := &models.SendWebhookSync{
		Url:  vn.URL,
		Body: string(b),
	}

	if err := bus.DispatchCtx(ctx, cmd); err != nil {
		vn.log.Error("Failed to send Victorops notification", "error", err, "webhook", vn.Name)
		return false, err
	}

	return true, nil
}

func (vn *VictoropsNotifier) SendResolved() bool {
	return !vn.GetDisableResolveMessage()
}
