#!/usr/bin/env zsh

# ===============================================================================
# REQUIRED SETUP:
#   1) define API keys in .env file (https://aider.chat/docs/config/api-keys.html)
#   2) specify model in .aider-model file
# ===============================================================================

MODEL=$(head -1 .aider-model) 
CHAT_MODE=ask

aider --model ${MODEL} \
     --chat-mode ${CHAT_MODE} \
     --no-auto-commits \
     --cache-prompts \
     --watch-files \
     --analytics-disable
