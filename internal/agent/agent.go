package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/gitwatch"
	"github.com/xaque208/znet/pkg/events"
)

// Agent is an RPC client worker bee.
type Agent struct {
	config     *config.Config
	conn       *grpc.ClientConn
	grpcServer *grpc.Server
}

// NewAgent returns a new *Agent from the given arguments.
func NewAgent(cfg *config.Config, conn *grpc.ClientConn) *Agent {

	if cfg.TLS == nil {
		log.Warn("nil TLS config")
	}

	if cfg.Vault == nil {
		log.Warn("nil Vault config")
	}

	if cfg.RPC == nil {
		log.Warn("nil RPC config")
	}

	if cfg.Agent == nil {
		log.Warn("nil Agent config")
	}

	a := &Agent{
		config: cfg,
		conn:   conn,
	}

	if cfg.RPC != nil {
		a.grpcServer = comms.StandardRPCServer(cfg.Vault, cfg.TLS)
	}

	return a
}

func (a *Agent) namedTimerHandler(name string, payload events.Payload) error {
	log.WithFields(log.Fields{
		"name":    name,
		"payload": string(payload),
	}).Warn("TODO")

	return nil
}

func (a *Agent) newTagHandler(name string, payload events.Payload) error {
	var x gitwatch.NewTag

	err := json.Unmarshal(payload, &x)
	if err != nil {
		log.Errorf("failed to unmarshal %T: %s", x, err)
	}

	return a.executeForGitEvent(x)
}

func (a *Agent) newCommitHandler(name string, payload events.Payload) error {
	log.Debugf("Agent.newCommitHandler: %+v", string(payload))
	log.Debugf("Agent.newCommitHandler config: %+v", a.config)

	var x gitwatch.NewCommit

	err := json.Unmarshal(payload, &x)
	if err != nil {
		log.Errorf("failed to unmarshal %T: %s", x, err)
	}

	return a.executeForGitEvent(x)
}

func (a *Agent) executeForGitEvent(x interface{}) error {
	log.Tracef("executeForGitEvent %+v", x)

	for _, execution := range a.config.Agent.Executions {

		for _, xx := range execution.Events {
			if xx != "" {
				var args []string

				// Render the args as template strings, passing the current x interface.
				for _, v := range execution.Args {
					tmpl, err := template.New("env").Parse(v)
					if err != nil {
						log.Errorf("failed to parse template %s: %s", v, err)
					}

					var buf bytes.Buffer

					err = tmpl.Execute(&buf, x)
					if err != nil {
						log.Error(err)
					}

					args = append(args, buf.String())
				}

				cmd := exec.Command(execution.Command, args...)

				if execution.Dir != "" {
					cmd.Dir = execution.Dir
				}

				var env []string

				// Render the values of the environment variables as templates using
				// the received event.
				for k, v := range execution.Environment {

					tmpl, err := template.New("env").Parse(v)
					if err != nil {
						log.Errorf("failed to parse template %s: %s", v, err)
					}

					var buf bytes.Buffer

					err = tmpl.Execute(&buf, x)
					if err != nil {
						log.Error(err)
					}

					env = append(env, fmt.Sprintf("%s=%s", k, buf.String()))
				}

				if len(env) > 0 {
					cmd.Env = append(os.Environ(), env...)
				}

				start := time.Now()
				// var out bytes.Buffer
				// cmd.Stdout = &out
				output, err := cmd.CombinedOutput()
				if err != nil {
					log.Errorf("command execution failed: %s", err)
				}

				now := time.Now()

				ev := ExecutionResult{
					Time:     &now,
					Command:  execution.Command,
					Args:     args,
					Dir:      execution.Dir,
					Output:   output,
					ExitCode: cmd.ProcessState.ExitCode(),
					Duration: time.Since(start),
				}

				err = events.ProduceEvent(a.conn, ev)
				if err != nil {
					log.Error(err)
				}

			}
		}
	}

	return nil
}

// Start calls start on the agent gRPC server.
func (a *Agent) Start() error {
	if a.config.RPC == nil {
		return fmt.Errorf("unable to start agent with nil RPC config")
	}

	if a.config.RPC.AgentListenAddress != "" {
		log.WithFields(log.Fields{
			"rpc_listen": a.config.RPC.AgentListenAddress,
		}).Debug("starting RPC listener")

		err := a.startRPCListener()
		if err != nil {
			return err
		}
	}

	return nil
}

// Stop calls stop on the agent gRPC server.
func (a *Agent) Stop() error {
	if a.grpcServer != nil {
		log.Debug("stopping RPC listener")
		a.grpcServer.Stop()
	}
	return nil
}

func (a *Agent) startRPCListener() error {
	info, err := readOSRelease()
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"id": info.ID,
	}).Debug("os-release")

	switch info.ID {
	case "freebsd":
		RegisterJailHostServer(a.grpcServer, &jailServer{})
		RegisterNodeServer(a.grpcServer, &nodeServer{})
	case "arch":
		RegisterJailHostServer(a.grpcServer, &notImplementedJailServer{})
		RegisterNodeServer(a.grpcServer, &nodeServer{})
	}

	go func() {
		lis, err := net.Listen("tcp", a.config.RPC.AgentListenAddress)
		if err != nil {
			log.Errorf("rpc failed to listen: %s", err)
		}

		err = a.grpcServer.Serve(lis)
		if err != nil {
			log.Error(err)
		}
	}()

	return nil
}
