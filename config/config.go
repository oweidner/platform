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

import (
	"fmt"
	"path"
)

// The Config struct defines the structure of the configuration file.
//
type Config struct {
	SERVER struct {
		Listen               string //`gcfg:"listen"`
		PlatformPrefix       string
		ApplicationPrefix    string
		EnablePlatformAPI    bool
		EnableApplicationAPI bool
		DisableAuth          bool   //`gcfg:"disableauth"`
		AdminAccount         string //`gcfg:"adminaccount"`
		AdminPassword        string //`gcfg:"adminpassword"`
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
func CheckConfig(config *Config, filename string, basepath string) error {

	// Set the default platform prefix.
	if config.SERVER.PlatformPrefix == "" {
		config.SERVER.PlatformPrefix = "platform"
	}

	// Set the default application prefix.
	if config.SERVER.ApplicationPrefix == "" {
		config.SERVER.ApplicationPrefix = "app"
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

	// Return an error if no SERVER.StorageBackend is defined.
	if config.SERVER.Listen == "" {
		return fmt.Errorf("%v: Configuration doesn't define mandatory SERVER.Listen", filename)
	}

	// Return an error if no SERVER.StorageBackend is defined.
	if config.SERVER.AdminAccount == "" {
		return fmt.Errorf("%v: Configuration doesn't define mandatory SERVER.AdminAccount", filename)
	}
	// Return an error if no SERVER.StorageBackend is defined.
	if config.SERVER.AdminPassword == "" {
		return fmt.Errorf("%v: Configuration doesn't define mandatory SERVER.AdminPassword", filename)
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
		return fmt.Errorf("%v: 'MYSQL.database is a required option", filename)
	}

	if config.MySQL.Host == "" {
		return fmt.Errorf("%v: 'MYSQL.host required option", filename)
	}

	if config.MySQL.Username == "" {
		return fmt.Errorf("%v: 'MYSQL.username is a required option", filename)
	}

	// tcpAddr := os.Getenv("MYSQL_PORT_3306_TCP_ADDR")
	// tcpPort := os.Getenv("MYSQL_PORT_3306_TCP_PORT")
	//
	// if tcpAddr != "" && tcpPort != "" {
	// 	fmt.Printf("$MYSQL_PORT_3306_TCP_ADDR and $MYSQL_PORT_3306_TCP_PORT set. Overrides MYSQL.host in %v", filename)
	// 	config.MySQL.Host = fmt.Sprintf("tcp(%v:%v)", tcpAddr, tcpPort)
	// }

	// check if various paths are relative. if so, add the base path
	if path.IsAbs(config.JWT.PublicKey) == false {
		config.JWT.PublicKey = path.Join(basepath, config.JWT.PublicKey)
	}
	if path.IsAbs(config.JWT.PrivateKey) == false {
		config.JWT.PrivateKey = path.Join(basepath, config.JWT.PrivateKey)
	}
	if path.IsAbs(config.TLS.KeyFile) == false {
		config.TLS.KeyFile = path.Join(basepath, config.TLS.KeyFile)
	}
	if path.IsAbs(config.TLS.CertFile) == false {
		config.TLS.CertFile = path.Join(basepath, config.TLS.CertFile)
	}

	return nil
}
