package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"

	"github.com/astronomerio/analytics-go"
	"github.com/tj/docopt"
)

const Usage = `
Analytics Go CLI

Usage:
  analytics track <event> [--properties=<properties>] [--context=<context>] [--appId=<appId>] [--userId=<userId>] [--anonymousId=<anonymousId>] [--integrations=<integrations>] [--timestamp=<timestamp>]
  analytics screen <name> [--properties=<properties>] [--context=<context>] [--appId=<appId>] [--userId=<userId>] [--anonymousId=<anonymousId>] [--integrations=<integrations>] [--timestamp=<timestamp>]
  analytics page <name> [--properties=<properties>] [--context=<context>] [--appId=<appId>] [--userId=<userId>] [--anonymousId=<anonymousId>] [--integrations=<integrations>] [--timestamp=<timestamp>]
  analytics identify [--traits=<traits>] [--context=<context>] [--appId=<appId>] [--userId=<userId>] [--anonymousId=<anonymousId>] [--integrations=<integrations>] [--timestamp=<timestamp>]
  analytics group --groupId=<groupId> [--traits=<traits>] [--properties=<properties>] [--context=<context>] [--appId=<appId>] [--userId=<userId>] [--anonymousId=<anonymousId>] [--integrations=<integrations>] [--timestamp=<timestamp>]
  analytics alias --userId=<userId> --previousId=<previousId> [--traits=<traits>] [--properties=<properties>] [--context=<context>] [--appId=<appId>] [--anonymousId=<anonymousId>] [--integrations=<integrations>] [--timestamp=<timestamp>]
  analytics -h | --help
  analytics --version

Options:
  -h --help     Show this screen.
  --version     Show version.
`

func main() {
	arguments, err := docopt.Parse(Usage, nil, true, "Anaytics Go CLI", false)
	check(err)

	appId := getOptionalString(arguments, "--appId")
	if appId == "" {
		appId = os.Getenv("ASTRONOMER_APP_ID")
		if appId == "" {
			log.Fatal("either $ASTRONOMER_APP_ID or --appId must be provided")
		}
	}

	client := analytics.New(appId)
	client.Size = 1
	client.Verbose = true

	if arguments["track"].(bool) {
		m := &analytics.Track{
			Event: arguments["<event>"].(string),
		}
		properties := getOptionalString(arguments, "--properties")
		if properties != "" {
			var parsedProperties map[string]interface{}
			err := json.Unmarshal([]byte(properties), &parsedProperties)
			check(err)
			m.Properties = parsedProperties
		}

		setCommonFields(m, arguments)

		check(client.Track(m))
	}

	if arguments["screen"].(bool) || arguments["page"].(bool) {
		m := &analytics.Page{
			Name: arguments["<name>"].(string),
		}
		/* Bug in Go library - page has traits not properties.
		properties := getOptionalString(arguments, "--properties")
		if properties != "" {
			var parsedProperties map[string]interface{}
			err := json.Unmarshal([]byte(properties), &parsedProperties)
			check(err)
			t.Properties = parsedProperties
		}
		*/

		setCommonFields(m, arguments)

		check(client.Page(m))
	}

	if arguments["identify"].(bool) {
		m := &analytics.Identify{}
		traits := getOptionalString(arguments, "--traits")
		if traits != "" {
			var parsedTraits map[string]interface{}
			err := json.Unmarshal([]byte(traits), &parsedTraits)
			check(err)
			m.Traits = parsedTraits
		}

		setCommonFields(m, arguments)

		check(client.Identify(m))
	}

	if arguments["group"].(bool) {
		m := &analytics.Group{
			GroupId: arguments["--groupId"].(string),
		}
		traits := getOptionalString(arguments, "--traits")
		if traits != "" {
			var parsedTraits map[string]interface{}
			err := json.Unmarshal([]byte(traits), &parsedTraits)
			check(err)
			m.Traits = parsedTraits
		}

		setCommonFields(m, arguments)

		check(client.Group(m))
	}

	if arguments["alias"].(bool) {
		m := &analytics.Alias{
			PreviousId: arguments["--previousId"].(string),
		}

		setCommonFields(m, arguments)

		check(client.Alias(m))
	}

	client.Close()
}

func setCommonFields(message interface{}, arguments map[string]interface{}) {
	userId := getOptionalString(arguments, "--userId")
	if userId != "" {
		setFieldValue(message, "UserId", userId)
	}
	anonymousId := getOptionalString(arguments, "--anonymousId")
	if anonymousId != "" {
		setFieldValue(message, "AnonymousId", anonymousId)
	}
	integrations := getOptionalString(arguments, "--integrations")
	if integrations != "" {
		var parsedIntegrations map[string]interface{}
		err := json.Unmarshal([]byte(integrations), &parsedIntegrations)
		check(err)
		setFieldValue(message, "Integrations", parsedIntegrations)
	}
	context := getOptionalString(arguments, "--context")
	if context != "" {
		var parsedContext map[string]interface{}
		err := json.Unmarshal([]byte(context), &parsedContext)
		check(err)
		setFieldValue(message, "Context", parsedContext)

	}
	timestamp := getOptionalString(arguments, "--timestamp")
	if timestamp != "" {
		setFieldValue(message, "Timestamp", timestamp)
	}
}

func setFieldValue(target interface{}, field string, value interface{}) {
	reflect.ValueOf(target).Elem().FieldByName(field).Set(reflect.ValueOf(value))
}

func getOptionalString(m map[string]interface{}, k string) string {
	v := m[k]
	if v == nil {
		return ""
	}
	return v.(string)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
