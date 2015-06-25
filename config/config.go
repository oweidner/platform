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

import "fmt"

// The Config struct defines the structure of the configuration file.
//
type Config struct {
	Server struct {
		Listen        string `gcfg:"listen"`
		APIPrefix     string
		DisableAuth   bool   `gcfg:"disableauth"`
		AdminAccount  string `gcfg:"adminaccount"`
		AdminPassword string `gcfg:"adminpassword"`
	}
	TLS struct {
		EnableTLS bool
		KeyFile   string
		CertFile  string
	}

	JWT struct {
		Expiration int    `gcfg:"expiration"`
		PublicKey  string `gcfg:"publickey"`
		PrivateKey string `gcfg:"privatekey"`
	}
	MySQL struct {
		Host     string
		Database string
		Username string
		Password string
	}
}

// CheckConfig checks the configuration file values and sets defaults
// wherever necessary.
func CheckConfig(config *Config, filename string) error {

	// Set the default to 12 hours if JWT:Expiration is not defined.
	if config.Server.APIPrefix == "" {
		config.Server.APIPrefix = "platform"
	}

	// Return an error if TLS is enabled and cert or Key are not provided
	if config.TLS.EnableTLS == true {
		if config.TLS.KeyFile == "" {
			return fmt.Errorf("%v: EnableTLS is set, but KeyFile is not defined", filename)
		}
		if config.TLS.CertFile == "" {
			return fmt.Errorf("%v: EnableTLS is set, but CertFile is not defined", filename)
		}

	}

	// Return an error if no Server.StorageBackend is defined.
	if config.Server.Listen == "" {
		return fmt.Errorf("%v: Configuration doesn't define mandatory Server.Listen", filename)
	}

	// Return an error if no Server.StorageBackend is defined.
	if config.Server.AdminAccount == "" {
		return fmt.Errorf("%v: Configuration doesn't define mandatory Server.AdminAccount", filename)
	}
	// Return an error if no Server.StorageBackend is defined.
	if config.Server.AdminPassword == "" {
		return fmt.Errorf("%v: Configuration doesn't define mandatory Server.AdminPassword", filename)
	}

	// Return an error if no JWT Public Key is defined.
	if config.JWT.PublicKey == "" {
		return fmt.Errorf("%v: Configuration doesn't define mandatory JWT.public_key", filename)
	}
	// Return an error if no JWT Private Key is defined.
	if config.JWT.PrivateKey == "" {
		return fmt.Errorf("%v: Configuration doesn't define mandatory JWT.private_key", filename)
	}

	// Set the default to 12 hours if JWT:Expiration is not defined.
	if config.JWT.Expiration <= 0 {
		config.JWT.Expiration = 12
	}

	// if config.MySQL.Host == "" {
	// 	return fmt.Errorf("%v: 'mongodb' storage backend requires MySQL.host config option", filename)
	// }
	if config.MySQL.Database == "" {
		return fmt.Errorf("%v: 'mongodb' storage backend requires MySQL.database config option", filename)
	}

	return nil
}
