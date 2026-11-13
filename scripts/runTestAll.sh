export LLM_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"
export IMAGE_QUERY_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"

export FILE_STORE_PATH=$PWD/dataVolume
export SERVABLE_DIRECTORY=$PWD/dataVolume/static

# this entry is for testing settings package
export SETTINGS_IGNORE_TESTCASE_ENDPOINT="Sample123"
export TTS_MODEL="mlx-community/Qwen3-TTS-12Hz-0.6B-Base-8bit"

go test ./pkg/...  -v
