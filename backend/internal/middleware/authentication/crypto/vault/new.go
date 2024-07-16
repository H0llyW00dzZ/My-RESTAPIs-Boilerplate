// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Note: Currently unused, and marked As TODO, previously it used for secure internal database.

package vault

import (
	"context"
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
)

// VConfig holds configuration parameters for the Vault client.
type VConfig struct {
	// Connection parameters
	Address             string
	ApproleRoleID       string
	ApproleSecretIDFile string
	Token               string
	// Transit secrets engine mount path
	TransitPath string
}

// VClient manages interaction with a Hashicorp Vault instance.
type VClient struct {
	client     *api.Client
	parameters VConfig
}

// NewVaultAppRoleClient creates a new Vault client authenticated using AppRole.
func NewVaultAppRoleClient(ctx context.Context, parameters VConfig) (*VClient, error) {
	log.LogInfof("connecting to vault @ %s", parameters.Address)

	config := api.DefaultConfig()
	config.Address = parameters.Address
	config.Timeout = 10 * time.Second

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize vault client: %w", err)
	}

	vault := &VClient{
		client:     client,
		parameters: parameters,
	}

	if err := vault.login(ctx); err != nil {
		return nil, fmt.Errorf("vault login error: %w", err)
	}

	log.LogInfof("connecting to vault: success!")
	return vault, nil
}

// login authenticates the Vault client using AppRole.
func (v *VClient) login(ctx context.Context) error {
	log.LogInfof("logging in to vault with approle auth; role id: %s", v.parameters.ApproleRoleID)

	approleSecretID := &approle.SecretID{
		FromFile: v.parameters.ApproleSecretIDFile,
	}

	appRoleAuth, err := approle.NewAppRoleAuth(
		v.parameters.ApproleRoleID,
		approleSecretID,
		approle.WithWrappingToken(),
	)
	if err != nil {
		return fmt.Errorf("unable to initialize approle authentication method: %w", err)
	}

	authInfo, err := v.client.Auth().Login(ctx, appRoleAuth)
	if err != nil {
		return fmt.Errorf("unable to login using approle auth method: %w", err)
	}
	if authInfo == nil {
		return fmt.Errorf("no approle info was returned after login")
	}

	log.LogInfof("logging in to vault with approle auth: success!")
	return nil
}
