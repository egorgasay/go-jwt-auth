package lib

import (
	"fmt"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type Logger struct {
	Sugared *zap.SugaredLogger
	*zap.Logger
}

// LogEvent log event for fx logger
func (l *Logger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Sugared.Debug("OnStart hook executing: ",
			zap.String("callee", e.FunctionName),
			zap.String("caller", e.CallerName),
		)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Sugared.Debug("OnStart hook failed: ",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.Error(e.Err),
			)
		} else {
			l.Sugared.Debug("OnStart hook executed: ",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.OnStopExecuting:
		l.Sugared.Debug("OnStop hook executing: ",
			zap.String("callee", e.FunctionName),
			zap.String("caller", e.CallerName),
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Sugared.Debug("OnStop hook failed: ",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.Error(e.Err),
			)
		} else {
			l.Sugared.Debug("OnStop hook executed: ",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.Supplied:
		l.Sugared.Debug("supplied: ", zap.String("type", e.TypeName), zap.Error(e.Err))
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Sugared.Debug("provided: ", e.ConstructorName, " => ", rtype)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.Sugared.Debug("decorated: ",
				zap.String("decorator", e.DecoratorName),
				zap.String("type", rtype),
			)
		}
	case *fxevent.Invoking:
		l.Sugared.Debug("invoking: ", e.FunctionName)
	case *fxevent.Started:
		if e.Err == nil {
			l.Sugared.Debug("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err == nil {
			l.Sugared.Debug("initialized: custom fxevent.Sugared -> ", e.ConstructorName)
		}
	}
}

func NewLogger() (Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return Logger{}, fmt.Errorf("can't create Logger instance: %v", err)
	}

	return Logger{Sugared: logger.Sugar(), Logger: logger}, nil
}
