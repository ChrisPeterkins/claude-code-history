package data

import (
	"os"
	"path/filepath"
	"strings"
)

// LoadSubagents discovers subagent conversation files for a session.
func LoadSubagents(session *Session) ([]SubagentInfo, error) {
	// Subagents live in {sessionID}/subagents/agent-{id}.jsonl
	sessionDir := filepath.Join(filepath.Dir(session.FilePath), session.ID)
	subagentDir := filepath.Join(sessionDir, "subagents")

	entries, err := os.ReadDir(subagentDir)
	if err != nil {
		// No subagents directory is normal
		return nil, nil
	}

	var agents []SubagentInfo
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".jsonl") {
			continue
		}

		agentID := strings.TrimSuffix(name, ".jsonl")
		agents = append(agents, SubagentInfo{
			AgentID:  agentID,
			FilePath: filepath.Join(subagentDir, name),
		})
	}

	return agents, nil
}

// LoadSubagentMessages loads the full conversation from a subagent file.
func LoadSubagentMessages(agent *SubagentInfo) ([]Message, error) {
	f, err := os.Open(agent.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var messages []Message
	scanner := newScanner(f, true)

	for scanner.Scan() {
		msg := parseMessage(scanner.Bytes())
		if msg != nil {
			messages = append(messages, *msg)
		}
	}

	if err := scanner.Err(); err != nil {
		return messages, err
	}

	PairToolInteractions(messages)
	return messages, nil
}
