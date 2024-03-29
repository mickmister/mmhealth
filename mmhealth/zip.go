package mmhealth

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/coltoneshaw/mmhealth/mmhealth/types"
	"github.com/mattermost/mattermost/server/public/model"
	"gopkg.in/yaml.v3"
)

func UnzipToMemory(zipReader *zip.Reader) (*types.PacketData, error) {

	fileContents := &types.PacketData{}

	// Read all data from the io.Reader into a byte slice

	for _, file := range zipReader.File {
		// Open each file in the zip archive
		zippedFile, err := file.Open()
		if err != nil {
			return nil, err
		}

		fmt.Println("Processing file: ", file.Name)

		defer zippedFile.Close()

		filename := filepath.Base(file.Name)

		switch filename {
		case "sanitized_config.json":
			config, err := processConfigFile(zippedFile)
			if err != nil {
				return nil, err
			}
			fileContents.Config = config
		case "plugins.json":
			plugins, err := processPluginFile(zippedFile)
			if err != nil {
				fmt.Println("Error processing plugins: ", err)
				return nil, err
			}
			fileContents.Plugins = plugins
		case "mattermost.log":
			logs, err := processMattermostLog(zippedFile)
			if err != nil {
				return nil, err
			}
			fileContents.Logs = logs
		case "notification.log":
			notifLogs, err := processNotificationLog(zippedFile)
			if err != nil {
				return nil, err
			}
			fileContents.NotificationLogs = notifLogs

		case "support_packet.yaml":
			packet, err := processPacketFile(zippedFile)
			if err != nil {
				return nil, err
			}
			fileContents.Packet = packet
		}

	}
	return fileContents, nil
}

func processConfigFile(file io.Reader) (model.Config, error) {
	var config model.Config
	err := json.NewDecoder(file).Decode(&config)
	if err != nil {
		return model.Config{}, err
	}
	return config, nil
}

func processPluginFile(file io.Reader) (model.PluginsResponse, error) {
	var plugins model.PluginsResponse
	err := json.NewDecoder(file).Decode(&plugins)
	if err != nil {
		return model.PluginsResponse{}, err
	}
	return plugins, nil
}

func processMattermostLog(file io.Reader) ([]types.MattermostLogEntry, error) {

	logs, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	parsedLogs, err := ParseLogs(logs)

	if err != nil {
		return nil, err
	}
	return parsedLogs, nil
}

func processNotificationLog(file io.Reader) ([]byte, error) {
	logs, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func processPacketFile(file io.Reader) (model.SupportPacket, error) {
	var packet model.SupportPacket
	packetBytes, err := io.ReadAll(file)
	if err != nil {
		return model.SupportPacket{}, err
	}
	// Unmarshal the YAML into the struct
	err = yaml.Unmarshal(packetBytes, &packet)
	if err != nil {
		return model.SupportPacket{}, err
	}

	return packet, nil
}
