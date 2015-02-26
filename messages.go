package webhookListener

import "github.com/mitchellh/mapstructure"

// GitlabPushMessage representation of a Gitlab push message
type GitlabPushMessage struct {
	Before            string
	After             string
	Ref               string
	UserID            int    `mapstructure:"user_id"`
	UserName          string `mapstructure:"user_name"`
	ProjectID         int    `mapstructure:"project_id"`
	TotalCommitsCount int    `mapstructure:"total_commits_count"`
	Repository        struct {
		Name        string
		URL         string `mapstructure:"url"`
		Description string
		Homepage    string
	}
	Commits []struct {
		ID        string `mapstructure:"id"`
		Message   string
		Timestamp string
		URL       string `mapstructure:"url"`
		Author    struct {
			Name  string
			Email string
		}
	}
}

// DecodeMessage from a generic map of string into a GitlabPushMessage structure
func DecodeMessage(message map[string]interface{}) (GitlabPushMessage, error) {
	var msg GitlabPushMessage

	err := mapstructure.Decode(message, &msg)

	return msg, err
}
