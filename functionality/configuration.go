package functionality

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/doublen987/Projects/MySite/server/persistence"
)

var (
	DBTypeDefault              = persistence.MONGODB
	DBConnectionDefault        = "mongodb://127.0.0.1:5656"
	FileStorageTypeDefault     = persistence.FILESYSTEM
	StoragePathDefault         = "images"
	HostDefault                = ""
	PortDefault                = "80"
	CertPEM                    = "cert.pem"
	KeyPEM                     = "key.pem"
	HostTLSDefault             = ""
	TLSPortDefault             = "443"
	MessageBrokerTypeDefault   = "amqp"
	AMQPMessageBrokerDefault   = "amqp://guest:guest@localhost:5672"
	KafkaMessageBrokersDefault = []string{"localhost:9092"}
)

type ServiceConfig struct {
	Databasetype        uint8    `json:"databasetype"`
	DBConnection        string   `json:"dbconnection"`
	FileStorageType     string   `json:"filestoragetype"`
	Storagepath         string   `json:"storagepath"`
	Host                string   `json:"host"`
	Port                string   `json:"port"`
	CertPEM             string   `json:"cert_pem"`
	KeyPEM              string   `json:"key_pem"`
	HostTLS             string   `json:"host_tls"`
	PortTLS             string   `json:"port_tls"`
	MessageBrokerType   string   `json:"message_broker_type"`
	AMQPMessageBroker   string   `json:"amqp_message_broker"`
	KafkaMessageBrokers []string `json:"kafka_message_brokers"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		FileStorageTypeDefault,
		StoragePathDefault,
		HostDefault,
		PortDefault,
		CertPEM,
		KeyPEM,
		HostTLSDefault,
		TLSPortDefault,
		MessageBrokerTypeDefault,
		AMQPMessageBrokerDefault,
		KafkaMessageBrokersDefault,
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values.")
		return conf, err
	}

	err = json.NewDecoder(file).Decode(&conf)
	if broker := os.Getenv("AMQP_URL"); broker != "" {
		conf.AMQPMessageBroker = broker
	}
	if dbURL := os.Getenv("DB_URL"); dbURL != "" {
		conf.DBConnection = dbURL
	}
	if port := os.Getenv("PORT"); port != "" {
		conf.Port = port
	}
	if portTLS := os.Getenv("TLS_PORT"); portTLS != "" {
		conf.PortTLS = portTLS
	}
	if CertPEM := os.Getenv("CERT_PEM"); CertPEM != "" {
		conf.CertPEM = CertPEM
	}
	if KeyPEM := os.Getenv("KEY_PEM"); KeyPEM != "" {
		conf.KeyPEM = KeyPEM
	}
	if DBType := os.Getenv("DB_TYPE"); DBType != "" {
		if i, err := strconv.ParseUint(DBType, 10, 8); err != nil {
			conf.Databasetype = uint8(i)
		}
	}
	if FileStorageType := os.Getenv("FILE_STORAGE_TYPE"); FileStorageType != "" {
		conf.FileStorageType = FileStorageType
	}
	return conf, err
}
