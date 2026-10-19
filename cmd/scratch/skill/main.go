package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

// SkillMetadata represents the YAML frontmatter of the SKILL.md standard
type SkillMetadata struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Triggers    []string `yaml:"triggers"` // Intent matching keywords
}

// AgentSkill holds the parsed metadata and raw instruction blocks
type AgentSkill struct {
	Metadata     SkillMetadata
	Instructions string
}

// LoadSkill reads and parses the skill file split into frontmatter and markdown body
func LoadSkill(filePath string) (*AgentSkill, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var yamlLines []string
	var markdownLines []string
	inFrontmatter := false
	frontmatterRead := false

	for scanner.Scan() {
		line := scanner.Text()

		// Detect YAML frontmatter boundaries (---)
		if strings.TrimSpace(line) == "---" {
			if !inFrontmatter && !frontmatterRead {
				inFrontmatter = true
				continue
			} else if inFrontmatter {
				inFrontmatter = false
				frontmatterRead = true
				continue
			}
		}

		if inFrontmatter {
			yamlLines = append(yamlLines, line)
		} else if frontmatterRead {
			markdownLines = append(markdownLines, line)
		}
	}

	// Parse YAML metadata
	var meta SkillMetadata
	err = yaml.Unmarshal([]byte(strings.Join(yamlLines, "\n")), &meta)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	return &AgentSkill{
		Metadata:     meta,
		Instructions: strings.Join(markdownLines, "\n"),
	}, nil
}

// Invoke evaluates user input and executes instructions if matched
func (as *AgentSkill) Invoke(userInput string) {
	fmt.Printf("🤖 Agent checking skill availability: '%s'...\n", as.Metadata.Name)

	// Simple routing mechanism matching keywords/triggers
	matched := false
	for _, trigger := range as.Metadata.Triggers {
		if strings.Contains(strings.ToLower(userInput), strings.ToLower(trigger)) {
			matched = true
			break
		}
	}

	if !matched {
		fmt.Println("❌ User request does not match this skill's scope.")
		return
	}

	fmt.Println("⚡ Match found! Loading instructions into agent context...")
	as.executeInstructions()
}

// Parse markdown instruction blocks and look for executable commands
func (as *AgentSkill) executeInstructions() {
	scanner := bufio.NewScanner(strings.NewReader(as.Instructions))
	inCommandBlock := false

	for scanner.Scan() {
		line := scanner.Text()

		// Match markdown bash/shell code blocks
		if strings.HasPrefix(line, "```bash") || strings.HasPrefix(line, "```sh") {
			inCommandBlock = true
			continue
		}
		if strings.HasPrefix(line, "```") && inCommandBlock {
			inCommandBlock = false
			continue
		}

		// If we are inside an instruction code block, execute it
		if inCommandBlock && strings.TrimSpace(line) != "" {
			fmt.Printf("⚙️ Executing workflow step: %s\n", line)

			// Simple shell string execution
			cmd := exec.Command("sh", "-c", line)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				fmt.Printf("⚠️ Step failed: %v\n", err)
				return
			}
		}
	}
	fmt.Println("✅ All instructions in skill successfully performed.")
}

func main() {
	// 1. Initialize the agent skill from disk
	skillFile := "./dataVolume/skills/skill.md"
	skill, err := LoadSkill(skillFile)
	if err != nil {
		fmt.Printf("Error loading %s: %v\n", skillFile, err)
		return
	}

	// 2. Mocking an invocation (In reality, this comes from a prompt/CLI argument)
	userPrompt := "Hey agent, please run the security audit setup right now."
	// userPrompt := "Hey agent, please run the project description"
	fmt.Printf("👤 User: \"%s\"\n\n", userPrompt)
	skill.Invoke(userPrompt)
}
