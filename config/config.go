//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package config

import "errors"

// The Config struct defines the structure of the configuration file.
//
type Config struct {
	Server struct {
		Listen        string `gcfg:"listen"`
		DisableAuth   bool   `gcfg:"disableauth"`
		AdminUser     string `gcfg:"adminuser"`
		AdminPassword string `gcfg:"adminpassword"`
	}
	JWT struct {
		Expiration int    `gcfg:"expiration"`
		PublicKey  string `gcfg:"publickey"`
		PrivateKey string `gcfg:"privatekey"`
	}
	MySQL struct {
		Host     string `gcfg:"host"`
		Database string `gcfg:"database"`
	}
}

// CheckConfig checks the configuration file values and sets defaults
// wherever necessary.
func CheckConfig(config *Config) error {

	// Return an error if no Server.StorageBackend is defined.
	if config.Server.Listen == "" {
		return errors.New("config: Configuration doesn't define mandatory Server.Listen")
	}

	// Return an error if no Server.StorageBackend is defined.
	if config.Server.AdminUser == "" {
		return errors.New("config: Configuration doesn't define mandatory Server.AdminUser")
	}
	// Return an error if no Server.StorageBackend is defined.
	if config.Server.AdminPassword == "" {
		return errors.New("config: Configuration doesn't define mandatory Server.AdminPassword")
	}

	// Return an error if no JWT Public Key is defined.
	if config.JWT.PublicKey == "" {
		return errors.New("config: Configuration doesn't define mandatory JWT.public_key")
	}
	// Return an error if no JWT Private Key is defined.
	if config.JWT.PrivateKey == "" {
		return errors.New("config: Configuration doesn't define mandatory JWT.private_key")
	}

	// Set the default to 12 hours if JWT:Expiration is not defined.
	if config.JWT.Expiration <= 0 {
		config.JWT.Expiration = 12
	}

	if config.MySQL.Host == "" {
		return errors.New("config: 'mongodb' storage backend requires MySQL.host config option")
	}
	if config.MySQL.Database == "" {
		return errors.New("config: 'mongodb' storage backend requires MySQL.database config option")
	}

	return nil
}
