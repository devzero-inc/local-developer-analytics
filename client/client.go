package client

import (
	"context"
	"fmt"
	gen "lda/gen/api/v1"
	"lda/logging"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Config holds configuration for the client connection.
type Config struct {
	Address          string // The server address
	SecureConnection bool   // True for secure (HTTPS), false for insecure (HTTP)
	CertFile         string // Optional path to the TLS cert file for secure connections
	Timeout          int    // Timeout in seconds for the connection
}

// Client is a struct that holds the connection to the server
type Client struct {
	conn    *grpc.ClientConn
	logger  *zerolog.Logger
	timeout time.Duration
}

// NewClient creates a new client and returns a pointer to it and an error
func NewClient(config Config) (*Client, error) {
	var opts []grpc.DialOption

	// Setup connection security based on config
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	if config.SecureConnection {
		if config.CertFile != "" {
			tlsFromFile, err := credentials.NewClientTLSFromFile(config.CertFile, "")
			if err != nil {
				return nil, fmt.Errorf("failed to create TLS credentials: %w", err)
			}
			creds = grpc.WithTransportCredentials(tlsFromFile)
		} else {
			creds = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
		}
	}
	opts = append(opts, creds)

	conn, err := grpc.Dial(config.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	// Set a default timeout of 60 seconds if not provided
	if config.Timeout == 0 {
		config.Timeout = 60
	}

	client := &Client{
		conn:    conn,
		logger:  &logging.Log,
		timeout: time.Duration(config.Timeout) * time.Second,
	}

	return client, nil
}

// SendCommands sends a list of commands to the server
func (c *Client) SendCommands(commands []*gen.Command, auth *gen.Auth) error {

	client := gen.NewCollectorServiceClient(c.conn)

	req := &gen.SendCommandsRequest{
		Commands: commands,
		Auth:     auth,
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := client.SendCommands(ctx, req)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to send commands")
	}

	return err
}

// SendProcesses sends a list of processes to the server
func (c *Client) SendProcesses(processes []*gen.Process, auth *gen.Auth) error {

	client := gen.NewCollectorServiceClient(c.conn)

	req := &gen.SendProcessesRequest{
		Processes: processes,
		Auth:      auth,
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := client.SendProcesses(ctx, req)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to send processes")
	}

	return err
}

// Close closes the connection to the server
func (c *Client) Close() {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			c.logger.Error().Err(err).Msg("Failed to close the connection")
		} else {
			c.logger.Info().Msg("Connection closed successfully")
		}
	}
}
