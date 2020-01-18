/**
 * Copyright 2019 Innodev LLC. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package hook

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/logging"
	joonix "github.com/innodv/log"
	"github.com/sirupsen/logrus"
)

func NewErrors(project string, application string) (logrus.Hook, error) {
	client, err := logging.NewClient(context.Background(), project)
	if err != nil {
		return nil, err
	}

	return &errHandler{
		lvls: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
		client: client,
		format: joonix.NewFormatter(),
		logger: client.Logger(application),
	}, nil
}

type errHandler struct {
	lvls   []logrus.Level
	client *logging.Client
	logger *logging.Logger
	format logrus.Formatter
}

func (eh *errHandler) Levels() []logrus.Level {
	return eh.lvls
}

func getSeverity(entry *logrus.Entry) logging.Severity {
	switch entry.Level {
	case logrus.PanicLevel:
		return logging.Critical
	case logrus.FatalLevel:
		return logging.Emergency
	case logrus.ErrorLevel:
		return logging.Error
	default:
		return logging.Warning
	}
}

func (eh *errHandler) Fire(entry *logrus.Entry) error {
	data, err := eh.format.Format(entry)
	if err != nil {
		return err
	}
	eh.logger.Log(logging.Entry{Severity: getSeverity(entry), Payload: json.RawMessage(data)})
	return nil
}
