package utils

import "strings"

// ImageIntent represents the detected action for an image-related prompt.
type PromptIntent string

const (
	IntentCreate   PromptIntent = "CREATE_IMAGE"
	IntentDescribe PromptIntent = "DESCRIBE_IMAGE"
	IntentEdit     PromptIntent = "EDIT_IMAGE"
	IntentNone     PromptIntent = "NONE"
)

// DetectImageIntent analyzes a text prompt and returns the likely image intent.
func DetectImageIntent(prompt string) PromptIntent {
	p := strings.ToLower(strings.TrimSpace(prompt))

	// Check for EDIT intent first (often references an existing image)
	editKeywords := []string{
		"edit this", "edit the", "modify this", "modify the",
		"change the background", "remove the", "add a", "replace the",
		"inpaint", "outpaint", "retouch", "crop", "resize",
		"make the image", "alter this", "transform this",
	}
	if containsAny(p, editKeywords) && hasImageContext(p) {
		return IntentEdit
	}

	// Check for DESCRIBE intent
	describeKeywords := []string{
		"describe", "what is in", "what's in", "explain this",
		"analyze this", "identify", "caption", "tell me about this",
		"read the text in", "ocr", "what does this show",
	}
	if containsAny(p, describeKeywords) && hasImageContext(p) {
		return IntentDescribe
	}

	// Check for CREATE intent
	createKeywords := []string{
		"generate an image", "generate a picture", "generate image",
		"create image", "create an image",
		"create photograph", "create a photograph", "create photo", "create a photo",
		"create a picture", "draw", "paint", "render", "illustrate",
		"sketch", "design a logo", "make an image", "text to image",
		"photo of", "illustration of", "3d render of",
	}
	if containsAny(p, createKeywords) {
		return IntentCreate
	}

	return IntentNone
}

// hasImageContext ensures words like "describe" or "remove" refer to a visual asset.
func hasImageContext(prompt string) bool {
	nouns := []string{"image", "picture", "photo", "pic", "drawing", "screenshot", "illustration"}
	return containsAny(prompt, nouns)
}

func containsAny(text string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(text, sub) {
			return true
		}
	}
	return false
}
