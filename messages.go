package webhookListener

import "github.com/mitchellh/mapstructure"

// GitlabPushMessage representation of a Gitlab push message
type GitlabPushMessage struct {
	Before            string
	After             string
	Ref               string
	UserID            int    `json:"user_id"`
	UserName          string `json:"user_name"`
	ProjectID         int    `json:"project_id"`
	TotalCommitsCount int    `json:"total_commits_count"`
	Repository        struct {
		Name        string
		URL         string `json:"url"`
		Description string
		Homepage    string
	}
	Commits []struct {
		ID        string `json:"id"`
		Message   string
		Timestamp string
		URL       string `json:"url"`
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
