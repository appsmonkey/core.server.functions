package airship

import (
	"encoding/json"
	"fmt"
)

// Message to be sent
type Message struct {
	data map[string]interface{}
}

// NewMessage to push to the clients
func newMessage(value string, lvl ChanelType, sensorValue string) []*Message {
	res := make([]*Message, 0)

	switch lvl {
	case Good:
		msg := new(Message)
		msg.data = map[string]interface{}{}
		msg.data["device_types"] = []string{"android", "ios"}
		msg.data["notification"] = map[string]interface{}{
			"alert": "Air quality is better now. " + sensorValue,
			"android": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_Good",
				},
			},
			"ios": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_Good",
				},
			},
		}
		msg.data["audience"] = map[string]interface{}{
			"OR": []map[string]string{
				{"tag": "last_Sensitive"},
				{"tag": "last_Unhealthy"},
				{"tag": "last_VeryUnhealthy"},
				{"tag": "last_Hazardous"},
			},
		}

		res = append(res, msg)

	case Sensitive:
		// Negative Message
		msg := new(Message)
		msg.data = map[string]interface{}{}
		msg.data["device_types"] = []string{"android", "ios"}
		msg.data["notification"] = map[string]interface{}{
			"alert": "Air quality is currently unhealthy for sensitive people, stay indoors. " + sensorValue,
			"android": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_Sensitive",
				},
			},
			"ios": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_Sensitive",
				},
			},
		}
		msg.data["audience"] = map[string]interface{}{
			"AND": []map[string]interface{}{
				{"tag": "setting_Sensitive"},
				{"tag": "last_Good"},
			},
		}

		res = append(res, msg)

		// Positive Message
		msg = new(Message)
		msg.data = map[string]interface{}{}
		msg.data["device_types"] = []string{"android", "ios"}
		msg.data["notification"] = map[string]interface{}{
			"alert": "Air quality is better now. " + sensorValue,
			"android": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_Sensitive",
				},
			},
			"ios": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_Sensitive",
				},
			},
		}
		msg.data["audience"] = map[string]interface{}{
			"AND": []map[string]interface{}{
				{
					"OR": []map[string]string{
						{"tag": "setting_Unhealthy"},
						{"tag": "setting_VeryUnhealthy"},
						{"tag": "setting_Hazardous"},
					},
				},
				{
					"OR": []map[string]string{
						{"tag": "last_Unhealthy"},
						{"tag": "last_VeryUnhealthy"},
						{"tag": "last_Hazardous"},
					},
				},
			},
		}

		res = append(res, msg)

	case Unhealthy:
		// Negative Message
		msg := new(Message)
		msg.data = map[string]interface{}{}
		msg.data["device_types"] = []string{"android", "ios"}
		msg.data["notification"] = map[string]interface{}{
			"alert": "Air quality is currently unhealthy, stay indoors. " + sensorValue,
			"android": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_Unhealthy",
				},
			},
			"ios": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_Unhealthy",
				},
			},
		}
		msg.data["audience"] = map[string]interface{}{
			"AND": []map[string]interface{}{
				{
					"OR": []map[string]string{
						{"tag": "setting_Sensitive"},
						{"tag": "setting_Unhealthy"},
					},
				},
				{
					"OR": []map[string]string{
						{"tag": "last_Good"},
						{"tag": "last_Sensitive"},
					},
				},
			},
		}

		res = append(res, msg)

		// Positive Message
		msg = new(Message)
		msg.data = map[string]interface{}{}
		msg.data["device_types"] = []string{"android", "ios"}
		msg.data["notification"] = map[string]interface{}{
			"alert": "Air quality is better now. " + sensorValue,
			"android": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_Unhealthy",
				},
			},
			"ios": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_Unhealthy",
				},
			},
		}
		msg.data["audience"] = map[string]interface{}{
			"AND": []map[string]interface{}{
				{
					"OR": []map[string]string{
						{"tag": "setting_VeryUnhealthy"},
						{"tag": "setting_Hazardous"},
					},
				},
				{
					"OR": []map[string]string{
						{"tag": "last_VeryUnhealthy"},
						{"tag": "last_Hazardous"},
					},
				},
			},
		}

		res = append(res, msg)

	case VeryUnhealthy:
		// Negative Message
		msg := new(Message)
		msg.data = map[string]interface{}{}
		msg.data["device_types"] = []string{"android", "ios"}
		msg.data["notification"] = map[string]interface{}{
			"alert": "Air quality is currently very unhealthy, stay indoors. " + sensorValue,
			"android": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_VeryUnhealthy",
				},
			},
			"ios": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_VeryUnhealthy",
				},
			},
		}
		msg.data["audience"] = map[string]interface{}{
			"AND": []map[string]interface{}{
				{
					"OR": []map[string]string{
						{"tag": "setting_Sensitive"},
						{"tag": "setting_Unhealthy"},
						{"tag": "setting_VeryUnhealthy"},
					},
				},
				{
					"OR": []map[string]string{
						{"tag": "last_Good"},
						{"tag": "last_Sensitive"},
						{"tag": "last_Unhealthy"},
					},
				},
			},
		}

		res = append(res, msg)

		// Positive Message
		msg = new(Message)
		msg.data = map[string]interface{}{}
		msg.data["device_types"] = []string{"android", "ios"}
		msg.data["notification"] = map[string]interface{}{
			"alert": "Air quality is better now. " + sensorValue,
			"android": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_VeryUnhealthy",
				},
			},
			"ios": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_VeryUnhealthy",
				},
			},
		}
		msg.data["audience"] = map[string]interface{}{
			"AND": []map[string]interface{}{
				{"tag": "setting_Hazardous"},
				{
					"OR": []map[string]string{
						{"tag": "last_Hazardous"},
					},
				},
			},
		}

		res = append(res, msg)

	case Hazardous:
		// Negative Message
		msg := new(Message)
		msg.data = map[string]interface{}{}
		msg.data["device_types"] = []string{"android", "ios"}
		msg.data["notification"] = map[string]interface{}{
			"alert": "Air quality is currently hazardous, stay indoors. " + sensorValue,
			"android": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_Hazardous",
				},
			},
			"ios": map[string]interface{}{
				"extra": map[string]string{
					"remove": "last_Good|last_Sensitive|last_Unhealthy|last_VeryUnhealthy|last_Hazardous",
					"add":    "last_Hazardous",
				},
			},
		}
		msg.data["audience"] = map[string]interface{}{
			"AND": []map[string]interface{}{
				{
					"OR": []map[string]string{
						{"tag": "setting_Sensitive"},
						{"tag": "setting_Unhealthy"},
						{"tag": "setting_VeryUnhealthy"},
						{"tag": "setting_Hazardous"},
					},
				},
				{
					"OR": []map[string]string{
						{"tag": "last_Good"},
						{"tag": "last_Sensitive"},
						{"tag": "last_Unhealthy"},
						{"tag": "last_VeryUnhealthy"},
					},
				},
			},
		}

		res = append(res, msg)
	}

	return res
}

// Marshal the message data
func (m *Message) Marshal() ([]byte, error) {
	b, err := json.Marshal(m.data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return b, nil
}
