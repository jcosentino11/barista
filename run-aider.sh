#!/usr/bin/env zsh

# ===============================================================================
# REQUIRED SETUP:
#   1) define API keys in .env file (https://aider.chat/docs/config/api-keys.html)
#   2) specify model in .aider-model file
# ===============================================================================

MODEL=$(cat .aider-model) 

aider --model ${MODEL} \
     --no-auto-commits \
     --cache-prompts \
     --watch-files \
     --analytics-disable
