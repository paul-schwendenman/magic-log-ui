package agent

type AgentConfig struct {
	Name         string
	ServerAddr   string
	LogFormat    string
	ParseRegex   string
	JqFilter     string
	CSVFieldsStr string
	HasCSVHeader bool
	Echo         bool
}

func GenerateAgentName() string {
	// Could use hostname or uuid if needed
	return "agent01"
}
