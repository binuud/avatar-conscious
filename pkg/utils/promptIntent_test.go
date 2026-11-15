package utils_test

import (
	"strings"
	"testing"

	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// THIS TEST CASE FILE WAS CREATED BY AI

// Helper functions needed for testing (mocking or implementing if not exported in your main file)
// Note: If containsAny and hasImageContext are not exported in your main file,
// you must either export them or move the logic into the test file as standalone functions.
// Assuming they are exported (capitalized) or you are moving the logic into the test file.

// If your original code has these functions unexported (lowercase), you cannot test them directly
// unless you move them to a separate _test.go file or export them.
// For this example, I assume they are exported or we will inline the logic for the test.

func Test_DetectPromptIntent(t *testing.T) {
	tests := []struct {
		name       string
		prompt     string
		wantIntent utils.PromptIntent
	}{
		// --- EDIT INTENT TESTS ---
		{
			name:       "edit keyword with image context",
			prompt:     "edit this image",
			wantIntent: utils.IntentEdit,
		},
		{
			name:       "edit keyword with photo context",
			prompt:     "modify the background of this photo",
			wantIntent: utils.IntentEdit,
		},
		{
			name:       "edit keyword with drawing context",
			prompt:     "retouch this drawing",
			wantIntent: utils.IntentEdit,
		},
		{
			name:       "case insensitive edit",
			prompt:     "EDIT THE IMAGE",
			wantIntent: utils.IntentEdit,
		},
		{
			name:       "whitespace insensitive edit",
			prompt:     "   remove the background   ",
			wantIntent: utils.IntentEdit,
		},

		// Specific EDIT keywords check
		{
			name:       "inpaint keyword with context",
			prompt:     "inpaint the image",
			wantIntent: utils.IntentEdit,
		},
		{
			name:       "outpaint keyword with context",
			prompt:     "outpaint the picture",
			wantIntent: utils.IntentEdit,
		},
		{
			name:       "replace keyword with context",
			prompt:     "replace the logo in the image",
			wantIntent: utils.IntentEdit,
		},
		{
			name:       "crop keyword with context",
			prompt:     "crop the image",
			wantIntent: utils.IntentEdit,
		},

		// --- DESCRIBE INTENT TESTS ---
		{
			name:       "describe keyword with image context",
			prompt:     "describe this image",
			wantIntent: utils.IntentDescribe,
		},
		{
			name:       "what is in keyword with image context",
			prompt:     "what is in this photo",
			wantIntent: utils.IntentDescribe,
		},
		{
			name:       "explain keyword with image context",
			prompt:     "explain this drawing",
			wantIntent: utils.IntentDescribe,
		},
		{
			name:       "ocr keyword with image context",
			prompt:     "ocr this image",
			wantIntent: utils.IntentDescribe,
		},
		{
			name:       "case insensitive describe",
			prompt:     "DESCRIBE THE PHOTO",
			wantIntent: utils.IntentDescribe,
		},
		{
			name:       "whitespace insensitive describe",
			prompt:     "   caption this screenshot   ",
			wantIntent: utils.IntentDescribe,
		},

		// --- CREATE INTENT TESTS ---

		{
			name:       "create keyword with context",
			prompt:     "create an image of a cat",
			wantIntent: utils.IntentCreate,
		},
		{
			name:       "generate keyword without context",
			prompt:     "generate an image",
			wantIntent: utils.IntentCreate,
		},
		{
			name:       "draw keyword without context",
			prompt:     "draw a dragon",
			wantIntent: utils.IntentCreate,
		},
		{
			name:       "text to image keyword",
			prompt:     "text to image prompt",
			wantIntent: utils.IntentCreate,
		},
		{
			name:       "illustrate keyword",
			prompt:     "illustrate this scene",
			wantIntent: utils.IntentCreate,
		},

		// --- PRECEDENCE TESTS ---
		// EDIT should take precedence over DESCRIBE if keywords overlap
		{
			name:       "edit overrides describe (remove vs describe)",
			prompt:     "remove this image instead of describing it",
			wantIntent: utils.IntentEdit, // "remove" triggers edit
		},
		// DESCRIBE should take precedence over CREATE if keywords overlap
		{
			name:       "describe overrides create (caption vs generate)",
			prompt:     "caption this image but don't generate one",
			wantIntent: utils.IntentDescribe, // "caption" triggers describe
		},

		// --- NONE INTENT TESTS ---
		{
			name:       "no keywords, no context",
			prompt:     "hello world",
			wantIntent: utils.IntentNone,
		},
		{
			name: "create keyword but missing image context (should fail edit/describe, succeed create)",
			// Note: If the user types "generate an image" without "this image",
			// our logic says CREATE. If they type "generate" without "image", it's NONE.
			prompt:     "generate an image",
			wantIntent: utils.IntentCreate,
		},
		{
			name:       "describe keyword but missing image context",
			prompt:     "describe the concept of love",
			wantIntent: utils.IntentNone, // "describe" present, but no "image", "photo", etc.
		},
		{
			name:       "edit keyword but missing image context",
			prompt:     "edit this document", // "document" is not in nouns list
			wantIntent: utils.IntentNone,
		},
		{
			name:       "empty prompt",
			prompt:     "",
			wantIntent: utils.IntentNone,
		},
		{
			name:       "whitespace only prompt",
			prompt:     "   ",
			wantIntent: utils.IntentNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.DetectImageIntent(tt.prompt)
			assert.Equal(t, tt.wantIntent, got)
		})
	}
}

// Helper functions implementation for the test file if they are not exported elsewhere.
// Copy these into the same file if your main logic file doesn't export them.
func containsAny(text string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(text, sub) {
			return true
		}
	}
	return false
}

// hasImageContext ensures words like "image" or "remove" refer to a visual asset.
func hasImageContext(prompt string) bool {
	nouns := []string{"image", "picture", "photo", "pic", "drawing", "screenshot", "illustration"}
	return containsAny(prompt, nouns)
}
